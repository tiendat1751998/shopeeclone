package validation

import "errors"

var (
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong  = errors.New("password must not exceed 128 characters")
	ErrPasswordNoUpper  = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLower  = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoDigit  = errors.New("password must contain at least one digit")
	ErrPasswordNoSpecial = errors.New("password must contain at least one special character")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrInvalidUsername  = errors.New("username must be 3-32 characters, alphanumeric and underscore only")
	ErrInvalidPhone     = errors.New("invalid phone number format")
)
