package sms

import "time"

type SMSProvider string

const (
	ProviderTwilio SMSProvider = "twilio"
	ProviderMock   SMSProvider = "mock"
)

type SMSMessage struct {
	ID        string       `json:"id"`
	To        string       `json:"to"`
	From      string       `json:"from"`
	Body      string       `json:"body"`
	Status    string       `json:"status"`
	Provider  SMSProvider  `json:"provider"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type SendSMSRequest struct {
	To   string `json:"to"`
	Body string `json:"body"`
}

type BulkSMSRequest struct {
	To   []string `json:"to"`
	Body string   `json:"body"`
}

type VerifyPhoneRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type VerifyPhoneResponse struct {
	Valid    bool   `json:"valid"`
	Message  string `json:"message"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}
