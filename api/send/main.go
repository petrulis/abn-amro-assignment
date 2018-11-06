package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/uuid"
	"github.com/petrulis/abn-amro-assignment/api"
	"github.com/petrulis/abn-amro-assignment/dynamodbdriver"
	"github.com/petrulis/abn-amro-assignment/model"
	"github.com/petrulis/abn-amro-assignment/validator"
	"os"
)

var (
	dd            *dynamodbdriver.DynamoDbDriver
	defaultRegion string
)

// init initializes long lived resources for Handler.
func init() {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(os.Getenv("REGION"))}))
	defaultRegion = os.Getenv("DEFAULT_REGION")
	dd = dynamodbdriver.New(sess, &dynamodbdriver.DriverConfig{
		MessageRequestTable: aws.String(os.Getenv("DDB_TABLE")),
	})
}

// Handler validates incoming MessageRequest and stores into Amazon DynamoDB database.
// Returned Error Codes:
//   * ErrBadRequest
//   Provided request payload was incorrect and therefore couldn't be processed.
//
//   * ErrInternal
//  The error indicates that something went wrong internally and MessageRequest
//  couldn't be processed and/or stored into database.
//
//   * ErrValidation
//  The error indicates that request payload couldn't pass MessageRequest validation
//  checks.
func Handler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req model.MessageRequest
	err := json.Unmarshal([]byte(event.Body), &req)
	if err != nil {
		return api.NewProxyErrorResponse(api.ErrBadRequest), nil
	}
	v := validator.NewMessageRequestValidator(defaultRegion)
	ok := v.Validate(&req)
	if !ok {
		return api.NewBatchedProxyValidationErrorResponse(v.Errors()), nil
	}
	req.RequestID = uuid.New().String()
	req.DeliveryStatus = model.DeliveryStatusScheduled
	if err = dd.Save(&req); err != nil {
		return api.NewProxyErrorResponse(api.ErrInternal), nil
	}
	return api.NewProxyOKResponse(req.Marshal()), nil
}

func main() {
	lambda.Start(Handler)
}
