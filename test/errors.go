package test

import (
	"errors"
)

// ErrUnknown is used when a test could not be found.
var ErrUnknown = errors.New("unknown test")

// ErrInvalidArgument is returned when one or more arguments are invalid.
var ErrInvalidArgument = errors.New("invalid argument")
