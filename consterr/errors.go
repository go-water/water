package consterr

import (
	"errors"
)

var (
	ErrInvalidKey = errors.New("key is invalid")
	ErrLimited    = errors.New("rate limit exceeded")
)
