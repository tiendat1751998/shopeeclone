package domain
import "context"

type NotificationRepository interface {
	Create(ctx context.Context, n *Notification) error
	FindByID(ctx context.Context, id string) (*Notification, error)
	FindByUserID(ctx context.Context, userID string, offset, limit int) ([]*Notification, int64, error)
	MarkRead(ctx context.Context, id string) error
	MarkDelivered(ctx context.Context, id string) error
	MarkFailed(ctx context.Context, id, reason string) error
}

type TemplateRepository interface {
	FindByName(ctx context.Context, name string) (*NotificationTemplate, error)
	Create(ctx context.Context, t *NotificationTemplate) error
}

type PreferenceRepository interface {
	GetUserPreferences(ctx context.Context, userID string) ([]*UserPreference, error)
	UpdatePreference(ctx context.Context, p *UserPreference) error
}

type DeliveryLogRepository interface {
	Create(ctx context.Context, log *DeliveryLog) error
	FindByNotificationID(ctx context.Context, notificationID string) ([]*DeliveryLog, error)
}
