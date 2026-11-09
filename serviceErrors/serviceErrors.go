package serviceErrors

import "errors"

var (
	ErrBadRequest            = errors.New("bad request")
	ErrNotFound              = errors.New("not found")
	ErrInternal              = errors.New("internal error")
	ErrInvalidId             = errors.New("invalid id")
	ErrForbidden             = errors.New("forbidden")
	ErrPublicationNotReached = errors.New("publication not reached")
	ErrNewsValidationFailed  = errors.New("news validation failed")
	ErrUnauthorized          = errors.New("unauthorized")
)
