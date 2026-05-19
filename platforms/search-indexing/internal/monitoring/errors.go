package monitoring

import "errors"

var (
	ErrIndexMetricNotFound = errors.New("index metric not found")
	ErrMetricAlreadyExists = errors.New("index metric already exists")
)
