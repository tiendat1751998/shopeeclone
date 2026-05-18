package domain
import "time"

type Notification struct { ID string `db:"id" json:"id"`; UserID string `db:"user_id" json:"user_id"`; Type string `db:"type" json:"type"`; Title string `db:"title" json:"title"`; Body string `db:"body" json:"body"`; Data string `db:"data" json:"data,omitempty"`; Channel string `db:"channel" json:"channel"`; Status string `db:"status" json:"status"`; Priority int `db:"priority" json:"priority"`; ReadAt *time.Time `db:"read_at" json:"read_at,omitempty"`; CreatedAt time.Time `db:"created_at" json:"created_at"` }

const (
	NotifyTypePush   = "push"
	NotifyTypeEmail  = "email"
	NotifyTypeSMS    = "sms"
	NotifyTypeInApp  = "inapp"
	NotifyStatusPending   = "pending"
	NotifyStatusSent      = "sent"
	NotifyStatusDelivered = "delivered"
	NotifyStatusFailed    = "failed"
	NotifyStatusRead      = "read"
)

type NotificationTemplate struct { ID string `db:"id" json:"id"`; Name string `db:"name" json:"name"`; Type string `db:"type" json:"type"`; Subject string `db:"subject" json:"subject"`; Body string `db:"body" json:"body"`; Variables string `db:"variables" json:"variables,omitempty"`; Version int `db:"version" json:"version"`; IsActive bool `db:"is_active" json:"is_active"`; CreatedAt time.Time `db:"created_at" json:"created_at"` }

type UserPreference struct { UserID string `db:"user_id" json:"user_id"`; Channel string `db:"channel" json:"channel"`; Enabled bool `db:"enabled" json:"enabled"`; QuietHours string `db:"quiet_hours" json:"quiet_hours,omitempty"` }

type DeliveryLog struct { ID string `db:"id" json:"id"`; NotificationID string `db:"notification_id" json:"notification_id"`; Channel string `db:"channel" json:"channel"`; Provider string `db:"provider" json:"provider"`; Status string `db:"status" json:"status"`; ErrorMsg string `db:"error_message" json:"error_message,omitempty"`; AttemptCount int `db:"attempt_count" json:"attempt_count"`; CreatedAt time.Time `db:"created_at" json:"created_at"` }

var ErrNotificationNotFound = ErrNotification("notification_not_found")
var ErrTemplateNotFound = ErrNotification("template_not_found")
type ErrNotification string
func (e ErrNotification) Error() string { return "notification: " + string(e) }
