package apperr

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrBadRequest   = errors.New("bad request")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

type Error struct {
	Kind    error
	Message string
}

func (e *Error) Error() string { return e.Message }

func (e *Error) Unwrap() error { return e.Kind }

func NotFound(msg string) error     { return &Error{Kind: ErrNotFound, Message: msg} }
func Conflict(msg string) error     { return &Error{Kind: ErrConflict, Message: msg} }
func BadRequest(msg string) error   { return &Error{Kind: ErrBadRequest, Message: msg} }
func Unauthorized(msg string) error { return &Error{Kind: ErrUnauthorized, Message: msg} }
func Forbidden(msg string) error    { return &Error{Kind: ErrForbidden, Message: msg} }
