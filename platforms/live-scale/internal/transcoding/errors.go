package transcoding

import "errors"

var (
	ErrJobNotFound      = errors.New("transcode job not found")
	ErrInvalidJobData   = errors.New("invalid transcode job data")
	ErrUnsupportedProfile = errors.New("unsupported video profile")
	ErrJobAlreadyStarted = errors.New("job already started")
)
