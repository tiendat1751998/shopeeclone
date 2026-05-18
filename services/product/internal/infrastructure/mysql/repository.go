package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shopee-clone/shopee/services/product/internal/domain"
)

// ProductRepo implements ProductRepository using MySQL
type ProductRepo struct {
	db *sqlx.DB
}

// NewProductRepo creates a new ProductRepo
func NewProductRepo(db *sqlx.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

// Create inserts a new product with its SKUs and images in a transaction
func (r *ProductRepo) Create(ctx context.Context, product *domain.Product) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	productQuery := `INSERT INTO products (spu_id, title, description, category_id, brand_id, seller_id, status, created_at, updated_at)
		VALUES (:spu_id, :title, :description, :category_id, :brand_id, :seller_id, :status, :created_at, :updated_at)`
	if _, err := tx.NamedExecContext(ctx, productQuery, product); err != nil {
		return fmt.Errorf("insert product: %w", err)
	}

	// Insert SKUs
	for i := range product.SKUs {
		product.SKUs[i].SPUID = product.SPUID
		skuQuery := `INSERT INTO skus (sku_id, spu_id, price, sale_price, stock, weight, length, width, height, status, created_at, updated_at)
			VALUES (:sku_id, :spu_id, :price, :sale_price, :stock, :weight, :length, :width, :height, :status, :created_at, :updated_at)`
		if _, err := tx.NamedExecContext(ctx, skuQuery, &product.SKUs[i]); err != nil {
			return fmt.Errorf("insert sku: %w", err)
		}
	}

	// Insert images
	for i := range product.Images {
		product.Images[i].SPUID = product.SPUID
		imgQuery := `INSERT INTO product_images (spu_id, url, alt_text, sort_order, is_primary, created_at)
			VALUES (:spu_id, :url, :alt_text, :sort_order, :is_primary, :created_at)`
		if _, err := tx.NamedExecContext(ctx, imgQuery, &product.Images[i]); err != nil {
			return fmt.Errorf("insert image: %w", err)
		}
	}

	return tx.Commit()
}

// GetByID gets a product by internal ID
func (r *ProductRepo) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	var product domain.Product
	query := `SELECT id, spu_id, title, description, category_id, brand_id, seller_id, status, created_at, updated_at, deleted_at
		FROM products WHERE id = ? AND deleted_at IS NULL`
	if err := r.db.GetContext(ctx, &product, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get product by id: %w", err)
	}
	return r.loadRelations(ctx, &product)
}

// GetBySPU gets a product by SPU ID
func (r *ProductRepo) GetBySPU(ctx context.Context, spuID string) (*domain.Product, error) {
	var product domain.Product
	query := `SELECT id, spu_id, title, description, category_id, brand_id, seller_id, status, created_at, updated_at, deleted_at
		FROM products WHERE spu_id = ? AND deleted_at IS NULL`
	if err := r.db.GetContext(ctx, &product, query, spuID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get product by spu: %w", err)
	}
	return r.loadRelations(ctx, &product)
}

