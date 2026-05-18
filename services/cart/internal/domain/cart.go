package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Cart represents a user's shopping cart
type Cart struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	SessionID string    `db:"session_id" json:"session_id,omitempty"`
	Status    string    `db:"status" json:"status"`
	Currency  string    `db:"currency" json:"currency"`
	ItemCount int       `db:"item_count" json:"item_count"`
	Subtotal  int64     `db:"subtotal" json:"subtotal"`
	Version   int64     `db:"version" json:"version"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

const (
	CartStatusActive   = "active"
	CartStatusMerged   = "merged"
	CartStatusExpired  = "expired"
	CartStatusCheckout = "checkout"
)

func NewCart(userID, sessionID, currency string, ttl time.Duration) *Cart {
	now := time.Now()
	return &Cart{
		ID:        uuid.New().String(),
		UserID:    userID,
		SessionID: sessionID,
		Status:    CartStatusActive,
		Currency:  currency,
		Version:   1,
		ExpiresAt: now.Add(ttl),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (c *Cart) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

func (c *Cart) IsActive() bool {
	return c.Status == CartStatusActive && !c.IsExpired()
}

func (c *Cart) Touch() {
	c.UpdatedAt = time.Now()
	c.Version++
}

func (c *Cart) MarkMerged() {
	c.Status = CartStatusMerged
	c.UpdatedAt = time.Now()
}

func (c *Cart) MarkCheckout() error {
	if c.Status != CartStatusActive {
		return fmt.Errorf("%w: cart not active", ErrInvalidCartState)
	}
	c.Status = CartStatusCheckout
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Cart) UpdateTotals(itemCount int, subtotal int64) {
	c.ItemCount = itemCount
	c.Subtotal = subtotal
	c.UpdatedAt = time.Now()
}

// CartItem represents an item in the cart
type CartItem struct {
	ID          string    `db:"id" json:"id"`
	CartID      string    `db:"cart_id" json:"cart_id"`
	SKU         string    `db:"sku" json:"sku"`
	ProductName string    `db:"product_name" json:"product_name"`
	ShopID      string    `db:"shop_id" json:"shop_id"`
	ShopName    string    `db:"shop_name" json:"shop_name"`
	Quantity    int       `db:"quantity" json:"quantity"`
	UnitPrice   int64     `db:"unit_price" json:"unit_price"`
	TotalPrice  int64     `db:"total_price" json:"total_price"`
	ImageURL    string    `db:"image_url" json:"image_url,omitempty"`
	Attributes  string    `db:"attributes" json:"attributes,omitempty"`
	IsSelected  bool      `db:"is_selected" json:"is_selected"`
	IsAvailable bool      `db:"is_available" json:"is_available"`
	AddedAt     time.Time `db:"added_at" json:"added_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

func NewCartItem(cartID, sku, productName, shopID, shopName string, quantity int, unitPrice int64, imageURL, attributes string) *CartItem {
	now := time.Now()
	return &CartItem{
		ID:          uuid.New().String(),
		CartID:      cartID,
		SKU:         sku,
		ProductName: productName,
		ShopID:      shopID,
		ShopName:    shopName,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
		TotalPrice:  int64(quantity) * unitPrice,
		ImageURL:    imageURL,
		Attributes:  attributes,
		IsSelected:  true,
		IsAvailable: true,
		AddedAt:     now,
		UpdatedAt:   now,
	}
}

func (ci *CartItem) UpdateQuantity(qty int) {
	ci.Quantity = qty
	ci.TotalPrice = int64(qty) * ci.UnitPrice
	ci.UpdatedAt = time.Now()
}

func (ci *CartItem) Recalculate() {
	ci.TotalPrice = int64(ci.Quantity) * ci.UnitPrice
}

