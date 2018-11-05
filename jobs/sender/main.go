package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-errors/errors"
	"github.com/petrulis/abn-amro-assignment/dynamodbdriver"
	"github.com/petrulis/abn-amro-assignment/model"
	"github.com/petrulis/abn-amro-assignment/sender"
	"os"
)

const charSet = "UTF-8"

var (
	dd  *dynamodbdriver.DynamoDbDriver
	mta *sender.EmailSender
	sta *sender.SMSSender
)

func init() {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(os.Getenv("REGION"))}))
	dd = dynamodbdriver.New(sess, &dynamodbdriver.DriverConfig{
		MessageRequestTable: aws.String(os.Getenv("DDB_TABLE")),
	})
	sta = sender.NewSMSSender(sess)
	mta = sender.NewEmailSender(sess, &sender.EmailSenderConfig{
		CharSet:    aws.String(charSet),
		SenderName: aws.String("info@segmnt.net"),
	})
}

func send(req *model.MessageRequest) error {
	switch req.IdentifierType {
	case model.IdentifierTypeSMS:
		return sta.Send(req)
	case model.IdentifierTypeEmail:
		return mta.Send(req)
	default:
		return errors.New("unsupported identifier type")
	}
}

func Handler(ctx context.Context, event events.SQSEvent) error {
	for _, rec := range event.Records {
		req, err := model.NewMessageRequestFromString(rec.Body)
		if err != nil {
			return err
		}
		req.DeliveryStatus = model.DeliveryStatusSent
		if err = dd.Update(req); err != nil {
			return err
		}
		if err = send(req); err != nil {
			req.DeliveryStatus = model.DeliveryStatusQueued
			if err = dd.Update(req); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}
