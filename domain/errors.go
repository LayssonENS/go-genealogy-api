package domain

import "errors"

var (
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("your requested Item is not found")
)

type ErrorResponse struct {
	ErrorMessage string `json:"error"`
}
