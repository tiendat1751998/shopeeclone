package domain

import "context"

// CartRepository defines cart data access
type CartRepository interface {
	FindByID(ctx context.Context, id string) (*Cart, error)
	FindByUserID(ctx context.Context, userID string) (*Cart, error)
	FindBySessionID(ctx context.Context, sessionID string) (*Cart, error)
	Create(ctx context.Context, cart *Cart) error
	Update(ctx context.Context, cart *Cart) error
	Delete(ctx context.Context, id string) error
	FindExpired(ctx context.Context, before string, limit int) ([]*Cart, error)
}

// CartItemRepository defines cart item data access
type CartItemRepository interface {
	FindByID(ctx context.Context, id string) (*CartItem, error)
	FindByCartID(ctx context.Context, cartID string) ([]*CartItem, error)
	FindByCartAndSKU(ctx context.Context, cartID, sku string) (*CartItem, error)
	Create(ctx context.Context, item *CartItem) error
	Update(ctx context.Context, item *CartItem) error
	Delete(ctx context.Context, id string) error
	DeleteByCartID(ctx context.Context, cartID string) error
	CountByCartID(ctx context.Context, cartID string) (int, error)
}

// CartSnapshotRepository defines snapshot data access
type CartSnapshotRepository interface {
	FindByID(ctx context.Context, id string) (*CartSnapshot, error)
	FindByCartID(ctx context.Context, cartID string) (*CartSnapshot, error)
	FindByIdempotencyKey(ctx context.Context, key string) (*CartSnapshot, error)
	Create(ctx context.Context, snapshot *CartSnapshot) error
	Delete(ctx context.Context, id string) error
}

// CartMergeHistoryRepository tracks merge history
type CartMergeHistoryRepository interface {
	Create(ctx context.Context, history *CartMergeHistory) error
	FindByUserID(ctx context.Context, userID string, limit int) ([]*CartMergeHistory, error)
}
