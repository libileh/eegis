package errors

import (
	"database/sql"
	"errors"
	"fmt"
)

type ErrorType string

const (
	NotFound       ErrorType = "NOT_FOUND_ERROR"
	InternalServer ErrorType = "INTERNAL_SERVER_ERROR"
	BadRequest     ErrorType = "BAD_REQUEST_ERROR"
)

// swagger:model CustomError
type CustomError struct {
	// The error message
	Message string `json:"error"`
	// The type of error (not included in JSON response)
	ErrType ErrorType `json:"-"`
}

// Error implements the error interface
func (e *CustomError) Error() string {
	return e.Message
}

// NewCustomError is a constructor function for CustomError
func NewCustomError(errType ErrorType, details string) *CustomError {
	return &CustomError{
		ErrType: errType,
		Message: details,
	}
}

// NewNotFoundError creates a new error for resource not found cases
func NewNotFoundError(details string) *CustomError {
	return NewCustomError(
		NotFound,
		fmt.Sprintf("resource not found: %s", details),
	)
}

// NewInternalServerError creates a new error for internal server error cases
func NewInternalServerError(details string, err ...error) *CustomError {
	var msg string
	if err != nil {
		msg = fmt.Sprintf("internal server error: %s, cause: %v", details, err)
	} else {
		msg = fmt.Sprintf("internal server error: %s", details)
	}
	return NewCustomError(InternalServer, msg)
}

func NewBadRequest(details string, err ...error) *CustomError {
	var msg string
	if err != nil {
		msg = fmt.Sprintf("bad request: %s, cause: %v", details, err)
	} else {
		msg = fmt.Sprintf("bad request: %s", details)
	}
	return NewCustomError(BadRequest, msg)
}

// Handle DB Errors converts database errors to CustomError
func HandleDBError(err error) *CustomError {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return NewNotFoundError(err.Error())
	default:
		return NewInternalServerError(err.Error())
	}
}
