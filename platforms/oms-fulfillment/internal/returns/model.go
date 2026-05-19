package returns

import "time"

type ReturnStatus string

const (
	ReturnStatusRequested ReturnStatus = "requested"
	ReturnStatusApproved  ReturnStatus = "approved"
	ReturnStatusRejected  ReturnStatus = "rejected"
	ReturnStatusReceived  ReturnStatus = "received"
	ReturnStatusRefunded  ReturnStatus = "refunded"
)

var validReturnTransitions = map[ReturnStatus][]ReturnStatus{
	ReturnStatusRequested: {ReturnStatusApproved, ReturnStatusRejected},
	ReturnStatusApproved:  {ReturnStatusReceived, ReturnStatusRejected},
	ReturnStatusRejected:  {},
	ReturnStatusReceived:  {ReturnStatusRefunded},
	ReturnStatusRefunded:  {},
}

func IsValidReturnTransition(from, to ReturnStatus) bool {
	allowed, ok := validReturnTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

type ReturnItem struct {
	ItemID    string `json:"item_id"`
	Quantity  int    `json:"quantity"`
	Condition string `json:"condition"`
}

type Return struct {
	ID           string       `json:"id"`
	OrderID      string       `json:"order_id"`
	UserID       string       `json:"user_id"`
	Items        []ReturnItem `json:"items"`
	Reason       string       `json:"reason"`
	Status       ReturnStatus `json:"status"`
	RefundAmount float64      `json:"refund_amount"`
	RMNumber     string       `json:"rma_number"`
	CreatedAt    time.Time    `json:"created_at"`
	ApprovedAt   *time.Time   `json:"approved_at,omitempty"`
	ReceivedAt   *time.Time   `json:"received_at,omitempty"`
	RefundedAt   *time.Time   `json:"refunded_at,omitempty"`
}
