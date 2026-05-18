package mysql
import ("context"; "database/sql"; "fmt"; "github.com/jmoiron/sqlx"; "github.com/shopee-clone/shopee/services/product-catalog/internal/domain")

type ProductRepository struct{ db *sqlx.DB }
func NewProductRepository(db *sqlx.DB) *ProductRepository { return &ProductRepository{db: db} }
func (r *ProductRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	var p domain.Product; err := r.db.GetContext(ctx, &p, "SELECT * FROM products WHERE id = ? AND deleted_at IS NULL", id)
	if err == sql.ErrNoRows { return nil, nil }; if err != nil { return nil, err }; return &p, nil
}
func (r *ProductRepository) FindByShopID(ctx context.Context, shopID string, offset, limit int) ([]*domain.Product, int64, error) {
	var total int64; r.db.GetContext(ctx, &total, "SELECT COUNT(*) FROM products WHERE shop_id = ? AND deleted_at IS NULL", shopID)
	var products []*domain.Product; err := r.db.SelectContext(ctx, &products, "SELECT * FROM products WHERE shop_id = ? AND deleted_at = NULL ORDER BY created_at DESC LIMIT ? OFFSET ?", shopID, limit, offset)
	return products, total, err
}
func (r *ProductRepository) FindByCategory(ctx context.Context, categoryID string, offset, limit int) ([]*domain.Product, int64, error) {
	var total int64; r.db.GetContext(ctx, &total, "SELECT COUNT(*) FROM products WHERE category_id = ? AND status = 'active' AND deleted_at IS NULL", categoryID)
	var products []*domain.Product; err := r.db.SelectContext(ctx, &products, "SELECT * FROM products WHERE category_id = ? AND status = 'active' AND deleted_at IS NULL ORDER BY created_at DESC LIMIT ? OFFSET ?", categoryID, limit, offset)
	return products, total, err
}
func (r *ProductRepository) Create(ctx context.Context, p *domain.Product) error {
	query := `INSERT INTO products (id, shop_id, name, description, category_id, status, currency, version, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, p.ID, p.ShopID, p.Name, p.Description, p.CategoryID, p.Status, p.Currency, p.Version, p.CreatedAt, p.UpdatedAt); return err
}
func (r *ProductRepository) Update(ctx context.Context, p *domain.Product) error {
	query := `UPDATE products SET name = ?, description = ?, category_id = ?, status = ?, version = version + 1, updated_at = ? WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, p.Name, p.Description, p.CategoryID, p.Status, p.UpdatedAt, p.ID); return err
}
func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE products SET deleted_at = NOW() WHERE id = ?", id); return err
}

type SKURepository struct{ db *sqlx.DB }
func NewSKURepository(db *sqlx.DB) *SKURepository { return &SKURepository{db: db} }
func (r *SKURepository) FindByID(ctx context.Context, id string) (*domain.SKU, error) {
	var s domain.SKU; err := r.db.GetContext(ctx, &s, "SELECT * FROM skus WHERE id = ?", id)
	if err == sql.ErrNoRows { return nil, nil }; return &s, err
}
func (r *SKURepository) FindByProductID(ctx context.Context, productID string) ([]*domain.SKU, error) {
	var skus []*domain.SKU; err := r.db.SelectContext(ctx, &skus, "SELECT * FROM skus WHERE product_id = ? ORDER BY sort_order", productID)
	return skus, err
}
func (r *SKURepository) Create(ctx context.Context, sku *domain.SKU) error {
	query := `INSERT INTO skus (id, product_id, sku_code, attributes, price, sale_price, stock, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, sku.ID, sku.ProductID, sku.SKUCode, sku.Attributes, sku.Price, sku.SalePrice, sku.Stock, sku.Status, sku.CreatedAt, sku.UpdatedAt); return err
}
func (r *SKURepository) Update(ctx context.Context, sku *domain.SKU) error {
	query := `UPDATE skus SET attributes = ?, price = ?, sale_price = ?, stock = ?, status = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, sku.Attributes, sku.Price, sku.SalePrice, sku.Stock, sku.Status, sku.UpdatedAt, sku.ID); return err
}
func (r *SKURepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM skus WHERE id = ?", id); return err
}

