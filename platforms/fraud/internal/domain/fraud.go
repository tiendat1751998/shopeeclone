package domain
import "time"

type FraudScore struct { ID string `json:"id"`; UserID string `json:"user_id"`; OrderID string `json:"order_id,omitempty"`; Score float64 `json:"score"`; RiskLevel string `json:"risk_level"`; Signals []string `json:"signals"`; CreatedAt time.Time `json:"created_at"` }

type FraudRule struct { ID string `db:"id" json:"id"`; Name string `db:"name" json:"name"`; RuleType string `db:"rule_type" json:"rule_type"`; Condition string `db:"condition" json:"condition"`; Action string `db:"action" json:"action"`; Priority int `db:"priority" json:"priority"`; IsActive bool `db:"is_active" json:"is_active"`; CreatedAt time.Time `db:"created_at" json:"created_at"` }

type FraudCase struct { ID string `db:"id" json:"id"`; UserID string `db:"user_id" json:"user_id"`; OrderID string `db:"order_id" json:"order_id,omitempty"`; Status string `db:"status" json:"status"`; Score float64 `db:"score" json:"score"`; Evidence string `db:"evidence" json:"evidence,omitempty"`; AssignedTo string `db:"assigned_to" json:"assigned_to,omitempty"`; CreatedAt time.Time `db:"created_at" json:"created_at"`; UpdatedAt time.Time `db:"updated_at" json:"updated_at"` }

const ( RiskLow = "low"; RiskMedium = "medium"; RiskHigh = "high"; RiskCritical = "critical" )
const ( CaseStatusOpen = "open"; CaseStatusReviewing = "reviewing"; CaseStatusConfirmed = "confirmed"; CaseStatusDismissed = "dismissed" )
var ErrFraudDetection = ErrFraud("fraud_detection_failed")
type ErrFraud string
func (e ErrFraud) Error() string { return "fraud: " + string(e) }
