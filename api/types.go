package api

import (
	"encoding/json"
	"fmt"
	"net/http"
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

type baseValidationError struct {

	// Classification of error
	Code string `json:"code"`

	// Detailed information about error
	Message string `json:"message"`

	// Http status code
	StatusCode int `json:"-"`
}

// newBaseValidationError returns an error object for the code, message, and errors.
func newBaseValidationError(code, message string) *baseValidationError {
	b := &baseValidationError{
		Code:       code,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
	return b
}

// Error returns the string representation of the error.
// Satisfies the error interface.
func (b baseValidationError) Error() string {
	return fmt.Sprintf("%s: %s", b.Code, b.Message)
}

// GetCode returns the short phrase depicting the classification of the error.
func (b baseValidationError) GetCode() string {
	return b.Code
}

// GetMessage returns the error details message.
func (b baseValidationError) GetMessage() string {
	return b.Message
}

// GetStatusCode returns http status code
func (b baseValidationError) GetStatusCode() int {
	return b.StatusCode
}

// GetJSON marshals the error and return sting representation of if.
func (b baseValidationError) GetJSON() string {
	bytes, _ := json.Marshal(b)
	return string(bytes)
}

// baseBatchedValidationError returns an error object for the code, message, and errors.
type baseBatchedValidationError struct {

	// Error list
	Errors []*baseValidationError `json:"errors"`

	// Classification of error
	Code string `json:"code"`

	// Http status code
	StatusCode int `json:"-"`
}

// newBaseBatchedValidationError returns an error object for the code, message, and errors.
func newBaseBatchedValidationError(code string, errs []*baseValidationError, statusCode int) *baseBatchedValidationError {
	b := &baseBatchedValidationError{
		Code:       code,
		Errors:     errs,
		StatusCode: statusCode,
	}
	return b
}

// Error returns the string representation of the error.
// Satisfies the error interface.
func (b baseBatchedValidationError) Error() string {
	return fmt.Sprintf("%d: %s", b.StatusCode, b.Errors)
}

// GetCode returns the short phrase depicting the classification of the error.
func (b baseBatchedValidationError) GetErrors() []*baseValidationError {
	return b.Errors
}

// GetCode returns the short phrase depicting the classification of the error.
func (b baseBatchedValidationError) GetCode() string {
	return b.Code
}

// GetMessage returns the error details message.
func (b baseBatchedValidationError) GetMessage() string {
	return ""
}

// GetStatusCode returns http status code
func (b baseBatchedValidationError) GetStatusCode() int {
	return b.StatusCode
}

// GetJSON marshals the error and return sting representation of if.
func (b baseBatchedValidationError) GetJSON() string {
	bytes, _ := json.Marshal(b)
	return string(bytes)
}
