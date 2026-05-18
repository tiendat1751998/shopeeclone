package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type StateMachine struct {
	transitions map[OrderStatus]map[OrderStatus]bool
}

func NewStateMachine() *StateMachine {
	sm := &StateMachine{
		transitions: make(map[OrderStatus]map[OrderStatus]bool),
	}
	sm.register(OrderStatusPending, OrderStatusAwaitingPayment, OrderStatusCancelled)
	sm.register(OrderStatusAwaitingPayment, OrderStatusPaid, OrderStatusCancelled)
	sm.register(OrderStatusPaid, OrderStatusProcessing, OrderStatusCancelled)
	sm.register(OrderStatusProcessing, OrderStatusPacked, OrderStatusCancelled)
	sm.register(OrderStatusPacked, OrderStatusShipped, OrderStatusCancelled)
	sm.register(OrderStatusShipped, OrderStatusDelivered, OrderStatusRefunded)
	sm.register(OrderStatusDelivered, OrderStatusCompleted, OrderStatusRefunded)
	sm.register(OrderStatusCompleted, OrderStatusRefunded)
	return sm
}

func (sm *StateMachine) register(from OrderStatus, to ...OrderStatus) {
	if sm.transitions[from] == nil {
		sm.transitions[from] = make(map[OrderStatus]bool)
	}
	for _, t := range to {
		sm.transitions[from][t] = true
	}
}

func (sm *StateMachine) CanTransition(from, to OrderStatus) bool {
	if targets, ok := sm.transitions[from]; ok {
		return targets[to]
	}
	return false
}

func (sm *StateMachine) GetValidTransitions(from OrderStatus) []OrderStatus {
	targets, ok := sm.transitions[from]
	if !ok {
		return nil
	}
	result := make([]OrderStatus, 0, len(targets))
	for t := range targets {
		result = append(result, t)
	}
	return result
}

func (sm *StateMachine) Transition(order *Order, target OrderStatus, actorID, actorType, reason string) (*LifecycleEvent, error) {
	if !sm.CanTransition(order.Status, target) {
		return nil, fmt.Errorf("%w: %s -> %s", ErrInvalidStateTransition, order.Status, target)
	}
	now := time.Now().UTC()
	fromStatus := order.Status
	order.Status = target
	order.Version++
	order.UpdatedAt = now
	return &LifecycleEvent{
		ID:               uuid.New().String(),
		OrderID:          order.ID,
		FromStatus:       fromStatus,
		ToStatus:         target,
		TransitionReason: reason,
		ActorID:          actorID,
		ActorType:        actorType,
		CreatedAt:        now,
	}, nil
}
