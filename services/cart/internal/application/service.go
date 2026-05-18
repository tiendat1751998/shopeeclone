package application

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopee-clone/shopee/services/cart/internal/domain"
	"github.com/shopee-clone/shopee/services/cart/internal/infrastructure/redis"
	"github.com/shopee-clone/shopee/services/cart/internal/metrics"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type CartService struct {
	cartRepo     domain.CartRepository
	itemRepo     domain.CartItemRepository
	snapshotRepo domain.CartSnapshotRepository
	mergeRepo    domain.CartMergeHistoryRepository
	redis        *redis.Store
	cartTTL      time.Duration
	maxCartItems int
	maxQuantity  int
	publisher    EventPublisher
}

type EventPublisher interface {
	Publish(ctx context.Context, event *domain.CartEvent) error
}

func NewCartService(
	cartRepo domain.CartRepository,
	itemRepo domain.CartItemRepository,
	snapshotRepo domain.CartSnapshotRepository,
	mergeRepo domain.CartMergeHistoryRepository,
	redisStore *redis.Store,
	cartTTL time.Duration,
	maxCartItems, maxQuantity int,
	publisher EventPublisher,
) *CartService {
	return &CartService{
		cartRepo:     cartRepo,
		itemRepo:     itemRepo,
		snapshotRepo: snapshotRepo,
		mergeRepo:    mergeRepo,
		redis:        redisStore,
		cartTTL:      cartTTL,
		maxCartItems: maxCartItems,
		maxQuantity:  maxQuantity,
		publisher:    publisher,
	}
}

// GetOrCreateCart gets existing cart or creates a new one
func (s *CartService) GetOrCreateCart(ctx context.Context, userID, sessionID, currency string) (*domain.Cart, error) {
	ctx, span := otel.Tracer("shopee-cart").Start(ctx, "cart.get_or_create")
	defer span.End()

	// Try to find existing active cart
	if userID != "" {
		cart, err := s.cartRepo.FindByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if cart != nil && cart.IsActive() {
			return cart, nil
		}
	}

	// Try session cart
	if sessionID != "" {
		cart, err := s.cartRepo.FindBySessionID(ctx, sessionID)
		if err != nil {
			return nil, err
		}
		if cart != nil && cart.IsActive() {
			return cart, nil
		}
	}

	// Create new cart
	cart := domain.NewCart(userID, sessionID, currency, s.cartTTL)
	if err := s.cartRepo.Create(ctx, cart); err != nil {
		return nil, fmt.Errorf("create cart: %w", err)
	}

	// Cache cart reference
	if userID != "" {
		s.redis.SetUserCart(ctx, userID, cart.ID, s.cartTTL)
	}
	if sessionID != "" {
		s.redis.SetSessionCart(ctx, sessionID, cart.ID, s.cartTTL)
	}

	metrics.CartsCreated.Inc()
	return cart, nil
}

// AddItem adds an item to the cart
func (s *CartService) AddItem(ctx context.Context, cartID string, req AddItemRequest) (*domain.CartItem, error) {
	ctx, span := otel.Tracer("shopee-cart").Start(ctx, "cart.add_item")
	defer span.End()

	span.SetAttributes(
		attribute.String("cart_id", cartID),
		attribute.String("sku", req.SKU),
		attribute.Int("quantity", req.Quantity),
	)

	cart, err := s.cartRepo.FindByID(ctx, cartID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, domain.ErrCartNotFound
	}
	if !cart.IsActive() {
		return nil, domain.ErrInvalidCartState
	}

	// Check cart capacity
	itemCount, err := s.itemRepo.CountByCartID(ctx, cartID)
	if err != nil {
		return nil, err
	}
	if itemCount >= s.maxCartItems {
		return nil, fmt.Errorf("%w: max %d items", domain.ErrCartFull, s.maxCartItems)
	}

	// Check if item already exists
	existing, err := s.itemRepo.FindByCartAndSKU(ctx, cartID, req.SKU)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		// Update quantity
		newQty := existing.Quantity + req.Quantity
		if newQty > s.maxQuantity {
			newQty = s.maxQuantity
		}
		existing.UpdateQuantity(newQty)
		if err := s.itemRepo.Update(ctx, existing); err != nil {
			return nil, err
		}
		s.recalculateCart(ctx, cart)
		s.redis.DeleteCart(ctx, cartID)
		metrics.ItemsUpdated.Inc()
		return existing, nil
	}

	// Validate quantity
	if req.Quantity <= 0 || req.Quantity > s.maxQuantity {
		return nil, fmt.Errorf("%w: quantity must be 1-%d", domain.ErrInvalidQuantity, s.maxQuantity)
	}

	// Create new item
	item := domain.NewCartItem(cartID, req.SKU, req.ProductName, req.ShopID, req.ShopName, req.Quantity, req.UnitPrice, req.ImageURL, req.Attributes)
	if err := s.itemRepo.Create(ctx, item); err != nil {
		return nil, fmt.Errorf("create cart item: %w", err)
	}

	s.recalculateCart(ctx, cart)
	s.redis.DeleteCart(ctx, cartID)

	metrics.ItemsAdded.Inc()

	if s.publisher != nil {
		s.publisher.Publish(ctx, &domain.CartEvent{
			EventType:     domain.EventItemAdded,
			AggregateType: "cart_item",
			AggregateID:   item.ID,
			Payload: domain.ItemAddedPayload{
				CartID: cartID, SKU: req.SKU, ProductName: req.ProductName,
				ShopID: req.ShopID, Quantity: req.Quantity, UnitPrice: req.UnitPrice,
			},
			CreatedAt: time.Now(),
		})
	}

	return item, nil
}

