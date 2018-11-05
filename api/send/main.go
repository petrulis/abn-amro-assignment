package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"github.com/petrulis/abn-amro-assignment/model"
	"net/http"
	"os"
)

var ddb *dynamodb.DynamoDB
var tbl *string
var defaultRegion string

func init() {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(os.Getenv("REGION"))}))
	ddb = dynamodb.New(sess)
	tbl = aws.String(os.Getenv("DDB_TABLE"))
	defaultRegion = os.Getenv("DEFAULT_REGION")
}

func Handler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var request model.MessageRequest
	err := json.Unmarshal([]byte(event.Body), &request)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: err.Error()}, nil
	}
	validator := model.NewMessageRequestValidator(defaultRegion)
	ok := validator.Validate(&request)
	if !ok {
		errors := validator.Errors()
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: string(errors.Marshal())}, nil
	}
	request.RequestID = uuid.New().String()
	request.DeliveryStatus = model.DeliveryStatusScheduled
	item, err := dynamodbattribute.MarshalMap(&request)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: err.Error()}, nil
	}
	input := &dynamodb.PutItemInput{
		TableName: tbl,
		Item:      item,
	}
	_, err = ddb.PutItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: err.Error()}, nil
	}
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: string(request.Marshal())}, nil
}

func main() {
	lambda.Start(Handler)
}
