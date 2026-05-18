package application

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopee-clone/shopee/packages/go-shared/pkg/errors"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
	"github.com/shopee-clone/shopee/services/product/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

// AttributeRepository defines the interface for attribute persistence.
type AttributeRepository interface {
	Create(ctx context.Context, attr *domain.Attribute) error
	GetByID(ctx context.Context, id int64) (*domain.Attribute, error)
	ListByCategory(ctx context.Context, categoryID int64) ([]domain.Attribute, error)
	Update(ctx context.Context, attr *domain.Attribute) error
	Delete(ctx context.Context, id int64) error
	CreateValue(ctx context.Context, value *domain.AttributeValue) error
	ListValues(ctx context.Context, attributeID int64) ([]domain.AttributeValue, error)
}

// AttributeCache defines the interface for attribute caching.
type AttributeCache interface {
	Get(ctx context.Context, id int64) (*domain.Attribute, error)
	Set(ctx context.Context, id int64, attr *domain.Attribute, ttl time.Duration) error
	Delete(ctx context.Context, id int64) error
	InvalidateByCategory(ctx context.Context, categoryID int64) error
}

// AttributeService handles product attribute use cases.
type AttributeService struct {
	repo      AttributeRepository
	cache     AttributeCache
	publisher EventPublisher
	logger    *zap.Logger
	tracer    otel.Tracer
}

// NewAttributeService creates a new AttributeService.
func NewAttributeService(
	repo AttributeRepository,
	cache AttributeCache,
	publisher EventPublisher,
	logger *zap.Logger,
) *AttributeService {
	return &AttributeService{
		repo:      repo,
		cache:     cache,
		publisher: publisher,
		logger:    logger,
		tracer:    otel.Tracer("attribute-service"),
	}
}

// CreateAttribute creates a new attribute definition for a category.
func (s *AttributeService) CreateAttribute(ctx context.Context, req CreateAttributeRequest) (*AttributeResponse, error) {
	ctx, span := s.tracer.Start(ctx, "application.attribute.create")
	defer span.End()

	log := observability.LogWithTrace(ctx)

	// Validate
	if err := validateCreateAttribute(req); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "validation failed")
		log.Warn("validation failed for CreateAttribute", zap.Error(err))
		return nil, err
	}

	span.SetAttributes(
		attribute.String("attribute.name", req.Name),
		attribute.String("attribute.type", req.Type),
	)

	// Build domain attribute
	now := time.Now().UTC()
	attr := &domain.Attribute{
		Name:         req.Name,
		Type:         domain.AttributeType(req.Type),
		IsRequired:   req.IsRequired,
		IsFilterable: req.IsFilterable,
		IsSearchable: req.IsSearchable,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Domain-level validation
	if err := attr.Validate(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "domain validation failed")
		log.Warn("domain validation failed for CreateAttribute", zap.Error(err))
		return nil, errors.NewValidation(err.Error())
	}

	// Persist
	if err := s.repo.Create(ctx, attr); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create attribute")
		observability.BusinessErrorsTotal.WithLabelValues("attribute", "CREATE_FAILED").Inc()
		log.Error("failed to create attribute",
			zap.Error(err),
			zap.String("name", req.Name),
		)
		return nil, errors.NewInternalError(err)
	}

	span.SetAttributes(attribute.Int64("attribute.id", attr.ID))

	log.Info("attribute created",
		zap.Int64("id", attr.ID),
		zap.String("name", attr.Name),
		zap.String("type", string(attr.Type)),
	)

	return ToAttributeResponse(attr), nil
}