// UpdateItemQuantity updates the quantity of a cart item
func (s *CartService) UpdateItemQuantity(ctx context.Context, cartID, itemID string, quantity int) error {
	ctx, span := otel.Tracer("shopee-cart").Start(ctx, "cart.update_item")
	defer span.End()

	item, err := s.itemRepo.FindByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item == nil {
		return domain.ErrItemNotFound
	}
	if item.CartID != cartID {
		return domain.ErrItemNotFound
	}

	if quantity <= 0 {
		return s.RemoveItem(ctx, cartID, itemID)
	}
	if quantity > s.maxQuantity {
		quantity = s.maxQuantity
	}

	item.UpdateQuantity(quantity)
	if err := s.itemRepo.Update(ctx, item); err != nil {
		return err
	}

	cart, _ := s.cartRepo.FindByID(ctx, cartID)
	if cart != nil {
		s.recalculateCart(ctx, cart)
	}
	s.redis.DeleteCart(ctx, cartID)

	metrics.ItemsUpdated.Inc()
	return nil
}

// RemoveItem removes an item from the cart
func (s *CartService) RemoveItem(ctx context.Context, cartID, itemID string) error {
	ctx, span := otel.Tracer("shopee-cart").Start(ctx, "cart.remove_item")
	defer span.End()

	item, err := s.itemRepo.FindByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item == nil || item.CartID != cartID {
		return domain.ErrItemNotFound
	}

	if err := s.itemRepo.Delete(ctx, itemID); err != nil {
		return err
	}

	cart, _ := s.cartRepo.FindByID(ctx, cartID)
	if cart != nil {
		s.recalculateCart(ctx, cart)
	}
	s.redis.DeleteCart(ctx, cartID)

	metrics.ItemsRemoved.Inc()
	return nil
}

// ClearCart removes all items from the cart
func (s *CartService) ClearCart(ctx context.Context, cartID string) error {
	ctx, span := otel.Tracer("shopee-cart").Start(ctx, "cart.clear")
	defer span.End()

	if err := s.itemRepo.DeleteByCartID(ctx, cartID); err != nil {
		return err
	}

	cart, _ := s.cartRepo.FindByID(ctx, cartID)
	if cart != nil {
		cart.UpdateTotals(0, 0)
		s.cartRepo.Update(ctx, cart)
	}
	s.redis.DeleteCart(ctx, cartID)

	if s.publisher != nil {
		s.publisher.Publish(ctx, &domain.CartEvent{
			EventType: domain.EventCartCleared, AggregateType: "cart",
			AggregateID: cartID, CreatedAt: time.Now(),
		})
	}

	return nil
}

