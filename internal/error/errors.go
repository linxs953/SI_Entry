package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

// Error represents a business error with code, message and stack trace
type Error struct {
	Code    ErrorCode
	Message string
	Err     error // Original error if any
	stack   error // Stack trace
}

// NewError creates a new error with error code and message
func NewError(code ErrorCode, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return &Error{
		Code:    code,
		Message: msg,
		stack:   errors.New(msg),
	}
}

// WrapError wraps an existing error with error code and message
func WrapError(code ErrorCode, err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(format, args...)
	return &Error{
		Code:    code,
		Message: msg,
		Err:     err,
		stack:   errors.Wrap(err, msg),
	}
}

// Wrap wraps an error with additional message
func Wrap(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	// If it's already our error type, just wrap the internal error
	if e, ok := err.(*Error); ok {
		msg := fmt.Sprintf(format, args...)
		e.stack = errors.Wrap(e.stack, msg)
		return e
	}

	// Otherwise create a new error with InternalError code
	return WrapError(InternalError, err, format, args...)
}

// Cause returns the underlying cause of the error
func Cause(err error) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*Error); ok {
		if e.Err != nil {
			return Cause(e.Err)
		}
		return e
	}

	return errors.Cause(err)
}

// GetErrorCode extracts the error code from an error
func GetErrorCode(err error) ErrorCode {
	if err == nil {
		return Success
	}

	if e, ok := err.(*Error); ok {
		return e.Code
	}

	return InternalError
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("error code: %d, message: %s, cause: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("error code: %d, message: %s", e.Code, e.Message)
}

// Stack returns the full stack trace of the error
func (e *Error) Stack() string {
	return fmt.Sprintf("%+v", e.stack)
}

// Is implements errors.Is interface for error comparison
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// Unwrap implements errors.Unwrap interface
func (e *Error) Unwrap() error {
	return e.Err
}