// GetAttribute retrieves a single attribute by ID (cache-through).
func (s *AttributeService) GetAttribute(ctx context.Context, id int64) (*AttributeResponse, error) {
	ctx, span := s.tracer.Start(ctx, "application.attribute.get")
	defer span.End()

	log := observability.LogWithTrace(ctx)

	if id <= 0 {
		return nil, errors.NewValidation("attribute id must be positive")
	}

	span.SetAttributes(attribute.Int64("attribute.id", id))

	// Try cache first
	cached, err := s.cache.Get(ctx, id)
	if err == nil && cached != nil {
		observability.CacheHitsTotal.WithLabelValues("attribute", "single").Inc()
		log.Debug("attribute cache hit", zap.Int64("id", id))
		return ToAttributeResponse(cached), nil
	}
	observability.CacheMissesTotal.WithLabelValues("attribute", "single").Inc()

	// Cache miss – fetch from repo
	attr, err := s.repo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get attribute")
		log.Error("failed to get attribute",
			zap.Error(err),
			zap.Int64("id", id),
		)
		return nil, errors.NewInternalError(err)
	}

	if attr == nil {
		span.SetStatus(codes.Error, "attribute not found")
		return nil, errors.NewNotFound(fmt.Sprintf("attribute %d not found", id))
	}

	// Populate cache
	if err := s.cache.Set(ctx, id, attr, 30*time.Minute); err != nil {
		log.Warn("failed to cache attribute",
			zap.Error(err),
			zap.Int64("id", id),
		)
	}

	return ToAttributeResponse(attr), nil
}

// ListAttributesByCategory returns all attributes for a given category.
func (s *AttributeService) ListAttributesByCategory(ctx context.Context, categoryID int64) ([]AttributeResponse, error) {
	ctx, span := s.tracer.Start(ctx, "application.attribute.list_by_category")
	defer span.End()

	log := observability.LogWithTrace(ctx)

	if categoryID <= 0 {
		return nil, errors.NewValidation("category_id must be positive")
	}

	span.SetAttributes(attribute.Int64("category.id", categoryID))

	attrs, err := s.repo.ListByCategory(ctx, categoryID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to list attributes by category")
		observability.BusinessErrorsTotal.WithLabelValues("attribute", "LIST_BY_CATEGORY_FAILED").Inc()
		log.Error("failed to list attributes by category",
			zap.Error(err),
			zap.Int64("category_id", categoryID),
		)
		return nil, errors.NewInternalError(err)
	}

	responses := make([]AttributeResponse, 0, len(attrs))
	for i := range attrs {
		responses = append(responses, *ToAttributeResponse(&attrs[i]))
	}

	span.SetAttributes(attribute.Int("attribute.count", len(responses)))

	return responses, nil
}

// UpdateAttribute modifies an existing attribute definition.
func (s *AttributeService) UpdateAttribute(ctx context.Context, id int64, req UpdateAttributeRequest) (*AttributeResponse, error) {
	ctx, span := s.tracer.Start(ctx, "application.attribute.update")
	defer span.End()

	log := observability.LogWithTrace(ctx)

	if id <= 0 {
		return nil, errors.NewValidation("attribute id must be positive")
	}

	span.SetAttributes(attribute.Int64("attribute.id", id))

	// Fetch existing
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to fetch attribute for update")
		log.Error("failed to fetch attribute for update",
			zap.Error(err),
			zap.Int64("id", id),
		)
		return nil, errors.NewInternalError(err)
	}
	if existing == nil {
		span.SetStatus(codes.Error, "attribute not found")
		return nil, errors.NewNotFound(fmt.Sprintf("attribute %d not found", id))
	}

	// Apply updates
	changedFields := make([]string, 0)

	if req.Name != "" && req.Name != existing.Name {
		existing.Name = req.Name
		changedFields = append(changedFields, "name")
	}
	if req.Type != "" && req.Type != string(existing.Type) {
		existing.Type = domain.AttributeType(req.Type)
		changedFields = append(changedFields, "type")
	}
	if req.IsRequired != existing.IsRequired {
		existing.IsRequired = req.IsRequired
		changedFields = append(changedFields, "is_required")
	}
	if req.IsFilterable != existing.IsFilterable {
		existing.IsFilterable = req.IsFilterable
		changedFields = append(changedFields, "is_filterable")
	}
	if req.IsSearchable != existing.IsSearchable {
		existing.IsSearchable = req.IsSearchable
		changedFields = append(changedFields, "is_searchable")
	}

	existing.UpdatedAt = time.Now().UTC()

	// Domain-level validation
	if err := existing.Validate(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "domain validation failed")
		return nil, errors.NewValidation(err.Error())
	}

	// Persist
	if err := s.repo.Update(ctx, existing); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to update attribute")
		observability.BusinessErrorsTotal.WithLabelValues("attribute", "UPDATE_FAILED").Inc()
		log.Error("failed to update attribute",
			zap.Error(err),
			zap.Int64("id", id),
		)
		return nil, errors.NewInternalError(err)
	}

	// Invalidate cache
	if err := s.cache.Delete(ctx, id); err != nil {
		log.Warn("failed to invalidate attribute cache",
			zap.Error(err),
			zap.Int64("id", id),
		)
	}
	if existing.CategoryID > 0 {
		if err := s.cache.InvalidateByCategory(ctx, existing.CategoryID); err != nil {
			log.Warn("failed to invalidate category attribute cache",
				zap.Error(err),
				zap.Int64("category_id", existing.CategoryID),
			)
		}
	}

	log.Info("attribute updated",
		zap.Int64("id", existing.ID),
		zap.Strings("changed_fields", changedFields),
	)

	return ToAttributeResponse(existing), nil
}