// List returns products with filtering and pagination
func (r *ProductRepo) List(ctx context.Context, filter domain.ProductFilter) (*domain.ProductList, error) {
	where := []string{"deleted_at IS NULL"}
	args := []interface{}{}

	if filter.CategoryID != "" {
		where = append(where, "category_id = ?")
		args = append(args, filter.CategoryID)
	}
	if filter.SellerID != "" {
		where = append(where, "seller_id = ?")
		args = append(args, filter.SellerID)
	}
	if filter.Status != "" {
		where = append(where, "status = ?")
		args = append(args, filter.Status)
	}
	if filter.MinPrice > 0 {
		where = append(where, "EXISTS (SELECT 1 FROM skus WHERE skus.spu_id = products.spu_id AND skus.price >= ?)")
		args = append(args, filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		where = append(where, "EXISTS (SELECT 1 FROM skus WHERE skus.spu_id = products.spu_id AND skus.price <= ?)")
		args = append(args, filter.MaxPrice)
	}

	whereClause := strings.Join(where, " AND ")

	// Count total
	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products WHERE %s", whereClause)
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, fmt.Errorf("count products: %w", err)
	}

	// Fetch page
	offset := (filter.Page - 1) * filter.Size
	query := fmt.Sprintf(`SELECT id, spu_id, title, description, category_id, brand_id, seller_id, status, created_at, updated_at
		FROM products WHERE %s ORDER BY %s %s LIMIT ? OFFSET ?`, whereClause, filter.SortBy, filter.SortOrder)
	args = append(args, filter.Size, offset)

	var products []domain.Product
	if err := r.db.SelectContext(ctx, &products, query, args...); err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}

	// Load relations for each product
	for i := range products {
		if _, err := r.loadRelations(ctx, &products[i]); err != nil {
			return nil, err
		}
	}

	return &domain.ProductList{
		Products: products,
		Total:    total,
		Page:     filter.Page,
		Size:     filter.Size,
	}, nil
}

// Update updates a product
func (r *ProductRepo) Update(ctx context.Context, product *domain.Product) error {
	query := `UPDATE products SET title = :title, description = :description, category_id = :category_id,
		brand_id = :brand_id, status = :status, updated_at = :updated_at WHERE spu_id = :spu_id AND deleted_at IS NULL`
	_, err := r.db.NamedExecContext(ctx, query, product)
	return err
}

// Delete soft-deletes a product
func (r *ProductRepo) Delete(ctx context.Context, id string) error {
	query := `UPDATE products SET deleted_at = ?, status = ? WHERE spu_id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, time.Now(), domain.ProductStatusDeleted, id)
	return err
}

// Search searches products by text (delegates to OpenSearch in production, fallback to LIKE here)
func (r *ProductRepo) Search(ctx context.Context, query string, filter domain.ProductFilter) (*domain.ProductList, error) {
	where := []string{"deleted_at IS NULL"}
	args := []interface{}{}

	if query != "" {
		where = append(where, "(title LIKE ? OR description LIKE ?)")
		args = append(args, "%"+query+"%", "%"+query+"%")
	}
	if filter.CategoryID != "" {
		where = append(where, "category_id = ?")
		args = append(args, filter.CategoryID)
	}
	if filter.SellerID != "" {
		where = append(where, "seller_id = ?")
		args = append(args, filter.SellerID)
	}

	whereClause := strings.Join(where, " AND ")

	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products WHERE %s", whereClause)
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, fmt.Errorf("count search: %w", err)
	}

	offset := (filter.Page - 1) * filter.Size
	selectQuery := fmt.Sprintf(`SELECT id, spu_id, title, description, category_id, brand_id, seller_id, status, created_at, updated_at
		FROM products WHERE %s ORDER BY updated_at DESC LIMIT ? OFFSET ?`, whereClause)
	args = append(args, filter.Size, offset)

	var products []domain.Product
	if err := r.db.SelectContext(ctx, &products, selectQuery, args...); err != nil {
		return nil, fmt.Errorf("search products: %w", err)
	}

	for i := range products {
		if _, err := r.loadRelations(ctx, &products[i]); err != nil {
			return nil, err
		}
	}

	return &domain.ProductList{
		Products: products,
		Total:    total,
		Page:     filter.Page,
		Size:     filter.Size,
	}, nil
}

// GetSKU gets a single SKU by ID
func (r *ProductRepo) GetSKU(ctx context.Context, skuID string) (*domain.SKU, error) {
	var sku domain.SKU
	query := `SELECT id, sku_id, spu_id, price, sale_price, stock, weight, length, width, height, status, created_at, updated_at
		FROM skus WHERE sku_id = ?`
	if err := r.db.GetContext(ctx, &sku, query, skuID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get sku: %w", err)
	}
	return &sku, nil
}

// BatchGetSKUs gets multiple SKUs by ID
func (r *ProductRepo) BatchGetSKUs(ctx context.Context, skuIDs []string) (map[string]*domain.SKU, error) {
	if len(skuIDs) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(skuIDs))
	args := make([]interface{}, len(skuIDs))
	for i, id := range skuIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`SELECT id, sku_id, spu_id, price, sale_price, stock, weight, length, width, height, status, created_at, updated_at
		FROM skus WHERE sku_id IN (%s)`, strings.Join(placeholders, ","))

	var skus []domain.SKU
	if err := r.db.SelectContext(ctx, &skus, query, args...); err != nil {
		return nil, fmt.Errorf("batch get skus: %w", err)
	}

	result := make(map[string]*domain.SKU, len(skus))
	for i := range skus {
		result[skus[i].SKUID] = &skus[i]
	}
	return result, nil
}

// CreateSKU creates a new SKU
func (r *ProductRepo) CreateSKU(ctx context.Context, sku *domain.SKU) error {
	query := `INSERT INTO skus (sku_id, spu_id, price, sale_price, stock, weight, length, width, height, status, created_at, updated_at)
		VALUES (:sku_id, :spu_id, :price, :sale_price, :stock, :weight, :length, :width, :height, :status, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, sku)
	return err
}

