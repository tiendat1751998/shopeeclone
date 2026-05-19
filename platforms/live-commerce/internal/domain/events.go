package domain

import "time"

const (
	EventLivestreamCreated   = "livestream.created"
	EventLivestreamStarted   = "livestream.started"
	EventLivestreamEnded     = "livestream.ended"
	EventLivestreamCancelled = "livestream.cancelled"
	EventChatMessageSent     = "chat.message.sent"
	EventChatMessageDeleted  = "chat.message.deleted"
	EventReactionAdded       = "reaction.added"
	EventGiftSent            = "gift.sent"
	EventProductPinned       = "product.pinned"
	EventProductUnpinned     = "product.unpinned"
	EventModerationAction    = "moderation.action"
	EventViewerJoined        = "viewer.joined"
	EventViewerLeft          = "viewer.left"
)

type Event struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Source    string      `json:"source"`
	Version   int         `json:"version"`
	Timestamp time.Time   `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}

type ChatMessageSentPayload struct {
	MessageID string `json:"message_id"`
	RoomID    string `json:"room_id"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

type ReactionAddedPayload struct {
	ReactionID string `json:"reaction_id"`
	RoomID     string `json:"room_id"`
	UserID     string `json:"user_id"`
	Type       string `json:"type"`
}

type GiftSentPayload struct {
	GiftID   string `json:"gift_id"`
	RoomID   string `json:"room_id"`
	UserID   string `json:"user_id"`
	GiftType string `json:"gift_type"`
	Amount   int64  `json:"amount"`
}

type LivestreamEventPayload struct {
	LivestreamID string `json:"livestream_id"`
	SellerID     string `json:"seller_id"`
	Title        string `json:"title,omitempty"`
	Timestamp    int64  `json:"timestamp"`
}

type ModerationActionPayload struct {
	ActionID    string `json:"action_id"`
	RoomID      string `json:"room_id"`
	UserID      string `json:"user_id"`
	Action      string `json:"action"`
	Reason      string `json:"reason"`
	ModeratedBy string `json:"moderated_by"`
}
