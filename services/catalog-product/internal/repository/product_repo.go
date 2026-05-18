package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
	"github.com/shopee-clone/shopee/services/catalog-product/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

const productsCollection = "products"

type ProductRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewProductRepository(client *mongo.Client, dbName string) *ProductRepository {
	db := client.Database(dbName)
	return &ProductRepository{
		db:         db,
		collection: db.Collection(productsCollection),
	}
}

func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
	ctx, span := otel.Tracer("catalog-product").Start(ctx, "repository.product.create")
	defer span.End()

	if product.SPUID == "" {
		product.SPUID = uuid.New().String()
	}
	product.Status = domain.ProductStatusActive
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	for i := range product.SKUs {
		if product.SKUs[i].SKUID == "" {
			product.SKUs[i].SKUID = uuid.New().String()
		}
		product.SKUs[i].SPUID = product.SPUID
		product.SKUs[i].Status = domain.SKUStatusActive
	}

	span.SetAttributes(attribute.String("spu_id", product.SPUID))

	_, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		observability.LogWithTrace(ctx).Error("failed to insert product",
			zap.Error(err),
			zap.String("spu_id", product.SPUID),
		)
		return err
	}

	return nil
}

func (r *ProductRepository) GetByID(ctx context.Context, spuID string) (*domain.Product, error) {
	ctx, span := otel.Tracer("catalog-product").Start(ctx, "repository.product.get_by_id")
	defer span.End()

	span.SetAttributes(attribute.String("spu_id", spuID))

	var product domain.Product
	err := r.collection.FindOne(ctx, bson.M{"spu_id": spuID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) List(ctx context.Context, filter domain.ProductFilter) (*domain.ProductList, error) {
	ctx, span := otel.Tracer("catalog-product").Start(ctx, "repository.product.list")
	defer span.End()

	query := bson.M{}

	if filter.CategoryID != "" {
		query["category_id"] = filter.CategoryID
	}
	if filter.SellerID != "" {
		query["seller_id"] = filter.SellerID
	}
	if filter.Search != "" {
		query["title"] = bson.M{"$regex": filter.Search, "$options": "i"}
	}
	if filter.MinPrice > 0 || filter.MaxPrice > 0 {
		priceQuery := bson.M{}
		if filter.MinPrice > 0 {
			priceQuery["$gte"] = filter.MinPrice
		}
		if filter.MaxPrice > 0 {
			priceQuery["$lte"] = filter.MaxPrice
		}
		query["skus.price"] = priceQuery
	}

	query["status"] = domain.ProductStatusActive

	skip := int64((filter.Page - 1) * filter.Size)
	limit := int64(filter.Size)

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Size <= 0 || filter.Size > 100 {
		filter.Size = 20
	}

	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.M{"created_at": -1})

	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	if products == nil {
		products = []domain.Product{}
	}

	return &domain.ProductList{
		Products: products,
		Total:    total,
		Page:     filter.Page,
		Size:     filter.Size,
	}, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *domain.Product) error {
	ctx, span := otel.Tracer("catalog-product").Start(ctx, "repository.product.update")
	defer span.End()

	span.SetAttributes(attribute.String("spu_id", product.SPUID))

	product.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"title":       product.Title,
			"description": product.Description,
			"category_id": product.CategoryID,
			"attributes":  product.Attributes,
			"images":      product.Images,
			"updated_at":  product.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"spu_id": product.SPUID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, spuID string) error {
	ctx, span := otel.Tracer("catalog-product").Start(ctx, "repository.product.delete")
	defer span.End()

	span.SetAttributes(attribute.String("spu_id", spuID))

	result, err := r.collection.DeleteOne(ctx, bson.M{"spu_id": spuID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *ProductRepository) GetSKU(ctx context.Context, skuID string) (*domain.SKU, error) {
	ctx, span := otel.Tracer("catalog-product").Start(ctx, "repository.product.get_sku")
	defer span.End()

	span.SetAttributes(attribute.String("sku_id", skuID))

	pipeline := bson.A{
		bson.M{"$unwind": "$skus"},
		bson.M{"$match": bson.M{"skus.sku_id": skuID}},
		bson.M{"$replaceRoot": bson.M{"newRoot": "$skus"}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var skus []domain.SKU
	if err := cursor.All(ctx, &skus); err != nil {
		return nil, err
	}

	if len(skus) == 0 {
		return nil, nil
	}

	return &skus[0], nil
}

func (r *ProductRepository) BatchGetSKUs(ctx context.Context, skuIDs []string) (map[string]*domain.SKU, error) {
	ctx, span := otel.Tracer("catalog-product").Start(ctx, "repository.product.batch_get_skus")
	defer span.End()

	pipeline := bson.A{
		bson.M{"$unwind": "$skus"},
		bson.M{"$match": bson.M{"skus.sku_id": bson.M{"$in": skuIDs}}},
		bson.M{"$replaceRoot": bson.M{"newRoot": "$skus"}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var skus []domain.SKU
	if err := cursor.All(ctx, &skus); err != nil {
		return nil, err
	}

	result := make(map[string]*domain.SKU, len(skus))
	for i := range skus {
		result[skus[i].SKUID] = &skus[i]
	}

	return result, nil
}
