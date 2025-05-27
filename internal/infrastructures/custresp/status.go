package custresp

import (
	"errors"
)

var (
	ErrTooManyRequest  = errors.New("too many request")
	ErrRequestTooEarly = errors.New("request too early")
	ErrInvalidRequest  = errors.New("invalid request")
)
