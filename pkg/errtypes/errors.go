// Package errtypes provides custom error types and error handling utilities.
//
// This package defines application-specific error types that can be used
// throughout the codebase for consistent error handling and reporting.
package errtypes

// CodeError represents an error with an associated error code for categorization.
type CodeError struct {
	Message string
	Code    string
}

// Error implements the error interface.
func (e *CodeError) Error() string {
	return e.Message
}

// New creates a new CodeError with the given message and code.
func New(message, code string) *CodeError {
	return &CodeError{
		Message: message,
		Code:    code,
	}
}