// MergeCarts merges a source cart into a target cart (e.g., guest -> user)
func (s *CartService) MergeCarts(ctx context.Context, sourceCartID, targetCartID, userID string) error {
	ctx, span := otel.Tracer("shopee-cart").Start(ctx, "cart.merge")
	defer span.End()

	sourceItems, err := s.itemRepo.FindByCartID(ctx, sourceCartID)
	if err != nil {
		return err
	}

	merged := 0
	for _, item := range sourceItems {
		existing, _ := s.itemRepo.FindByCartAndSKU(ctx, targetCartID, item.SKU)
		if existing != nil {
			newQty := existing.Quantity + item.Quantity
			if newQty > s.maxQuantity {
				newQty = s.maxQuantity
			}
			existing.UpdateQuantity(newQty)
			s.itemRepo.Update(ctx, existing)
		} else {
			item.CartID = targetCartID
			item.ID = ""
			s.itemRepo.Create(ctx, item)
		}
		merged++
	}

	// Mark source cart as merged
	sourceCart, _ := s.cartRepo.FindByID(ctx, sourceCartID)
	if sourceCart != nil {
		sourceCart.MarkMerged()
		s.cartRepo.Update(ctx, sourceCart)
	}

	// Recalculate target cart
	targetCart, _ := s.cartRepo.FindByID(ctx, targetCartID)
	if targetCart != nil {
		s.recalculateCart(ctx, targetCart)
	}
	s.redis.DeleteCart(ctx, targetCartID)

	// Record merge history
	mergeHistory := &domain.CartMergeHistory{
		ID:           fmt.Sprintf("merge_%d", time.Now().UnixNano()),
		SourceCartID: sourceCartID,
		TargetCartID: targetCartID,
		UserID:       userID,
		MergeType:    domain.MergeTypeGuestToUser,
		ItemsMerged:  merged,
		CreatedAt:    time.Now(),
	}
	s.mergeRepo.Create(ctx, mergeHistory)

	metrics.CartsMerged.Inc()

	if s.publisher != nil {
		s.publisher.Publish(ctx, &domain.CartEvent{
			EventType: domain.EventCartMerged, AggregateType: "cart",
			AggregateID: targetCartID,
			Payload: domain.CartMergedPayload{
				SourceCartID: sourceCartID, TargetCartID: targetCartID,
				UserID: userID, ItemsMerged: merged,
			},
			CreatedAt: time.Now(),
		})
	}

	observability.LogWithTrace(ctx).Info("carts merged",
		zap.String("source", sourceCartID),
		zap.String("target", targetCartID),
		zap.Int("items_merged", merged),
	)

	return nil
}

// PrepareCheckout creates a checkout preview with seller grouping
func (s *CartService) PrepareCheckout(ctx context.Context, cartID, userID, idempotencyKey string) (*domain.CheckoutPreview, error) {
	ctx, span := otel.Tracer("shopee-cart").Start(ctx, "cart.prepare_checkout")
	defer span.End()

	// Check idempotency
	if idempotencyKey != "" {
		if existing, err := s.snapshotRepo.FindByIdempotencyKey(ctx, idempotencyKey); err == nil && existing != nil {
			metrics.IdempotentRequests.Inc()
			// Return cached preview
			return s.buildPreviewFromSnapshot(ctx, existing)
		}
	}

	cart, err := s.cartRepo.FindByID(ctx, cartID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, domain.ErrCartNotFound
	}
	if !cart.IsActive() {
		return nil, domain.ErrInvalidCartState
	}

	items, err := s.itemRepo.FindByCartID(ctx, cartID)
	if err != nil {
		return nil, err
	}

	selectedItems := make([]*domain.CartItem, 0)
	for _, item := range items {
		if item.IsSelected && item.IsAvailable {
			selectedItems = append(selectedItems, item)
		}
	}

	if len(selectedItems) == 0 {
		return nil, fmt.Errorf("no items selected for checkout")
	}

	// Group by seller
	sellerGroups := groupBySeller(selectedItems)

	// Build preview
	preview := &domain.CheckoutPreview{
		ID:             fmt.Sprintf("preview_%d", time.Now().UnixNano()),
		CartID:         cartID,
		UserID:         userID,
		SellerGroups:   sellerGroups,
		Subtotal:       cart.Subtotal,
		Currency:       cart.Currency,
		IdempotencyKey: idempotencyKey,
		ExpiresAt:      time.Now().Add(15 * time.Minute),
		CreatedAt:      time.Now(),
	}

	// Create snapshot
	itemsJSON, _ := json.Marshal(selectedItems)
	sellerGroupsJSON, _ := json.Marshal(sellerGroups)
	snapshot := domain.NewCartSnapshot(cartID, userID, string(itemsJSON), string(sellerGroupsJSON), cart.Subtotal, len(selectedItems), cart.Currency, idempotencyKey, 15*time.Minute)
	s.snapshotRepo.Create(ctx, snapshot)

	// Cache preview
	previewData, _ := json.Marshal(preview)
	s.redis.SetCheckoutPreview(ctx, preview.ID, previewData, 15*time.Minute)

	metrics.CheckoutPreviewsCreated.Inc()

	if s.publisher != nil {
		s.publisher.Publish(ctx, &domain.CartEvent{
			EventType: domain.EventCheckoutPrepared, AggregateType: "cart",
			AggregateID: cartID,
			Payload: domain.CheckoutPreparedPayload{
				CartID: cartID, UserID: userID, Subtotal: cart.Subtotal, ItemCount: len(selectedItems),
			},
			CreatedAt: time.Now(),
		})
	}

	return preview, nil
}

