package domain
import "time"

type ClickEvent struct { ID string `json:"id"`; UserID string `json:"user_id"`; SessionID string `json:"session_id"`; EventType string `json:"event_type"`; ProductID string `json:"product_id,omitempty"`; PageURL string `json:"page_url"`; Referrer string `json:"referrer,omitempty"`; DeviceType string `json:"device_type"`; IPAddress string `json:"ip_address"`; UserAgent string `json:"user_agent"`; Metadata map[string]interface{} `json:"metadata,omitempty"`; Timestamp time.Time `json:"timestamp"` }

type Session struct { ID string `db:"id" json:"id"`; UserID string `db:"user_id" json:"user_id"`; StartTime time.Time `db:"start_time" json:"start_time"`; EndTime *time.Time `db:"end_time" json:"end_time,omitempty"`; PageViews int `db:"page_views" json:"page_views"`; Events int `db:"events" json:"events"`; DeviceType string `db:"device_type" json:"device_type"`; CreatedAt time.Time `db:"created_at" json:"created_at"` }

type AggregatedMetric struct { MetricName string `db:"metric_name" json:"metric_name"`; Dimension string `db:"dimension" json:"dimension"`; Value int64 `db:"value" json:"value"`; WindowStart time.Time `db:"window_start" json:"window_start"`; WindowEnd time.Time `db:"window_end" json:"window_end"` }

const ( EventPageView = "page_view"; EventProductView = "product_view"; EventClick = "click"; EventAddToCart = "add_to_cart"; EventCheckout = "checkout"; EventSearch = "search"; EventImpression = "impression" )
var ErrEventValidation = ErrBehavior("event_validation_failed")
type ErrBehavior string
func (e ErrBehavior) Error() string { return "behavior: " + string(e) }
