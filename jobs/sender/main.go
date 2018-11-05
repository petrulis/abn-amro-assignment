package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/go-errors/errors"
	"github.com/petrulis/abn-amro-assignment/model"
	"os"
)

const sender = "no-reply@segmnt.net"
const charSet = "UTF-8"

var ddb *dynamodb.DynamoDB
var tbl *string
var snc *sns.SNS
var sec *ses.SES

func init() {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(os.Getenv("REGION"))}))
	ddb = dynamodb.New(sess)
	tbl = aws.String(os.Getenv("DDB_TABLE"))
	snc = sns.New(sess)
	sec = ses.New(sess)
}

func sendSMS(request *model.MessageRequest) error {
	input := &sns.PublishInput{
		PhoneNumber: aws.String(request.RecipientIdentifier),
		Message:     aws.String(request.Body),
	}
	_, err := snc.Publish(input)
	return err
}

func sendEmail(request *model.MessageRequest) error {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(request.RecipientIdentifier)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(request.Body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charSet),
				Data:    aws.String("Email from ABN AMRO Assignment"),
			},
		},
		Source: aws.String(sender),
	}
	_, err := sec.SendEmail(input)
	return err
}

func updateStatus(request *model.MessageRequest) error {
	input := &dynamodb.UpdateItemInput{
		TableName: tbl,
		Key: map[string]*dynamodb.AttributeValue{
			"RecipientIdentifier": {S: aws.String(request.RecipientIdentifier)},
			"RequestId":           {S: aws.String(request.RequestID)},
		},
		ExpressionAttributeNames: map[string]*string{
			"#d": aws.String("DeliveryStatus"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":d": {S: aws.String(model.DeliveryStatusSent)},
		},
		UpdateExpression: aws.String("#d = :d"),
	}
	_, err := ddb.UpdateItem(input)
	return err
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
		err = updateStatus(&request)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}
