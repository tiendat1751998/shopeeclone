package domain

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type QRCodeType string

const (
	QRCodePickup   QRCodeType = "pickup"
	QRCodeDelivery QRCodeType = "delivery"
)

type QRCodeStatus string

const (
	QRCodeStatusActive   QRCodeStatus = "active"
	QRCodeStatusScanned  QRCodeStatus = "scanned"
	QRCodeStatusExpired  QRCodeStatus = "expired"
	QRCodeStatusRevoked  QRCodeStatus = "revoked"
)

type QRCode struct {
	ID          string       `db:"id" json:"id"`
	ShipmentID  string       `db:"shipment_id" json:"shipment_id"`
	Type        QRCodeType   `db:"type" json:"type"`
	Code        string       `db:"code" json:"code"`
	Status      QRCodeStatus `db:"status" json:"status"`
	SignedToken string       `db:"signed_token" json:"signed_token"`
	ExpiresAt   time.Time    `db:"expires_at" json:"expires_at"`
	ScannedAt   *time.Time   `db:"scanned_at" json:"scanned_at,omitempty"`
	ScannedBy   string       `db:"scanned_by" json:"scanned_by,omitempty"`
	ScanCount   int32        `db:"scan_count" json:"scan_count"`
	CreatedAt   time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at" json:"updated_at"`
}

type ScanEvent struct {
	ID          string    `db:"id" json:"id"`
	QRCodeID    string    `db:"qr_code_id" json:"qr_code_id"`
	ShipmentID  string    `db:"shipment_id" json:"shipment_id"`
	ShipperID   string    `db:"shipper_id" json:"shipper_id"`
	ShipperName string    `db:"shipper_name" json:"shipper_name,omitempty"`
	ShipperRole string    `db:"shipper_role" json:"shipper_role"`
	ScanType    string    `db:"scan_type" json:"scan_type"`
	Latitude    float64   `db:"latitude" json:"latitude"`
	Longitude   float64   `db:"longitude" long:"longitude"`
	DeviceInfo  string    `db:"device_info" json:"device_info,omitempty"`
	IPAddress   string    `db:"ip_address" json:"ip_address,omitempty"`
	IsValid     bool      `db:"is_valid" json:"is_valid"`
	FailReason  string    `db:"fail_reason" json:"fail_reason,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

func NewQRCode(shipmentID string, qrType QRCodeType, ttl time.Duration, secret string) *QRCode {
	now := time.Now().UTC()
	code := uuid.New().String()
	return &QRCode{
		ID:          uuid.New().String(),
		ShipmentID:  shipmentID,
		Type:        qrType,
		Code:        code,
		Status:      QRCodeStatusActive,
		SignedToken: SignQRCode(code, shipmentID, secret),
		ExpiresAt:   now.Add(ttl),
		ScanCount:   0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func NewScanEvent(qrCodeID, shipmentID, shipperID, shipperName, shipperRole, scanType string, lat, lng float64, deviceInfo, ip string, isValid bool, failReason string) *ScanEvent {
	return &ScanEvent{
		ID:          uuid.New().String(),
		QRCodeID:    qrCodeID,
		ShipmentID:  shipmentID,
		ShipperID:   shipperID,
		ShipperName: shipperName,
		ShipperRole: shipperRole,
		ScanType:    scanType,
		Latitude:    lat,
		Longitude:   lng,
		DeviceInfo:  deviceInfo,
		IPAddress:   ip,
		IsValid:     isValid,
		FailReason:  failReason,
		CreatedAt:   time.Now().UTC(),
	}
}

func SignQRCode(code, shipmentID, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(code + ":" + shipmentID))
	return hex.EncodeToString(mac.Sum(nil))
}

func (q *QRCode) VerifySignature(secret string) bool {
	expected := SignQRCode(q.Code, q.ShipmentID, secret)
	return hmac.Equal([]byte(q.SignedToken), []byte(expected))
}

func (q *QRCode) IsExpired() bool {
	return time.Now().UTC().After(q.ExpiresAt)
}

func (q *QRCode) IsUsable() bool {
	return q.Status == QRCodeStatusActive && !q.IsExpired()
}

func (q *QRCode) MarkScanned(shipperID string) error {
	if !q.IsUsable() {
		if q.Status == QRCodeStatusScanned {
			return fmt.Errorf("%w: QR code already scanned", ErrQRCodeAlreadyScanned)
		}
		if q.Status == QRCodeStatusRevoked {
			return fmt.Errorf("%w: QR code has been revoked", ErrQRCodeRevoked)
		}
		if q.IsExpired() {
			return fmt.Errorf("%w: QR code has expired", ErrQRCodeExpired)
		}
		return fmt.Errorf("%w: QR code is not active", ErrQRCodeInactive)
	}
	now := time.Now().UTC()
	q.Status = QRCodeStatusScanned
	q.ScannedAt = &now
	q.ScannedBy = shipperID
	q.ScanCount++
	q.UpdatedAt = now
	return nil
}

func (q *QRCode) Revoke() {
	q.Status = QRCodeStatusRevoked
	q.UpdatedAt = time.Now().UTC()
}
