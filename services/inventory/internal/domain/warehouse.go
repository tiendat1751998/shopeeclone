package domain

import "time"

type Warehouse struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Code      string    `db:"code" json:"code"`
	Address   string    `db:"address" json:"address"`
	City      string    `db:"city" json:"city"`
	Region    string    `db:"region" json:"region"`
	Priority  int       `db:"priority" json:"priority"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
