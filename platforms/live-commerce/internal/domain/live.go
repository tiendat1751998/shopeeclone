package domain

import (
	"fmt"
	"time"
	"github.com/google/uuid"
)

type Livestream struct {
	ID          string     `db:"id" json:"id"`
	SellerID    string     `db:"seller_id" json:"seller_id"`
	Title       string     `db:"title" json:"title"`
	Description string     `db:"description" json:"description,omitempty"`
	CoverURL    string     `db:"cover_url" json:"cover_url,omitempty"`
	Status      string     `db:"status" json:"status"`
	ViewerCount int64      `db:"viewer_count" json:"viewer_count"`
	PeakViewers int64      `db:"peak_viewers" json:"peak_viewers"`
	TotalLikes  int64      `db:"total_likes" json:"total_likes"`
	TotalGifts  int64      `db:"total_gifts" json:"total_gifts"`
	TotalShares int64      `db:"total_shares" json:"total_shares"`
	Category    string     `db:"category" json:"category,omitempty"`
	Tags        []string   `db:"tags" json:"tags,omitempty"`
	StartedAt   *time.Time `db:"started_at" json:"started_at,omitempty"`
	EndedAt     *time.Time `db:"ended_at" json:"ended_at,omitempty"`
	ScheduledAt *time.Time `db:"scheduled_at" json:"scheduled_at,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

const (
	LiveStatusScheduled = "scheduled"
	LiveStatusLive      = "live"
	LiveStatusEnded     = "ended"
	LiveStatusCancelled = "cancelled"
)

func NewLivestream(sellerID, title, description, coverURL, category string, tags []string, scheduledAt *time.Time) *Livestream {
	now := time.Now()
	if tags == nil {
		tags = []string{}
	}
	return &Livestream{
		ID:          uuid.New().String(),
		SellerID:    sellerID,
		Title:       title,
		Description: description,
		CoverURL:    coverURL,
		Category:    category,
		Tags:        tags,
		Status:      LiveStatusScheduled,
		ScheduledAt: scheduledAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (l *Livestream) Start() error {
	if l.Status != LiveStatusScheduled {
		return fmt.Errorf("%w: cannot start from %s", ErrInvalidLiveState, l.Status)
	}
	now := time.Now()
	l.Status = LiveStatusLive
	l.StartedAt = &now
	l.UpdatedAt = now
	return nil
}

func (l *Livestream) End() error {
	if l.Status != LiveStatusLive {
		return fmt.Errorf("%w: cannot end from %s", ErrInvalidLiveState, l.Status)
	}
	now := time.Now()
	l.Status = LiveStatusEnded
	l.EndedAt = &now
	l.UpdatedAt = now
	return nil
}

func (l *Livestream) Cancel() error {
	if l.Status != LiveStatusScheduled {
		return fmt.Errorf("%w: cannot cancel from %s", ErrInvalidLiveState, l.Status)
	}
	l.Status = LiveStatusCancelled
	l.UpdatedAt = time.Now()
	return nil
}

func (l *Livestream) IsLive() bool { return l.Status == LiveStatusLive }

func (l *Livestream) UpdateViewers(count int64) {
	l.ViewerCount = count
	if count > l.PeakViewers {
		l.PeakViewers = count
	}
	l.UpdatedAt = time.Now()
}

func (l *Livestream) AddLike()      { l.TotalLikes++; l.UpdatedAt = time.Now() }
func (l *Livestream) AddGift(amt int64) { l.TotalGifts += amt; l.UpdatedAt = time.Now() }
func (l *Livestream) AddShare()     { l.TotalShares++; l.UpdatedAt = time.Now() }

type ChatMessage struct {
	ID          string    `json:"id"`
	RoomID      string    `json:"room_id"`
	UserID      string    `json:"user_id"`
	Username    string    `json:"username"`
	Content     string    `json:"content"`
	Type        string    `json:"type"`
	IsModerated bool      `json:"is_moderated"`
	Sequence    int64     `json:"sequence"`
	Timestamp   time.Time `json:"timestamp"`
}

const (
	MsgTypeText   = "text"
	MsgTypeSystem = "system"
	MsgTypeGift   = "gift"
	MsgTypeNotice = "notice"
	MaxMsgLength  = 500
)

func NewChatMessage(roomID, userID, username, content, msgType string) *ChatMessage {
	return &ChatMessage{
		ID:        uuid.New().String(),
		RoomID:    roomID,
		UserID:    userID,
		Username:  username,
		Content:   content,
		Type:      msgType,
		Timestamp: time.Now(),
	}
}

type Reaction struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"room_id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

const (
	ReactionLike  = "like"
	ReactionLove  = "love"
	ReactionWow   = "wow"
	ReactionLaugh = "laugh"
	ReactionSad   = "sad"
	ReactionAngry = "angry"
)

func NewReaction(roomID, userID, reactionType string) *Reaction {
	return &Reaction{
		ID:        uuid.New().String(),
		RoomID:    roomID,
		UserID:    userID,
		Type:      reactionType,
		Timestamp: time.Now(),
	}
}

type Gift struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"room_id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	GiftType  string    `json:"gift_type"`
	Amount    int64     `json:"amount"`
	Currency  string    `json:"currency"`
	Timestamp time.Time `json:"timestamp"`
}

type PinnedProduct struct {
	ID           string    `db:"id" json:"id"`
	LivestreamID string    `db:"livestream_id" json:"livestream_id"`
	ProductID    string    `db:"product_id" json:"product_id"`
	ProductName  string    `db:"product_name" json:"product_name"`
	Price        int64     `db:"price" json:"price"`
	ImageURL     string    `db:"image_url" json:"image_url"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	PinnedAt     time.Time `db:"pinned_at" json:"pinned_at"`
}

