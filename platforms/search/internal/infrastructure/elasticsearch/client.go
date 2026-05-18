package elasticsearch
import ("context"; "encoding/json"; "fmt"; "github.com/shopee-clone/shopee/platforms/search/internal/config"; "github.com/shopee-clone/shopee/platforms/search/internal/domain"; "github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"; "go.uber.org/zap")

const IndexProducts = "products"
const IndexAutocomplete = "autocomplete"

type Client struct { cfg config.ESConfig }

func NewClient(cfg config.ESConfig) *Client { return &Client{cfg: cfg} }

func (c *Client) Search(ctx context.Context, query domain.SearchQuery) (*domain.SearchResult, error) {
	// In production: use elastic/go-elasticsearch client
	// For now, return mock result
	observability.LogWithTrace(ctx).Info("ES search", zap.String("query", query.Query))
	return &domain.SearchResult{Products: []domain.ProductHit{}, Total: 0, Page: query.Page, Limit: query.Limit, TookMs: 5}, nil
}

func (c *Client) Autocomplete(ctx context.Context, prefix string, limit int) (*domain.AutocompleteResult, error) {
	return &domain.AutocompleteResult{Suggestions: []domain.Suggestion{}, TookMs: 2}, nil
}

func (c *Client) Index(ctx context.Context, doc *domain.IndexDocument) error {
	data, _ := json.Marshal(doc)
	observability.LogWithTrace(ctx).Info("ES index", zap.String("id", doc.ID), zap.Int("size", len(data)))
	return nil
}

func (c *Client) BulkIndex(ctx context.Context, docs []*domain.IndexDocument) error {
	for _, doc := range docs { c.Index(ctx, doc) }
	return nil
}

func (c *Client) Delete(ctx context.Context, id string) error {
	observability.LogWithTrace(ctx).Info("ES delete", zap.String("id", id))
	return nil
}

func (c *Client) CreateIndex(ctx context.Context, index string) error {
	observability.LogWithTrace(ctx).Info("ES create index", zap.String("index", index))
	return nil
}

func (c *Client) HealthCheck(ctx context.Context) error {
	// In production: ping ES cluster
	return nil
}
