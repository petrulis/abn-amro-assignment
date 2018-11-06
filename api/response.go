package api

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/petrulis/abn-amro-assignment/validator"
	"net/http"
)

// NewProxyOkResponse return APIGatewayProxyResponse with status 200 OK and
// body passed as an argument.
func NewProxyOKResponse(body []byte) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: string(body)}
}

// NewProxyErrorResponse builds APIGatewayProxyResponse from Error.
func NewProxyErrorResponse(err Error) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{StatusCode: err.GetStatusCode(), Body: err.GetJSON()}
}

// NewBatchedProxyValidationErrorResponse builds APIGatewayProxyResponse with status 400 from multiple
// validation errors.
func NewBatchedProxyValidationErrorResponse(errs validator.ValidationErrors) events.APIGatewayProxyResponse {
	baseErrs := make([]*baseValidationError, len(errs))
	for i, err := range errs {
		baseErrs[i] = newBaseValidationError(err.Code, err.Message)
	}
	err := newBaseBatchedValidationError("ErrValidation", baseErrs, http.StatusBadRequest)
	return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: err.GetJSON()}
}
