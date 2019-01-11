package errors

import "errors"

var (
	ErrPropNotFound = errors.New("property not found")
	ErrNotValidULID = errors.New("not valid ULID")
)
