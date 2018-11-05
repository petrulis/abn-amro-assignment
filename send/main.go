package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/petrulis/abn-amro-assignment/model"
	"encoding/json"
	"net/http"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"os"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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
	item, err := dynamodbattribute.MarshalMap(&request)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: err.Error()}, nil
	}
	input := &dynamodb.PutItemInput{
		TableName: tbl,
		Item: item,
	}
	_, err = ddb.PutItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: err.Error()}, nil
	}
	return events.APIGatewayProxyResponse{}, nil
}

func main() {
	lambda.Start(Handler)
}
