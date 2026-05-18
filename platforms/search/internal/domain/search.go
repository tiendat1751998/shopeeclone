package domain
import "time"

type SearchQuery struct { Query string `json:"query"`; CategoryID string `json:"category_id"`; ShopID string `json:"shop_id"`; Filters map[string]interface{} `json:"filters"`; SortBy string `json:"sort_by"`; Page int `json:"page"`; Limit int `json:"limit"` }

type SearchResult struct { Products []ProductHit `json:"products"`; Total int64 `json:"total"`; Page int `json:"page"`; Limit int `json:"limit"`; TookMs int64 `json:"took_ms"`; Suggestions []string `json:"suggestions,omitempty"` }

type ProductHit struct { ID string `json:"id"`; Name string `json:"name"`; ShopID string `json:"shop_id"`; ShopName string `json:"shop_name"`; Price int64 `json:"price"`; ImageURL string `json:"image_url"`; Score float64 `json:"score"`; Highlights map[string]string `json:"highlights,omitempty"` }

type AutocompleteResult struct { Suggestions []Suggestion `json:"suggestions"`; TookMs int64 `json:"took_ms"` }

type Suggestion struct { Text string `json:"text"`; Score float64 `json:"score"`; Type string `json:"type"` }

type IndexDocument struct { ID string `json:"id"`; Name string `json:"name"`; Description string `json:"description"`; CategoryID string `json:"category_id"`; CategoryPath string `json:"category_path"`; ShopID string `json:"shop_id"`; ShopName string `json:"shop_name"`; Price int64 `json:"price"`; Attributes map[string]string `json:"attributes"`; ImageURL string `json:"image_url"`; Status string `json:"status"`; CreatedAt time.Time `json:"created_at"`; UpdatedAt time.Time `json:"updated_at"` }

type RankingConfig struct { RelevanceWeight float64 `json:"relevance_weight"`; PopularityWeight float64 `json:"popularity_weight"`; FreshnessWeight float64 `json:"freshness_weight"`; SponsoredWeight float64 `json:"sponsored_weight"` }

var ErrSearchFailed = ErrSearch("search_failed")
var ErrIndexNotFound = ErrSearch("index_not_found")
type ErrSearch string
func (e ErrSearch) Error() string { return "search: " + string(e) }
