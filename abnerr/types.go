package abnerr

import (
	"encoding/json"
	"fmt"
)

type baseError struct {
	// Classification of error
	Code string `json:"code"`

	// Detailed information about error
	Message string `json:"message"`

	// Http status code
	StatusCode int `json:"-"`
}

// newBaseError returns an error object for the code, message, and errors.
//
// code is a short no whitespace phrase depicting the classification of
// the error that is being created.
//
// message is the free flow string containing detailed information about the
// error.
func newBaseError(code, message string, statusCode int) *baseError {
	b := &baseError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
	return b
}

// Error returns the string representation of the error.
// Satisfies the error interface.
func (b baseError) Error() string {
	return fmt.Sprintf("%s: %s", b.Code, b.Message)
}

// GetCode returns the short phrase depicting the classification of the error.
func (b baseError) GetCode() string {
	return b.Code
}

// GetMessage returns the error details message.
func (b baseError) GetMessage() string {
	return b.Message
}

// GetStatusCode returns http status code
func (b baseError) GetStatusCode() int {
	return b.StatusCode
}

// GetJSON marshals the error and return sting representation of if.
func (b baseError) GetJSON() string {
	bytes, _ := json.Marshal(b)
	return string(bytes)
}
