// Package apperr is the error type handlers and logic return; middleware.Response turns it into
// { "error": { "code", "message", "fields" } } with the matching HTTP status.
package apperr

import "net/http"

type Error struct {
	Status  int               `json:"-"`
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func (e *Error) Error() string { return e.Code + ": " + e.Message }

// Validation reports bad input; fields maps the JSON field name to its message.
func Validation(fields map[string]string) *Error {
	return &Error{http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Please fix the highlighted fields.", fields}
}

func Unauthenticated() *Error {
	return &Error{http.StatusUnauthorized, "UNAUTHENTICATED", "Please log in.", nil}
}

func Forbidden() *Error {
	return &Error{http.StatusForbidden, "FORBIDDEN", "You do not have access to this.", nil}
}

// NotFound is also the answer for another customer's record, so its existence is never revealed.
func NotFound() *Error {
	return &Error{http.StatusNotFound, "NOT_FOUND", "Not found.", nil}
}

func Conflict(message string) *Error {
	return &Error{http.StatusConflict, "CONFLICT", message, nil}
}

func Internal() *Error {
	return &Error{http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong. Please try again.", nil}
}
