package collabvector

import "errors"

var (
	ErrNotEnoughInteractions = errors.New("collabvector: not enough interactions for factorization")
	ErrUserNotFound          = errors.New("collabvector: user not found in matrix")
)
