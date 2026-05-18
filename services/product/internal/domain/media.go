package domain

import (
	"time"
)

// MediaType represents the type of media asset.
type MediaType string

const (
	MediaTypeImage    MediaType = "IMAGE"
	MediaTypeVideo    MediaType = "VIDEO"
	MediaTypeDocument MediaType = "DOCUMENT"
)

// Media represents a media asset associated with a product or SKU.
type Media struct {
	ID           string    `db:"id"             json:"id"`
	SPUID        string    `db:"spu_id"         json:"spu_id"`
	SKUID        string    `db:"sku_id"         json:"sku_id,omitempty"`
	Type         MediaType `db:"type"           json:"type"`
	URL          string    `db:"url"            json:"url"`
	ThumbnailURL string    `db:"thumbnail_url"  json:"thumbnail_url,omitempty"`
	AltText      string    `db:"alt_text"       json:"alt_text,omitempty"`
	SortOrder    int       `db:"sort_order"     json:"sort_order"`
	MimeType     string    `db:"mime_type"      json:"mime_type"`
	FileSize     int64     `db:"file_size"      json:"file_size"`
	CreatedAt    time.Time `db:"created_at"     json:"created_at"`
}

// IsImage returns true if the media is an image.
func (m *Media) IsImage() bool {
	return m.Type == MediaTypeImage
}

// IsVideo returns true if the media is a video.
func (m *Media) IsVideo() bool {
	return m.Type == MediaTypeVideo
}

// IsSKULevel returns true if the media is associated with a specific SKU.
func (m *Media) IsSKULevel() bool {
	return m.SKUID != ""
}

// FileSizeInKB returns the file size in kilobytes.
func (m *Media) FileSizeInKB() float64 {
	return float64(m.FileSize) / 1024.0
}

// FileSizeInMB returns the file size in megabytes.
func (m *Media) FileSizeInMB() float64 {
	return float64(m.FileSize) / (1024.0 * 1024.0)
}
