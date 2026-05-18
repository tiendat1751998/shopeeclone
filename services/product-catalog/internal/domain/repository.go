package domain

import "context"

type ProductRepository interface {
	FindByID(ctx context.Context, id string) (*Product, error)
	FindByShopID(ctx context.Context, shopID string, offset, limit int) ([]*Product, int64, error)
	FindByCategory(ctx context.Context, categoryID string, offset, limit int) ([]*Product, int64, error)
	Create(ctx context.Context, p *Product) error
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id string) error
}

type SKURepository interface {
	FindByID(ctx context.Context, id string) (*SKU, error)
	FindByProductID(ctx context.Context, productID string) ([]*SKU, error)
	Create(ctx context.Context, sku *SKU) error
	Update(ctx context.Context, sku *SKU) error
	Delete(ctx context.Context, id string) error
}

type CategoryRepository interface {
	FindByID(ctx context.Context, id string) (*Category, error)
	FindByParentID(ctx context.Context, parentID string) ([]*Category, error)
	GetTree(ctx context.Context) ([]*Category, error)
	Create(ctx context.Context, c *Category) error
	Update(ctx context.Context, c *Category) error
	Delete(ctx context.Context, id string) error
}

type AttributeRepository interface {
	FindByCategoryID(ctx context.Context, categoryID string) ([]*Attribute, error)
	Create(ctx context.Context, a *Attribute) error
	Update(ctx context.Context, a *Attribute) error
}

type ProductMediaRepository interface {
	FindByProductID(ctx context.Context, productID string) ([]*ProductMedia, error)
	Create(ctx context.Context, m *ProductMedia) error
	Delete(ctx context.Context, id string) error
}
