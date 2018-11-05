package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/petrulis/abn-amro-assignment/dynamodbdriver"
	"github.com/petrulis/abn-amro-assignment/model"
	"net/http"
	"os"
)

var (
	dd *dynamodbdriver.DynamoDbDriver
)

func init() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION")),
	}))
	dd = dynamodbdriver.New(sess, &dynamodbdriver.DriverConfig{
		MessageRequestTable: aws.String(os.Getenv("DDB_TABLE")),
	})
}

// Handler ...
func Handler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	rid, ok := event.QueryStringParameters["requestIdentifier"]
	if !ok {
		return events.APIGatewayProxyResponse{StatusCode: errBadRequest.GetStatusCode(), Body: errBadRequest.GetJSON()}, nil
	}
	token, ok := event.QueryStringParameters["token"]
	var exclusiveStartKey *model.Key
	if ok {
		key, err := model.NewKeyFromString(token)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: errBadRequest.GetStatusCode(), Body: errBadRequest.GetJSON()}, nil
		}
		exclusiveStartKey = key
	}
	requests, lastEvaluatedKey, err := dd.FindByRecipientIdentifier(aws.String(rid), exclusiveStartKey)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: errInternal.GetStatusCode(), Body: errInternal.GetJSON()}, nil
	}
	response := newResponse(requests, lastEvaluatedKey)
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: response.String()}, nil
}

func main() {
	lambda.Start(Handler)
}
