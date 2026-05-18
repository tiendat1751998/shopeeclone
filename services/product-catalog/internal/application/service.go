package application
import ("context"; "encoding/json"; "fmt"; "time"; "github.com/shopee-clone/shopee/services/product-catalog/internal/domain"; "github.com/shopee-clone/shopee/services/product-catalog/internal/infrastructure/redis"; "github.com/shopee-clone/shopee/services/product-catalog/internal/metrics"; "github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"; "go.opentelemetry.io/otel"; "go.opentelemetry.io/otel/attribute"; "go.uber.org/zap")

type CatalogService struct {
	productRepo  domain.ProductRepository
	skuRepo      domain.SKURepository
	categoryRepo domain.CategoryRepository
	attrRepo     domain.AttributeRepository
	mediaRepo    domain.ProductMediaRepository
	redis        *redis.Store
	productTTL   time.Duration
	categoryTTL  time.Duration
	publisher    EventPublisher
}

type EventPublisher interface { Publish(ctx context.Context, event *domain.CatalogEvent) error }

func NewCatalogService(pr domain.ProductRepository, sr domain.SKURepository, cr domain.CategoryRepository, ar domain.AttributeRepository, mr domain.ProductMediaRepository, rs *redis.Store, pt, ct time.Duration, pub EventPublisher) *CatalogService {
	return &CatalogService{productRepo: pr, skuRepo: sr, categoryRepo: cr, attrRepo: ar, mediaRepo: mr, redis: rs, productTTL: pt, categoryTTL: ct, publisher: pub}
}

func (s *CatalogService) CreateProduct(ctx context.Context, shopID, name, description, categoryID, currency, idempotencyKey string) (*domain.Product, error) {
	ctx, span := otel.Tracer("shopee-catalog").Start(ctx, "catalog.create_product"); defer span.End()
	if idempotencyKey != "" {
		if ok, _ := s.redis.CheckIdempotency(ctx, "create_product:"+idempotencyKey, 24*time.Hour); !ok {
			metrics.IdempotentRequests.Inc(); return nil, domain.ErrInvalidState
		}
	}
	p := domain.NewProduct(shopID, name, description, categoryID, currency)
	if err := s.productRepo.Create(ctx, p); err != nil { return nil, err }
	data, _ := json.Marshal(p); s.redis.SetProduct(ctx, p.ID, data, s.productTTL)
	metrics.ProductsCreated.Inc()
	if s.publisher != nil { s.publisher.Publish(ctx, &domain.CatalogEvent{EventType: domain.EventProductCreated, AggregateType: "product", AggregateID: p.ID, Payload: p, CreatedAt: time.Now()}) }
	return p, nil
}

func (s *CatalogService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	ctx, span := otel.Tracer("shopee-catalog").Start(ctx, "catalog.get_product"); defer span.End()
	if data, err := s.redis.GetProduct(ctx, id); err == nil && len(data) > 0 {
		metrics.CacheHitsTotal.WithLabelValues("catalog", "redis").Inc()
		var p domain.Product; json.Unmarshal(data, &p); return &p, nil
	}
	metrics.CacheMissesTotal.WithLabelValues("catalog", "redis").Inc()
	p, err := s.productRepo.FindByID(ctx, id); if err != nil { return nil, err }
	if p == nil { return nil, domain.ErrProductNotFound }
	data, _ := json.Marshal(p); s.redis.SetProduct(ctx, id, data, s.productTTL)
	return p, nil
}

func (s *CatalogService) UpdateProduct(ctx context.Context, id, name, description, categoryID string) error {
	p, err := s.productRepo.FindByID(ctx, id); if err != nil { return err }
	if p == nil { return domain.ErrProductNotFound }
	p.Update(name, description); if categoryID != "" { p.CategoryID = categoryID }
	if err := s.productRepo.Update(ctx, p); err != nil { return err }
	s.redis.DeleteProduct(ctx, id)
	metrics.ProductsUpdated.Inc()
	if s.publisher != nil { s.publisher.Publish(ctx, &domain.CatalogEvent{EventType: domain.EventProductUpdated, AggregateType: "product", AggregateID: id, CreatedAt: time.Now()}) }
	return nil
}

func (s *CatalogService) ArchiveProduct(ctx context.Context, id string) error {
	p, err := s.productRepo.FindByID(ctx, id); if err != nil { return err }
	if p == nil { return domain.ErrProductNotFound }
	if err := p.Archive(); err != nil { return err }
	if err := s.productRepo.Update(ctx, p); err != nil { return err }
	s.redis.DeleteProduct(ctx, id)
	if s.publisher != nil { s.publisher.Publish(ctx, &domain.CatalogEvent{EventType: domain.EventProductArchived, AggregateType: "product", AggregateID: id, CreatedAt: time.Now()}) }
	return nil
}

func (s *CatalogService) GetCategoryTree(ctx context.Context) ([]*domain.Category, error) {
	if data, err := s.redis.GetCategoryTree(ctx); err == nil && len(data) > 0 {
		var cats []*domain.Category; json.Unmarshal(data, &cats); return cats, nil
	}
	cats, err := s.categoryRepo.GetTree(ctx); if err != nil { return nil, err }
	data, _ := json.Marshal(cats); s.redis.SetCategoryTree(ctx, data, s.categoryTTL)
	return cats, nil
}

func (s *CatalogService) CreateCategory(ctx context.Context, parentID, name, slug string, level, sortOrder int) (*domain.Category, error) {
	c := domain.NewCategory(parentID, name, slug, level, sortOrder)
	if err := s.categoryRepo.Create(ctx, c); err != nil { return nil, err }
	s.redis.InvalidateCategory(ctx, c.ID)
	return c, nil
}

func (s *CatalogService) AddSKU(ctx context.Context, productID, skuCode, attributes string, price int64) (*domain.SKU, error) {
	sku := domain.NewSKU(productID, skuCode, attributes, price)
	if err := s.skuRepo.Create(ctx, sku); err != nil { return nil, err }
	s.redis.DeleteProduct(ctx, productID)
	metrics.SKUsCreated.Inc()
	if s.publisher != nil { s.publisher.Publish(ctx, &domain.CatalogEvent{EventType: domain.EventSKUUpdated, AggregateType: "sku", AggregateID: sku.ID, CreatedAt: time.Now()}) }
	return sku, nil
}

func (s *CatalogService) ListProductsByShop(ctx context.Context, shopID string, offset, limit int) ([]*domain.Product, int64, error) {
	return s.productRepo.FindByShopID(ctx, shopID, offset, limit)
}
