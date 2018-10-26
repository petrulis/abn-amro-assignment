package abnerr

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