type ModerationAction struct {
	ID          string    `json:"id"`
	RoomID      string    `json:"room_id"`
	UserID      string    `json:"user_id"`
	Action      string    `json:"action"`
	Reason      string    `json:"reason"`
	ModeratedBy string    `json:"moderated_by"`
	DurationSec int64     `json:"duration_sec,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

const (
	ModActionMute    = "mute"
	ModActionBan     = "ban"
	ModActionRemove  = "remove_message"
	ModActionWarning = "warning"
)

type Room struct {
	ID           string    `json:"id"`
	LivestreamID string    `json:"livestream_id"`
	Status       string    `json:"status"`
	ViewerCount  int64     `json:"viewer_count"`
	CreatedAt    time.Time `json:"created_at"`
}

type ViewerSession struct {
	UserID      string    `json:"user_id"`
	Username    string    `json:"username"`
	LivestreamID string   `json:"livestream_id"`
	ConnectedAt time.Time `json:"connected_at"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
}

type LivestreamState string

const (
	StateScheduled LivestreamState = "scheduled"
	StateLive     LivestreamState = "live"
	StateEnded    LivestreamState = "ended"
	StateCancelled LivestreamState = "cancelled"
)

var (
	ErrLivestreamNotFound = ErrLive("livestream_not_found")
	ErrInvalidLiveState   = ErrLive("invalid_live_state")
	ErrRoomNotFound       = ErrLive("room_not_found")
	ErrAlreadyEnded       = ErrLive("already_ended")
	ErrNotLive            = ErrLive("not_live")
	ErrProductNotPinned   = ErrLive("product_not_pinned")
	ErrMessageTooLong     = ErrLive("message_too_long")
	ErrUserMuted          = ErrLive("user_muted")
	ErrUserBanned         = ErrLive("user_banned")
	ErrInvalidReaction    = ErrLive("invalid_reaction")
	ErrInvalidGiftType    = ErrLive("invalid_gift_type")
	ErrInsufficientCoins  = ErrLive("insufficient_coins")
	ErrDuplicateReaction  = ErrLive("duplicate_reaction")
)

type ErrLive string
func (e ErrLive) Error() string { return "live: " + string(e) }
