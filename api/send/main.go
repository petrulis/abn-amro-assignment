package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/uuid"
	"github.com/petrulis/abn-amro-assignment/dynamodbdriver"
	"github.com/petrulis/abn-amro-assignment/model"
	"github.com/petrulis/abn-amro-assignment/validator"
	"net/http"
	"os"
)

var (
	dd            *dynamodbdriver.DynamoDbDriver
	defaultRegion string
)

func init() {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(os.Getenv("REGION"))}))
	defaultRegion = os.Getenv("DEFAULT_REGION")
	dd = dynamodbdriver.New(sess, &dynamodbdriver.DriverConfig{
		MessageRequestTable: aws.String(os.Getenv("DDB_TABLE")),
	})
}

func Handler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req model.MessageRequest
	err := json.Unmarshal([]byte(event.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: err.Error()}, nil
	}
	v := validator.NewMessageRequestValidator(defaultRegion)
	ok := v.Validate(&req)
	if !ok {
		errors := v.Errors()
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: string(errors.Marshal())}, nil
	}
	req.RequestID = uuid.New().String()
	req.DeliveryStatus = model.DeliveryStatusScheduled
	dd.Save(&req)
	if err = dd.Save(&req); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: err.Error()}, nil
	}
	body := string(req.Marshal())
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: body}, nil
}

func main() {
	lambda.Start(Handler)
}
