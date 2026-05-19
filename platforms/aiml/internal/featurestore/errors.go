package featurestore

import "errors"

var (
	ErrFeatureNotFound      = errors.New("featurestore: feature not found")
	ErrFeatureAlreadyExists = errors.New("featurestore: feature already exists")
	ErrFeatureValueNotFound = errors.New("featurestore: feature value not found")
	ErrInvalidValueType     = errors.New("featurestore: invalid value type")
)
