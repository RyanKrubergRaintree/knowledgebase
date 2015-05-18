package kb

import "errors"

var (
	ErrInvalid       = errors.New("invalid argument")
	ErrUnknownAction = errors.New("unknown action")
)