// DeleteAttribute removes an attribute definition.
func (s *AttributeService) DeleteAttribute(ctx context.Context, id int64) error {
	ctx, span := s.tracer.Start(ctx, "application.attribute.delete")
	defer span.End()

	log := observability.LogWithTrace(ctx)

	if id <= 0 {
		return errors.NewValidation("attribute id must be positive")
	}

	span.SetAttributes(attribute.Int64("attribute.id", id))

	// Check existence
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to fetch attribute for deletion")
		log.Error("failed to fetch attribute for deletion",
			zap.Error(err),
			zap.Int64("id", id),
		)
		return errors.NewInternalError(err)
	}
	if existing == nil {
		span.SetStatus(codes.Error, "attribute not found")
		return errors.NewNotFound(fmt.Sprintf("attribute %d not found", id))
	}

	// Delete
	if err := s.repo.Delete(ctx, id); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to delete attribute")
		observability.BusinessErrorsTotal.WithLabelValues("attribute", "DELETE_FAILED").Inc()
		log.Error("failed to delete attribute",
			zap.Error(err),
			zap.Int64("id", id),
		)
		return errors.NewInternalError(err)
	}

	// Invalidate cache
	if err := s.cache.Delete(ctx, id); err != nil {
		log.Warn("failed to invalidate attribute cache",
			zap.Error(err),
			zap.Int64("id", id),
		)
	}
	if existing.CategoryID > 0 {
		if err := s.cache.InvalidateByCategory(ctx, existing.CategoryID); err != nil {
			log.Warn("failed to invalidate category attribute cache",
				zap.Error(err),
				zap.Int64("category_id", existing.CategoryID),
			)
		}
	}

	log.Info("attribute deleted", zap.Int64("id", id))

	return nil
}

