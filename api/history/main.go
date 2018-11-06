package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/petrulis/abn-amro-assignment/api"
	"github.com/petrulis/abn-amro-assignment/dynamodbdriver"
	"github.com/petrulis/abn-amro-assignment/model"
	"os"
)

var (
	dd *dynamodbdriver.DynamoDbDriver
)

// init initializes long lived resources for Handler.
func init() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION")),
	}))
	dd = dynamodbdriver.New(sess, &dynamodbdriver.DriverConfig{
		MessageRequestTable: aws.String(os.Getenv("DDB_TABLE")),
	})
}

// Handler validates incoming request and searches for MessageRequests by RecipientIdentifier
// and base64 encoded token.
// Returned Error Codes:
//   * ErrBadRequest
//   Provided request payload was incorrect and therefore couldn't be processed.
func Handler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	rid, ok := event.QueryStringParameters["recipientIdentifier"]
	if !ok {
		return api.NewProxyErrorResponse(api.ErrBadRequest), nil
	}
	token, ok := event.QueryStringParameters["token"]
	var exclusiveStartKey *model.Key
	if ok {
		key, err := model.NewKeyFromString(token)
		if err != nil {
			return api.NewProxyErrorResponse(api.ErrInvalidToken), nil
		}
		exclusiveStartKey = key
	}
	requests, lastEvaluatedKey, err := dd.FindByRecipientIdentifier(aws.String(rid), exclusiveStartKey)
	if err != nil {
		return api.NewProxyErrorResponse(api.ErrInternal), nil
	}
	response := newResponse(requests, lastEvaluatedKey)
	return api.NewProxyOKResponse(response.Marshal()), nil
}

func main() {
	lambda.Start(Handler)
}
