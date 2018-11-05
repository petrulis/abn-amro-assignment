package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/petrulis/abn-amro-assignment/model"
	"encoding/json"
	"github.com/go-errors/errors"
	"github.com/aws/aws-sdk-go/service/sns"
)

var ddb *dynamodb.DynamoDB
var tbl *string
var snc *sns.SNS

func init() {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(os.Getenv("REGION"))}))
	ddb = dynamodb.New(sess)
	tbl = aws.String(os.Getenv("DDB_TABLE"))
	snc = sns.New(sess)
}

func sendSMS(request *model.MessageRequest) error {
	input := &sns.PublishInput{
		PhoneNumber: aws.String(request.RecipientIdentifier),
		Message: aws.String(request.Body),
	}
	_, err := snc.Publish(input)
	return err
}

func sendEmail(request *model.MessageRequest) error {
	return nil
}

func send(request *model.MessageRequest) error {
	switch request.IdentifierType {
	case model.IdentifierTypeSMS:
		return sendSMS(request)
	case model.IdentifierTypeEmail:
		return sendEmail(request)
	default:
		return errors.New("unsupported message type")
	}
}

func Handler(ctx context.Context, event events.SQSEvent) error {
	for _, record := range event.Records {
		body := record.Body
		var request model.MessageRequest
		err := json.Unmarshal([]byte(body), &request)
		if err != nil {
			return err
		}
		err = send(&request)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}