// UpdateSKU updates a SKU
func (r *ProductRepo) UpdateSKU(ctx context.Context, sku *domain.SKU) error {
	query := `UPDATE skus SET price = :price, sale_price = :sale_price, stock = :stock, status = :status, updated_at = :updated_at
		WHERE sku_id = :sku_id`
	_, err := r.db.NamedExecContext(ctx, query, sku)
	return err
}

// ListSKUsByProduct lists all SKUs for a product
func (r *ProductRepo) ListSKUsByProduct(ctx context.Context, spuID string) ([]domain.SKU, error) {
	query := `SELECT id, sku_id, spu_id, price, sale_price, stock, weight, length, width, height, status, created_at, updated_at
		FROM skus WHERE spu_id = ? ORDER BY created_at ASC`
	var skus []domain.SKU
	if err := r.db.SelectContext(ctx, &skus, query, spuID); err != nil {
		return nil, fmt.Errorf("list skus: %w", err)
	}
	return skus, nil
}

// loadRelations loads SKUs and images for a product
func (r *ProductRepo) loadRelations(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	// Load SKUs
	skus, err := r.ListSKUsByProduct(ctx, product.SPUID)
	if err != nil {
		return nil, err
	}
	product.SKUs = skus

	// Load images
	var images []domain.ProductImage
	imgQuery := `SELECT id, spu_id, url, alt_text, sort_order, is_primary, created_at FROM product_images
		WHERE spu_id = ? ORDER BY sort_order ASC`
	if err := r.db.SelectContext(ctx, &images, imgQuery, product.SPUID); err != nil {
		return nil, fmt.Errorf("load images: %w", err)
	}
	product.Images = images

	// Load attributes
	var attrValues []struct {
		AttributeID string `db:"attribute_id"`
		ValueID     string `db:"value_id"`
		CustomValue string `db:"custom_value"`
	}
	attrQuery := `SELECT attribute_id, value_id, custom_value FROM product_attribute_values WHERE spu_id = ?`
	if err := r.db.SelectContext(ctx, &attrValues, attrQuery, product.SPUID); err != nil {
		return nil, fmt.Errorf("load attributes: %w", err)
	}

	return product, nil
}

// CategoryRepo implements CategoryRepository using MySQL
type CategoryRepo struct {
	db *sqlx.DB
}

