package domain

import (
	"time"

	"github.com/google/uuid"
)

type FraudCheckResult struct {
	ID          string    `db:"id" json:"id"`
	PaymentID   string    `db:"payment_id" json:"payment_id"`
	UserID      string    `db:"user_id" json:"user_id"`
	RiskScore   int       `db:"risk_score" json:"risk_score"`
	RiskLevel   string    `db:"risk_level" json:"risk_level"`
	IsFraud     bool      `db:"is_fraud" json:"is_fraud"`
	Reasons     []byte    `db:"reasons" json:"reasons,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

func NewFraudCheckResult(paymentID, userID string, riskScore int, isFraud bool) *FraudCheckResult {
	level := "low"
	if riskScore >= 80 { level = "high" } else if riskScore >= 50 { level = "medium" }
	return &FraudCheckResult{
		ID:        uuid.New().String(),
		PaymentID: paymentID,
		UserID:    userID,
		RiskScore: riskScore,
		RiskLevel: level,
		IsFraud:   isFraud,
		CreatedAt: time.Now().UTC(),
	}
}
