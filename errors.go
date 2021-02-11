package msgfmt

import (
	"errors"
)

var (
	ErrNotImplemented  = errors.New("not implemented")
	ErrEmptyKey        = errors.New("empty key")
	ErrInvalidTemplate = errors.New("invalid template")
)
