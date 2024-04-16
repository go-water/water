package consterr

import (
	"errors"
)

var (
	ErrLimited = errors.New("rate limit exceeded")
)
