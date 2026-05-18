package domain

import (
	"time"

	"github.com/google/uuid"
)

type LifecycleEvent struct {
	ID               string      `db:"id" json:"id"`
	OrderID          string      `db:"order_id" json:"order_id"`
	FromStatus       OrderStatus `db:"from_state" json:"from_status"`
	ToStatus         OrderStatus `db:"to_state" json:"to_status"`
	TransitionReason string      `db:"transition_reason" json:"transition_reason"`
	ActorID          string      `db:"actor_id" json:"actor_id"`
	ActorType        string      `db:"actor_type" json:"actor_type"`
	Metadata         []byte      `db:"metadata" json:"metadata,omitempty"`
	CreatedAt        time.Time   `db:"created_at" json:"created_at"`
}

func NewLifecycleEvent(orderID string, from, to OrderStatus, reason, actorID, actorType string) *LifecycleEvent {
	return &LifecycleEvent{
		ID:               uuid.New().String(),
		OrderID:          orderID,
		FromStatus:       from,
		ToStatus:         to,
		TransitionReason: reason,
		ActorID:          actorID,
		ActorType:        actorType,
		CreatedAt:        time.Now().UTC(),
	}
}
