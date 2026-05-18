package domain

import (
	"time"

	"github.com/google/uuid"
)

type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
)

type MediaStatus string

const (
	MediaStatusPending  MediaStatus = "pending"
	MediaStatusActive   MediaStatus = "active"
	MediaStatusFailed   MediaStatus = "failed"
	MediaStatusDeleted  MediaStatus = "deleted"
)

type Media struct {
	ID         string      `db:"id" json:"id"`
	ProductID  string      `db:"product_id" json:"product_id"`
	SKUID      string      `db:"sku_id" json:"sku_id,omitempty"`
	Type       MediaType   `db:"type" json:"type"`
	URL        string      `db:"url" json:"url"`
	ThumbnailURL string    `db:"thumbnail_url" json:"thumbnail_url,omitempty"`
	AltText    string      `db:"alt_text" json:"alt_text,omitempty"`
	SortOrder  int         `db:"sort_order" json:"sort_order"`
	Status     MediaStatus `db:"status" json:"status"`
	Metadata   []byte      `db:"metadata" json:"metadata,omitempty"`
	Version    int         `db:"version" json:"version"`
	CreatedAt  time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time   `db:"updated_at" json:"updated_at"`
}

func NewMedia(productID, mediaType, url string, sortOrder int) *Media {
	now := time.Now().UTC()
	return &Media{
		ID: uuid.New().String(), ProductID: productID, Type: MediaType(mediaType),
		URL: url, SortOrder: sortOrder, Status: MediaStatusPending,
		Version: 1, CreatedAt: now, UpdatedAt: now,
	}
}
