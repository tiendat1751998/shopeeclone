package featurestore

import "time"

type ValueType string

const (
	TypeString  ValueType = "string"
	TypeNumber  ValueType = "number"
	TypeBoolean ValueType = "boolean"
	TypeVector  ValueType = "vector"
)

type EntityType string

const (
	EntityUser    EntityType = "user"
	EntityProduct EntityType = "product"
	EntityOrder   EntityType = "order"
)

type Feature struct {
	Name        string     `json:"name"`
	ValueType   ValueType  `json:"value_type"`
	Entity      EntityType `json:"entity"`
	Source      string     `json:"source"`
	Description string     `json:"description"`
	IsOnline    bool       `json:"is_online"`
	CreatedAt   time.Time  `json:"created_at"`
}

type FeatureValue struct {
	FeatureName string      `json:"feature_name"`
	EntityID    string      `json:"entity_id"`
	Value       interface{} `json:"value"`
	Timestamp   time.Time   `json:"timestamp"`
	Version     int64       `json:"version"`
}