// GetCartWithItems retrieves a cart with all its items (cache-first)
func (s *CartService) GetCartWithItems(ctx context.Context, cartID string) (*domain.Cart, []*domain.CartItem, error) {
	cart, err := s.cartRepo.FindByID(ctx, cartID)
	if err != nil {
		return nil, nil, err
	}
	if cart == nil {
		return nil, nil, domain.ErrCartNotFound
	}

	items, err := s.itemRepo.FindByCartID(ctx, cartID)
	if err != nil {
		return nil, nil, err
	}

	return cart, items, nil
}

// recalculateCart updates cart totals based on items
func (s *CartService) recalculateCart(ctx context.Context, cart *domain.Cart) {
	items, err := s.itemRepo.FindByCartID(ctx, cart.ID)
	if err != nil {
		observability.LogWithTrace(ctx).Warn("failed to recalculate cart", zap.Error(err))
		return
	}

	total := int64(0)
	count := 0
	for _, item := range items {
		if item.IsSelected {
			total += item.TotalPrice
			count += item.Quantity
		}
	}
	cart.UpdateTotals(len(items), total)
	s.cartRepo.Update(ctx, cart)
}

// buildPreviewFromSnapshot rebuilds a preview from a cached snapshot
func (s *CartService) buildPreviewFromSnapshot(ctx context.Context, snapshot *domain.CartSnapshot) (*domain.CheckoutPreview, error) {
	var items []domain.CartItem
	if err := json.Unmarshal([]byte(snapshot.Items), &items); err != nil {
		return nil, err
	}

	var sellerGroups []domain.SellerGroup
	if err := json.Unmarshal([]byte(snapshot.SellerGroups), &sellerGroups); err != nil {
		return nil, err
	}

	return &domain.CheckoutPreview{
		ID:            snapshot.ID,
		CartID:        snapshot.CartID,
		UserID:        snapshot.UserID,
		SellerGroups:  sellerGroups,
		Subtotal:      snapshot.Subtotal,
		Currency:      snapshot.Currency,
		IdempotencyKey: snapshot.IdempotencyKey,
		ExpiresAt:     snapshot.ExpiresAt,
		CreatedAt:     snapshot.CreatedAt,
	}, nil
}

// groupBySeller groups cart items by their shop
func groupBySeller(items []*domain.CartItem) []domain.SellerGroup {
	groupMap := make(map[string]*domain.SellerGroup)
	for _, item := range items {
		if group, ok := groupMap[item.ShopID]; ok {
			group.Items = append(group.Items, *item)
			group.Subtotal += item.TotalPrice
		} else {
			groupMap[item.ShopID] = &domain.SellerGroup{
				ShopID:   item.ShopID,
				ShopName: item.ShopName,
				Items:    []domain.CartItem{*item},
				Subtotal: item.TotalPrice,
			}
		}
	}

	groups := make([]domain.SellerGroup, 0, len(groupMap))
	for _, g := range groupMap {
		groups = append(groups, *g)
	}
	return groups
}

// Request types

type AddItemRequest struct {
	SKU         string
	ProductName string
	ShopID      string
	ShopName    string
	Quantity    int
	UnitPrice   int64
	ImageURL    string
	Attributes  string
}
