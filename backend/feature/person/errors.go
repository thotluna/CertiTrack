package person

import (
	"errors"
	"fmt"
)

var (
	ErrPersonNotFound     = errors.New("person not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidStatus      = errors.New("invalid status")
	ErrInvalidPersonData  = errors.New("invalid person data")
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

func IsValidationError(err error) bool {
	_, ok := err.(ValidationError)
	return ok
}

type NotFoundError struct {
	Resource string
	ID       string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s not found", e.Resource, e.ID)
}

func IsNotFoundError(err error) bool {
	_, ok := err.(NotFoundError)
	return ok || errors.Is(err, ErrPersonNotFound)
}
