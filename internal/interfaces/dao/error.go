package dao

import "errors"

var (
	ErrNoUpdateHappened = errors.New("no affected row during update")
	ErrNilParam         = errors.New("parameter is nil")
	ErrDuplicate        = errors.New("duplicate")
	ErrNoResult         = errors.New("no result")
)