// NewCategoryRepo creates a new CategoryRepo
func NewCategoryRepo(db *sqlx.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) Create(ctx context.Context, category *domain.Category) error {
	query := `INSERT INTO categories (category_id, name, slug, parent_id, level, sort_order, image_url, is_active, created_at, updated_at)
		VALUES (:category_id, :name, :slug, :parent_id, :level, :sort_order, :image_url, :is_active, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, category)
	return err
}

func (r *CategoryRepo) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	var category domain.Category
	query := `SELECT id, category_id, name, slug, parent_id, level, sort_order, image_url, is_active, created_at, updated_at
		FROM categories WHERE category_id = ? AND is_active = true`
	if err := r.db.GetContext(ctx, &category, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get category: %w", err)
	}
	return &category, nil
}

func (r *CategoryRepo) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	var category domain.Category
	query := `SELECT id, category_id, name, slug, parent_id, level, sort_order, image_url, is_active, created_at, updated_at
		FROM categories WHERE slug = ? AND is_active = true`
	if err := r.db.GetContext(ctx, &category, query, slug); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get category by slug: %w", err)
	}
	return &category, nil
}

func (r *CategoryRepo) GetTree(ctx context.Context) (*domain.CategoryTree, error) {
	query := `SELECT id, category_id, name, slug, parent_id, level, sort_order, image_url, is_active, created_at, updated_at
		FROM categories WHERE is_active = true ORDER BY level ASC, sort_order ASC`
	var categories []domain.Category
	if err := r.db.SelectContext(ctx, &categories, query); err != nil {
		return nil, fmt.Errorf("get categories: %w", err)
	}

	tree := &domain.CategoryTree{}
	nodeMap := make(map[string]*domain.CategoryTreeNode)

	for i := range categories {
		node := &domain.CategoryTreeNode{
			Category: categories[i],
			Children: nil,
		}
		nodeMap[categories[i].CategoryID] = node
	}

	for _, node := range nodeMap {
		if node.ParentID == "" {
			tree.Roots = append(tree.Roots, node)
		} else {
			if parent, ok := nodeMap[node.ParentID]; ok {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	return tree, nil
}

func (r *CategoryRepo) List(ctx context.Context, parentID string) ([]domain.Category, error) {
	var categories []domain.Category
	var query string
	var args []interface{}
	if parentID == "" {
		query = `SELECT id, category_id, name, slug, parent_id, level, sort_order, image_url, is_active, created_at, updated_at
			FROM categories WHERE parent_id IS NULL AND is_active = true ORDER BY sort_order ASC`
	} else {
		query = `SELECT id, category_id, name, slug, parent_id, level, sort_order, image_url, is_active, created_at, updated_at
			FROM categories WHERE parent_id = ? AND is_active = true ORDER BY sort_order ASC`
		args = append(args, parentID)
	}
	if err := r.db.SelectContext(ctx, &categories, query, args...); err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	return categories, nil
}

func (r *CategoryRepo) Update(ctx context.Context, category *domain.Category) error {
	query := `UPDATE categories SET name = :name, slug = :slug, parent_id = :parent_id, sort_order = :sort_order,
		image_url = :image_url, is_active = :is_active, updated_at = :updated_at WHERE category_id = :category_id`
	_, err := r.db.NamedExecContext(ctx, query, category)
	return err
}

func (r *CategoryRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE categories SET is_active = false, updated_at = ? WHERE category_id = ?", time.Now(), id)
	return err
}

// AttributeRepo implements attribute repository using MySQL
type AttributeRepo struct {
	db *sqlx.DB
}

func NewAttributeRepo(db *sqlx.DB) *AttributeRepo {
	return &AttributeRepo{db: db}
}

func (r *AttributeRepo) Create(ctx context.Context, attr *domain.Attribute) error {
	query := `INSERT INTO attributes (attribute_id, category_id, name, type, is_required, is_filterable, is_searchable, sort_order, created_at, updated_at)
		VALUES (:attribute_id, :category_id, :name, :type, :is_required, :is_filterable, :is_searchable, :sort_order, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, attr)
	return err
}

