package domain

import (
	"time"

	"github.com/google/uuid"
)

type SellerSplit struct {
	ID            string     `db:"id" json:"id"`
	ParentOrderID string     `db:"parent_order_id" json:"parent_order_id"`
	SellerID      string     `db:"seller_id" json:"seller_id"`
	SubOrderID    string     `db:"sub_order_id" json:"sub_order_id"`
	Status        OrderStatus `db:"status" json:"status"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
}

func NewSellerSplit(parentOrderID, sellerID, subOrderID string) *SellerSplit {
	now := time.Now().UTC()
	return &SellerSplit{
		ID:            uuid.New().String(),
		ParentOrderID: parentOrderID,
		SellerID:      sellerID,
		SubOrderID:    subOrderID,
		Status:        OrderStatusPending,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// SplitOrderBySellers groups items by seller and creates sub-orders
func SplitOrderBySellers(order *Order) map[string][]OrderItem {
	sellerItems := make(map[string][]OrderItem)
	for _, item := range order.Items {
		sellerItems[item.ShopID] = append(sellerItems[item.ShopID], item)
	}
	return sellerItems
}