type CategoryRepository struct{ db *sqlx.DB }
func NewCategoryRepository(db *sqlx.DB) *CategoryRepository { return &CategoryRepository{db: db} }
func (r *CategoryRepository) FindByID(ctx context.Context, id string) (*domain.Category, error) {
	var c domain.Category; err := r.db.GetContext(ctx, &c, "SELECT * FROM categories WHERE id = ?", id)
	if err == sql.ErrNoRows { return nil, nil }; return &c, err
}
func (r *CategoryRepository) FindByParentID(ctx context.Context, parentID string) ([]*domain.Category, error) {
	var cats []*domain.Category; err := r.db.SelectContext(ctx, &cats, "SELECT * FROM categories WHERE parent_id = ? AND is_active = true ORDER BY sort_order", parentID)
	return cats, err
}
func (r *CategoryRepository) GetTree(ctx context.Context) ([]*domain.Category, error) {
	var cats []*domain.Category; err := r.db.SelectContext(ctx, &cats, "SELECT * FROM categories WHERE is_active = true ORDER BY level, sort_order")
	return cats, err
}
func (r *CategoryRepository) Create(ctx context.Context, c *domain.Category) error {
	query := `INSERT INTO categories (id, parent_id, name, slug, level, sort_order, is_active, metadata, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, c.ID, c.ParentID, c.Name, c.Slug, c.Level, c.SortOrder, c.IsActive, c.Metadata, c.CreatedAt, c.UpdatedAt); return err
}
func (r *CategoryRepository) Update(ctx context.Context, c *domain.Category) error {
	query := `UPDATE categories SET name = ?, slug = ?, sort_order = ?, is_active = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, c.Name, c.Slug, c.SortOrder, c.IsActive, c.UpdatedAt, c.ID); return err
}
func (r *CategoryRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM categories WHERE id = ?", id); return err
}

type AttributeRepository struct{ db *sqlx.DB }
func NewAttributeRepository(db *sqlx.DB) *AttributeRepository { return &AttributeRepository{db: db} }
func (r *AttributeRepository) FindByCategoryID(ctx context.Context, categoryID string) ([]*domain.Attribute, error) {
	var attrs []*domain.Attribute; err := r.db.SelectContext(ctx, &attrs, "SELECT * FROM attributes WHERE category_id = ? AND is_active = true ORDER BY sort_order", categoryID)
	return attrs, err
}
func (r *AttributeRepository) Create(ctx context.Context, a *domain.Attribute) error {
	query := `INSERT INTO attributes (id, category_id, name, display_name, type, required, options, sort_order, is_active) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, a.ID, a.CategoryID, a.Name, a.DisplayName, a.AttrType, a.Required, a.Options, a.SortOrder, a.IsActive); return err
}
func (r *AttributeRepository) Update(ctx context.Context, a *domain.Attribute) error {
	query := `UPDATE attributes SET name = ?, display_name = ?, required = ?, options = ?, sort_order = ?, is_active = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, a.Name, a.DisplayName, a.Required, a.Options, a.SortOrder, a.IsActive, a.ID); return err
}

type ProductMediaRepository struct{ db *sqlx.DB }
func NewProductMediaRepository(db *sqlx.DB) *ProductMediaRepository { return &ProductMediaRepository{db: db} }
func (r *ProductMediaRepository) FindByProductID(ctx context.Context, productID string) ([]*domain.ProductMedia, error) {
	var media []*domain.ProductMedia; err := r.db.SelectContext(ctx, &media, "SELECT * FROM product_media WHERE product_id = ? ORDER BY sort_order", productID)
	return media, err
}
func (r *ProductMediaRepository) Create(ctx context.Context, m *domain.ProductMedia) error {
	query := `INSERT INTO product_media (id, product_id, media_type, url, thumbnail, sort_order, is_primary, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, m.ID, m.ProductID, m.MediaType, m.URL, m.Thumbnail, m.SortOrder, m.IsPrimary, m.CreatedAt); return err
}
func (r *ProductMediaRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM product_media WHERE id = ?", id); return err
}
