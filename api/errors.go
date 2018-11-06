package api

var (
	// ErrNotFound implements Error interface and means that requested resource
	// couldn't be found in Amazon DynamoDB.
	ErrNotFound = New("NotFound", "Requested resource couldn't be found", 404)

	// ErrInternal implements Error interface and means that something bad happend internallly
	// and request couldn't be processed.
	ErrInternal = New("Internal", "Unexpected error occurred. Please try again later", 500)

	// ErrBadRequest implements Error interface and means that request payload was invalid.
	ErrBadRequest = New("BadRequest", "Request payload is not valid", 400)

	// ErrInvalidToken implements Error interface and means provided base64 encoded string coudn't be
	// transformed into model.Key.
	ErrInvalidToken = New("InvalidToken", "Provided token is not valid.", 400)
)

// Error represents api error
type Error interface {
	// Satisfy the generic error interface.
	error

	// Returns the short phrase depicting the classification of the error.
	GetCode() string

	// Returns the error details message.
	GetMessage() string

	// Returns the error http status code
	GetStatusCode() int

	// Return the error json string
	GetJSON() string
}

// New returns an Error object described by the code, message.
//
// If origErr satisfies the Error interface it will not be wrapped within a new
// Error object and will instead be returned.
func New(code, message string, statusCode int) Error {
	return newBaseError(code, message, statusCode)
}

// New returns an Error object described by the code, message.
//
// If origErr satisfies the Error interface it will not be wrapped within a new
// Error object and will instead be returned.
func NewValidationError(code, message string) Error {
	return newBaseValidationError(code, message)
}

// New returns an Error object described by the code, message.
//
// If origErr satisfies the Error interface it will not be wrapped within a new
// Error object and will instead be returned.
func NewBatchedValidationError(code string, errs []*baseValidationError, statusCode int) Error {
	return newBaseBatchedValidationError(code, errs, statusCode)
}
