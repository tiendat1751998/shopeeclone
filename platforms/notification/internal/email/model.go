package email

import "time"

type EmailStatus string

const (
	EmailStatusPending   EmailStatus = "pending"
	EmailStatusSent      EmailStatus = "sent"
	EmailStatusDelivered EmailStatus = "delivered"
	EmailStatusBounced   EmailStatus = "bounced"
	EmailStatusOpened    EmailStatus = "opened"
	EmailStatusFailed    EmailStatus = "failed"
)

type EmailAttachment struct {
	Filename string `json:"filename"`
	Content  []byte `json:"content"`
	MimeType string `json:"mime_type"`
}

type EmailMessage struct {
	ID        string           `json:"id"`
	To        []string         `json:"to"`
	CC        []string         `json:"cc,omitempty"`
	BCC       []string         `json:"bcc,omitempty"`
	From      string           `json:"from"`
	ReplyTo   string           `json:"reply_to,omitempty"`
	Subject   string           `json:"subject"`
	PlainText string           `json:"plain_text"`
	HTML      string           `json:"html"`
	Status    EmailStatus      `json:"status"`
	Attachments []EmailAttachment `json:"attachments,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type EmailProvider string

const (
	ProviderSMTP     EmailProvider = "smtp"
	ProviderSendGrid EmailProvider = "sendgrid"
	ProviderMailgun  EmailProvider = "mailgun"
)

type SendEmailRequest struct {
	To          []string          `json:"to"`
	CC          []string          `json:"cc,omitempty"`
	BCC         []string          `json:"bcc,omitempty"`
	Subject     string            `json:"subject"`
	PlainText   string            `json:"plain_text"`
	HTML        string            `json:"html"`
	ReplyTo     string            `json:"reply_to,omitempty"`
	Attachments []EmailAttachment `json:"attachments,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	TemplateID  string            `json:"template_id,omitempty"`
	TemplateVars map[string]interface{} `json:"template_vars,omitempty"`
}
