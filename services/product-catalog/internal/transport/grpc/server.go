package grpc

import (
	"context"

	"github.com/shopee-clone/shopee/services/product-catalog/internal/application"
	"github.com/shopee-clone/shopee/services/product-catalog/internal/domain"
	pb "github.com/shopee-clone/shopee/services/product-catalog/proto/productcatalog/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CatalogGRPCServer struct {
	pb.UnimplementedCatalogServiceServer
	catalogService *application.CatalogService
}

func NewCatalogGRPCServer(svc *application.CatalogService) *CatalogGRPCServer {
	return &CatalogGRPCServer{catalogService: svc}
}

func (s *CatalogGRPCServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	r := &application.CreateProductRequest{
		ShopID: req.ShopId, Name: req.Name, Description: req.Description,
		CategoryID: req.CategoryId, Brand: req.Brand, Condition: req.Condition,
		IdempotencyKey: req.IdempotencyKey,
	}
	product, err := s.catalogService.CreateProduct(ctx, r)
	if err != nil { return nil, toGRPCError(err) }
	return toProductResponse(product), nil
}

func (s *CatalogGRPCServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	product, err := s.catalogService.GetProduct(ctx, req.ProductId)
	if err != nil { return nil, toGRPCError(err) }
	return toProductResponse(product), nil
}

func (s *CatalogGRPCServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, total, err := s.catalogService.ListProducts(ctx, req.ShopId, req.Status, int(req.Page), int(req.PageSize))
	if err != nil { return nil, toGRPCError(err) }
	pbProducts := make([]*pb.ProductResponse, 0, len(products))
	for _, p := range products { pbProducts = append(pbProducts, toProductResponse(p)) }
	return &pb.ListProductsResponse{Products: pbProducts, Total: int32(total), Page: req.Page}, nil
}

func toProductResponse(p *domain.Product) *pb.ProductResponse {
	return &pb.ProductResponse{
		Id: p.ID, ShopId: p.ShopID, Name: p.Name, Description: p.Description,
		CategoryId: p.CategoryID, Brand: p.Brand, Status: string(p.Status),
		CreatedAt: p.CreatedAt.String(),
	}
}

func toGRPCError(err error) error {
	switch err {
	case domain.ErrProductNotFound: return status.Error(codes.NotFound, err.Error())
	case domain.ErrCategoryNotFound: return status.Error(codes.NotFound, err.Error())
	case domain.ErrUnauthorized: return status.Error(codes.PermissionDenied, err.Error())
	case domain.ErrInvalidProductData: return status.Error(codes.InvalidArgument, err.Error())
	default:
		zap.L().Error("unexpected gRPC error", zap.Error(err))
		return status.Error(codes.Internal, "internal server error")
	}
}
