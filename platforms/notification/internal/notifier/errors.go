package notifier

import "errors"

var (
	ErrNotFound         = errors.New("notification not found")
	ErrInvalidType      = errors.New("invalid notification type")
	ErrInvalidChannel   = errors.New("invalid channel")
	ErrSendFailed       = errors.New("send failed")
	ErrPreferenceBlock  = errors.New("notification blocked by user preferences")
	ErrQuietHours       = errors.New("notification blocked by quiet hours")
	ErrRateLimited      = errors.New("rate limit exceeded")
	ErrTemplateNotFound = errors.New("template not found")
	ErrTemplateRender   = errors.New("template render failed")
	ErrDeviceNotFound   = errors.New("device not found")
)
