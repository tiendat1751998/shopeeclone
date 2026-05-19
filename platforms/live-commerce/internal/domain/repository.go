package domain

import "context"

type LivestreamRepository interface {
	Create(ctx context.Context, ls *Livestream) error
	GetByID(ctx context.Context, id string) (*Livestream, error)
	Update(ctx context.Context, ls *Livestream) error
	ListBySeller(ctx context.Context, sellerID string, offset, limit int) ([]*Livestream, int64, error)
	ListActive(ctx context.Context, offset, limit int) ([]*Livestream, int64, error)
}

type ChatMessageRepository interface {
	Save(ctx context.Context, msg *ChatMessage) error
	GetByRoom(ctx context.Context, roomID string, offset, limit int) ([]*ChatMessage, int64, error)
	GetLastSequence(ctx context.Context, roomID string) (int64, error)
	MarkModerated(ctx context.Context, messageID string) error
	Delete(ctx context.Context, messageID string) error
}

type ReactionRepository interface {
	Save(ctx context.Context, reaction *Reaction) error
	GetCountByRoom(ctx context.Context, roomID string, reactionType string) (int64, error)
	GetSummaryByRoom(ctx context.Context, roomID string) (map[string]int64, error)
}

type GiftRepository interface {
	Save(ctx context.Context, gift *Gift) error
	GetLeaderboardByRoom(ctx context.Context, roomID string, limit int) ([]*GiftLeaderboardEntry, error)
}

type PinnedProductRepository interface {
	Pin(ctx context.Context, pp *PinnedProduct) error
	Unpin(ctx context.Context, livestreamID, productID string) error
	GetActiveByLivestream(ctx context.Context, livestreamID string) ([]*PinnedProduct, error)
}

type ModerationRepository interface {
	SaveAction(ctx context.Context, action *ModerationAction) error
	IsUserMuted(ctx context.Context, roomID, userID string) (bool, error)
	GetMuteDuration(ctx context.Context, roomID, userID string) (int64, error)
}

type GiftLeaderboardEntry struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Total    int64  `json:"total"`
	Rank     int    `json:"rank"`
}
