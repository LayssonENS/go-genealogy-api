package domain

import "errors"

var (
	// ErrRegistrationNotFound will throw if the searched person registration cannot be found
	ErrRegistrationNotFound = errors.New("person registration not found")
	// ErrInvalidSelfRelation will throw if the attempted relationship is with oneself.
	ErrInvalidSelfRelation = errors.New("cannot create a relationship with oneself")
	// ErrDuplicateRelation will throw if the attempted relationship already exists.
	ErrDuplicateRelation = errors.New("relationship already exists")
	// ErrRelationNotFound will throw if the searched relationship cannot be found.
	ErrRelationNotFound = errors.New("relationship not found")
	// ErrIncestuousRelation will throw if the attempted relationship is incestuous.
	ErrIncestuousRelation = errors.New("cannot create incestuous relationship")
)

type ErrorResponse struct {
	ErrorMessage string `json:"error"`
}
