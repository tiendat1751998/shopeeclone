package domain
import "time"

type Livestream struct { ID string `db:"id" json:"id"`; SellerID string `db:"seller_id" json:"seller_id"`; Title string `db:"title" json:"title"`; Status string `db:"status" json:"status"`; ViewerCount int64 `db:"viewer_count" json:"viewer_count"`; PeakViewers int64 `db:"peak_viewers" json:"peak_viewers"`; StartedAt *time.Time `db:"started_at" json:"started_at,omitempty"`; EndedAt *time.Time `db:"ended_at" json:"ended_at,omitempty"`; CreatedAt time.Time `db:"created_at" json:"created_at"` }

type ChatMessage struct { ID string `json:"id"`; RoomID string `json:"room_id"`; UserID string `json:"user_id"`; Content string `json:"content"`; Type string `json:"type"`; Timestamp time.Time `json:"timestamp"` }

type Reaction struct { ID string `json:"id"`; RoomID string `json:"room_id"`; UserID string `json:"user_id"`; Type string `json:"type"`; Timestamp time.Time `json:"timestamp"` }

type Gift struct { ID string `json:"id"`; RoomID string `json:"room_id"`; UserID string `json:"user_id"`; GiftType string `json:"gift_type"`; Amount int64 `json:"amount"`; Timestamp time.Time `json:"timestamp"` }

type PinnedProduct struct { ID string `db:"id" json:"id"`; LivestreamID string `db:"livestream_id" json:"livestream_id"`; ProductID string `db:"product_id" json:"product_id"`; PinnedAt time.Time `db:"pinned_at" json:"pinned_at"` }

const ( LiveStatusScheduled = "scheduled"; LiveStatusLive = "live"; LiveStatusEnded = "ended" )
var ErrLivestreamNotFound = ErrLive("livestream_not_found")
type ErrLive string
func (e ErrLive) Error() string { return "live: " + string(e) }
