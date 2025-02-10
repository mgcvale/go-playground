package util

import (
	"errors"
	"fmt"
)

var (
	BadRequestError   = errors.New("Bad Request")
	InternalError     = errors.New("Internal Server Error")
	ConflictError     = errors.New("Conflict")
	UnauthorizedError = errors.New("Unauthorized")
	ForbiddenError    = errors.New("Forbidden")
	NotFoundError     = errors.New("Not Found")
)

type AppError struct {
	Err         error
	description string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("AppError: %s due to %s", e.Err.Error(), e.description)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(err error, description string) *AppError {
	return &AppError{Err: err, description: description}
}