// CreateAttributeValue creates a new predefined value for an attribute.
func (s *AttributeService) CreateAttributeValue(ctx context.Context, req CreateAttributeValueRequest) (*AttributeValueResponse, error) {
	ctx, span := s.tracer.Start(ctx, "application.attribute.create_value")
	defer span.End()

	log := observability.LogWithTrace(ctx)

	// Validate
	if req.AttributeID <= 0 {
		return nil, errors.NewValidation("attribute_id must be positive")
	}
	if req.Value == "" {
		return nil, errors.NewValidation("value is required")
	}

	span.SetAttributes(
		attribute.Int64("attribute.id", req.AttributeID),
		attribute.String("attribute.value", req.Value),
	)

	// Verify attribute exists
	attr, err := s.repo.GetByID(ctx, req.AttributeID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to fetch parent attribute")
		log.Error("failed to fetch parent attribute for value creation",
			zap.Error(err),
			zap.Int64("attribute_id", req.AttributeID),
		)
		return nil, errors.NewInternalError(err)
	}
	if attr == nil {
		span.SetStatus(codes.Error, "parent attribute not found")
		return nil, errors.NewNotFound(fmt.Sprintf("attribute %d not found", req.AttributeID))
	}

	// Build value
	now := time.Now().UTC()
	value := &domain.AttributeValue{
		AttributeID:  req.AttributeID,
		Value:        req.Value,
		DisplayValue: req.DisplayValue,
		SortOrder:    req.SortOrder,
	}

	// Persist
	if err := s.repo.CreateValue(ctx, value); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create attribute value")
		observability.BusinessErrorsTotal.WithLabelValues("attribute", "CREATE_VALUE_FAILED").Inc()
		log.Error("failed to create attribute value",
			zap.Error(err),
			zap.Int64("attribute_id", req.AttributeID),
		)
		return nil, errors.NewInternalError(err)
	}

	span.SetAttributes(attribute.Int64("attribute_value.id", value.ID))

	// Invalidate attribute cache to reflect new value
	if err := s.cache.Delete(ctx, req.AttributeID); err != nil {
		log.Warn("failed to invalidate attribute cache after value creation",
			zap.Error(err),
			zap.Int64("attribute_id", req.AttributeID),
		)
	}

	log.Info("attribute value created",
		zap.Int64("id", value.ID),
		zap.Int64("attribute_id", req.AttributeID),
		zap.String("value", req.Value),
	)

	return &AttributeValueResponse{
		ID:           value.ID,
		AttributeID:  value.AttributeID,
		Value:        value.Value,
		DisplayValue: value.DisplayValue,
		SortOrder:    value.SortOrder,
	}, nil
}

// ListAttributeValues returns all predefined values for an attribute.
func (s *AttributeService) ListAttributeValues(ctx context.Context, attributeID int64) ([]AttributeValueResponse, error) {
	ctx, span := s.tracer.Start(ctx, "application.attribute.list_values")
	defer span.End()

	log := observability.LogWithTrace(ctx)

	if attributeID <= 0 {
		return nil, errors.NewValidation("attribute_id must be positive")
	}

	span.SetAttributes(attribute.Int64("attribute.id", attributeID))

	values, err := s.repo.ListValues(ctx, attributeID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to list attribute values")
		observability.BusinessErrorsTotal.WithLabelValues("attribute", "LIST_VALUES_FAILED").Inc()
		log.Error("failed to list attribute values",
			zap.Error(err),
			zap.Int64("attribute_id", attributeID),
		)
		return nil, errors.NewInternalError(err)
	}

	responses := make([]AttributeValueResponse, 0, len(values))
	for i := range values {
		responses = append(responses, AttributeValueResponse{
			ID:           values[i].ID,
			AttributeID:  values[i].AttributeID,
			Value:        values[i].Value,
			DisplayValue: values[i].DisplayValue,
			SortOrder:    values[i].SortOrder,
		})
	}

	span.SetAttributes(attribute.Int("attribute_value.count", len(responses)))

	return responses, nil
}

// -----------------------------------------------------------------------------
// Validation Helpers
// -----------------------------------------------------------------------------

func validateCreateAttribute(req CreateAttributeRequest) error {
	if req.Name == "" {
		return errors.NewValidation("attribute name is required")
	}
	validTypes := map[string]bool{
		"text": true, "number": true, "boolean": true,
		"select": true, "multi_select": true, "color": true,
	}
	if !validTypes[req.Type] {
		return errors.NewValidation(fmt.Sprintf("invalid attribute type: %s", req.Type))
	}
	return nil
}
