package callback

import "errors"

var (
	ErrInvalidCallbackType = errors.New("invalid callback type")
	ErrInvalidData         = errors.New("invalid data")
)