// CartSnapshot is an immutable snapshot for checkout
type CartSnapshot struct {
	ID              string    `db:"id" json:"id"`
	CartID          string    `db:"cart_id" json:"cart_id"`
	UserID          string    `db:"user_id" json:"user_id"`
	Items           string    `db:"items" json:"items"`
	SellerGroups    string    `db:"seller_groups" json:"seller_groups"`
	Subtotal        int64     `db:"subtotal" json:"subtotal"`
	ItemCount       int       `db:"item_count" json:"item_count"`
	Currency        string    `db:"currency" json:"currency"`
	IdempotencyKey  string    `db:"idempotency_key" json:"idempotency_key"`
	ExpiresAt       time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

func NewCartSnapshot(cartID, userID string, items, sellerGroups string, subtotal int64, itemCount int, currency, idempotencyKey string, ttl time.Duration) *CartSnapshot {
	now := time.Now()
	return &CartSnapshot{
		ID:             uuid.New().String(),
		CartID:         cartID,
		UserID:         userID,
		Items:          items,
		SellerGroups:   sellerGroups,
		Subtotal:       subtotal,
		ItemCount:      itemCount,
		Currency:       currency,
		IdempotencyKey: idempotencyKey,
		ExpiresAt:      now.Add(ttl),
		CreatedAt:      now,
	}
}

// CheckoutPreview represents the checkout preparation result
type CheckoutPreview struct {
	ID              string           `json:"id"`
	CartID          string           `json:"cart_id"`
	UserID          string           `json:"user_id"`
	SellerGroups    []SellerGroup    `json:"seller_groups"`
	Subtotal        int64            `json:"subtotal"`
	DiscountTotal   int64            `json:"discount_total"`
	ShippingTotal   int64            `json:"shipping_total"`
	GrandTotal      int64            `json:"grand_total"`
	Currency        string           `json:"currency"`
	PromotionPreview []PromotionPreview `json:"promotion_preview"`
	InventoryStatus []InventoryStatus   `json:"inventory_status"`
	IdempotencyKey  string           `json:"idempotency_key"`
	ExpiresAt       time.Time        `json:"expires_at"`
	CreatedAt       time.Time        `json:"created_at"`
}

// SellerGroup groups cart items by seller
type SellerGroup struct {
	ShopID      string      `json:"shop_id"`
	ShopName    string      `json:"shop_name"`
	Items       []CartItem  `json:"items"`
	Subtotal    int64       `json:"subtotal"`
	ShippingFee int64       `json:"shipping_fee"`
	Discount    int64       `json:"discount"`
}

// PromotionPreview shows applicable promotions
type PromotionPreview struct {
	PromotionID   string `json:"promotion_id"`
	Type          string `json:"type"`
	Description   string `json:"description"`
	DiscountValue int64  `json:"discount_value"`
}

// InventoryStatus shows stock availability
type InventoryStatus struct {
	SKU       string `json:"sku"`
	Available bool   `json:"available"`
	Stock     int64  `json:"stock"`
}

// CartMergeHistory tracks cart merge operations
type CartMergeHistory struct {
	ID           string    `db:"id" json:"id"`
	SourceCartID string    `db:"source_cart_id" json:"source_cart_id"`
	TargetCartID string    `db:"target_cart_id" json:"target_cart_id"`
	UserID       string    `db:"user_id" json:"user_id"`
	MergeType    string    `db:"merge_type" json:"merge_type"`
	ItemsMerged  int       `db:"items_merged" json:"items_merged"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

const (
	MergeTypeGuestToUser = "guest_to_user"
	MergeTypeSession     = "session"
	MergeTypeConflict    = "conflict_resolution"
)

// Domain errors
var (
	ErrCartNotFound      = ErrCart("cart_not_found")
	ErrCartExpired       = ErrCart("cart_expired")
	ErrCartFull          = ErrCart("cart_full")
	ErrItemNotFound      = ErrCart("item_not_found")
	ErrInvalidQuantity   = ErrCart("invalid_quantity")
	ErrInvalidCartState  = ErrCart("invalid_cart_state")
	ErrDuplicateItem     = ErrCart("duplicate_item")
	ErrMergeConflict     = ErrCart("merge_conflict")
	ErrSnapshotNotFound  = ErrCart("snapshot_not_found")
	ErrIdempotencyConflict = ErrCart("idempotency_conflict")
)

type ErrCart string

func (e ErrCart) Error() string {
	return "cart: " + string(e)
}
