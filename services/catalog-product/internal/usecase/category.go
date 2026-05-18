package usecase

import (
	"context"

	"github.com/shopee-clone/shopee/packages/go-shared/pkg/errors"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/kafka"
	"github.com/shopee-clone/shopee/services/catalog-product/internal/domain"
	"go.opentelemetry.io/otel"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) error
	GetByID(ctx context.Context, categoryID string) (*domain.Category, error)
	List(ctx context.Context, parentID string, level int32) ([]domain.Category, error)
	Update(ctx context.Context, category *domain.Category) error
}

type CategoryUseCase struct {
	repo     CategoryRepository
	producer *kafka.Producer
}

func NewCategoryUseCase(repo CategoryRepository, producer *kafka.Producer) *CategoryUseCase {
	return &CategoryUseCase{
		repo:     repo,
		producer: producer,
	}
}

func (uc *CategoryUseCase) Create(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	ctx, span := otel.Tracer("catalog-product").Start(ctx, "usecase.category.create")
	defer span.End()

	if category.Name == "" {
		return nil, errors.NewValidation("category name is required")
	}

	if err := uc.repo.Create(ctx, category); err != nil {
		return nil, errors.NewInternalError(err)
	}

	return category, nil
}

func (uc *CategoryUseCase) GetByID(ctx context.Context, categoryID string) (*domain.Category, error) {
	ctx, span := otel.Tracer("catalog-product").Start(ctx, "usecase.category.get_by_id")
	defer span.End()

	category, err := uc.repo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	if category == nil {
		return nil, domain.ErrCategoryNotFound
	}

	return category, nil
}

func (uc *CategoryUseCase) List(ctx context.Context, parentID string, level int32) ([]domain.Category, error) {
	ctx, span := otel.Tracer("catalog-product").Start(ctx, "usecase.category.list")
	defer span.End()

	categories, err := uc.repo.List(ctx, parentID, level)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	return categories, nil
}
