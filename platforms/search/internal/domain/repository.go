package domain
import "context"

type SearchRepository interface {
	Search(ctx context.Context, query SearchQuery) (*SearchResult, error)
	Autocomplete(ctx context.Context, prefix string, limit int) (*AutocompleteResult, error)
	Index(ctx context.Context, doc *IndexDocument) error
	BulkIndex(ctx context.Context, docs []*IndexDocument) error
	Delete(ctx context.Context, id string) error
}

type RankingRepository interface {
	GetConfig(ctx context.Context) (*RankingConfig, error)
	UpdateConfig(ctx context.Context, config *RankingConfig) error
}

type AnalyticsRepository interface {
	RecordQuery(ctx context.Context, query string, resultCount int64, tookMs int64) error
	RecordClick(ctx context.Context, query string, productID string, position int) error
	GetTrendingQueries(ctx context.Context, limit int) ([]string, error)
	GetZeroResultQueries(ctx context.Context, limit int) ([]string, error)
}