func (r *AttributeRepo) GetByID(ctx context.Context, id string) (*domain.Attribute, error) {
	var attr domain.Attribute
	query := `SELECT * FROM attributes WHERE attribute_id = ?`
	if err := r.db.GetContext(ctx, &attr, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &attr, nil
}

func (r *AttributeRepo) ListByCategory(ctx context.Context, categoryID string) ([]domain.Attribute, error) {
	query := `SELECT * FROM attributes WHERE category_id = ? ORDER BY sort_order ASC`
	var attrs []domain.Attribute
	if err := r.db.SelectContext(ctx, &attrs, query, categoryID); err != nil {
		return nil, err
	}
	return attrs, nil
}

func (r *AttributeRepo) Update(ctx context.Context, attr *domain.Attribute) error {
	query := `UPDATE attributes SET name = :name, type = :type, is_required = :is_required, is_filterable = :is_filterable,
		is_searchable = :is_searchable, sort_order = :sort_order, updated_at = :updated_at WHERE attribute_id = :attribute_id`
	_, err := r.db.NamedExecContext(ctx, query, attr)
	return err
}

func (r *AttributeRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM attributes WHERE attribute_id = ?", id)
	return err
}

func (r *AttributeRepo) CreateValue(ctx context.Context, val *domain.AttributeValue) error {
	query := `INSERT INTO attribute_values (attribute_id, value, display_value, sort_order, created_at)
		VALUES (:attribute_id, :value, :display_value, :sort_order, :created_at)`
	_, err := r.db.NamedExecContext(ctx, query, val)
	return err
}

func (r *AttributeRepo) ListValues(ctx context.Context, attributeID string) ([]domain.AttributeValue, error) {
	query := `SELECT * FROM attribute_values WHERE attribute_id = ? ORDER BY sort_order ASC`
	var values []domain.AttributeValue
	if err := r.db.SelectContext(ctx, &values, query, attributeID); err != nil {
		return nil, err
	}
	return values, nil
}

// ModerationRepo implements moderation repository
type ModerationRepo struct {
	db *sqlx.DB
}

func NewModerationRepo(db *sqlx.DB) *ModerationRepo {
	return &ModerationRepo{db: db}
}

func (r *ModerationRepo) Create(ctx context.Context, record *domain.ModerationRecord) error {
	query := `INSERT INTO moderation_records (spu_id, status, reason, reviewer_id, created_at, updated_at)
		VALUES (:spu_id, :status, :reason, :reviewer_id, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, record)
	return err
}

func (r *ModerationRepo) GetByProduct(ctx context.Context, spuID string) (*domain.ModerationRecord, error) {
	var record domain.ModerationRecord
	query := `SELECT * FROM moderation_records WHERE spu_id = ? ORDER BY created_at DESC LIMIT 1`
	if err := r.db.GetContext(ctx, &record, query, spuID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

func (r *ModerationRepo) Update(ctx context.Context, record *domain.ModerationRecord) error {
	query := `UPDATE moderation_records SET status = :status, reason = :reason, reviewer_id = :reviewer_id, updated_at = :updated_at
		WHERE spu_id = :spu_id AND created_at = :created_at`
	_, err := r.db.NamedExecContext(ctx, query, record)
	return err
}

// Ensure interfaces are satisfied
var _ ProductRepository = (*ProductRepo)(nil)
var _ CategoryRepository = (*CategoryRepo)(nil)

// ProductRepository interface matching
type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	GetBySPU(ctx context.Context, spuID string) (*domain.Product, error)
	List(ctx context.Context, filter domain.ProductFilter) (*domain.ProductList, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, query string, filter domain.ProductFilter) (*domain.ProductList, error)
	GetSKU(ctx context.Context, skuID string) (*domain.SKU, error)
	BatchGetSKUs(ctx context.Context, skuIDs []string) (map[string]*domain.SKU, error)
	CreateSKU(ctx context.Context, sku *domain.SKU) error
	UpdateSKU(ctx context.Context, sku *domain.SKU) error
	ListSKUsByProduct(ctx context.Context, spuID string) ([]domain.SKU, error)
}

// jsonMarshal is a helper for event marshaling
func jsonMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
