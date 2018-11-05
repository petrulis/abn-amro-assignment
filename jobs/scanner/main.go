package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/petrulis/abn-amro-assignment/dynamodbdriver"
	"os"
	"time"
)

var (
	sqc      *sqs.SQS
	queueUrl *string
	dd       *dynamodbdriver.DynamoDbDriver
)

func init() {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(os.Getenv("REGION"))}))
	sqc = sqs.New(sess)
	queueUrl = aws.String(os.Getenv("QUEUE_URL"))
	dd = dynamodbdriver.New(sess, &dynamodbdriver.DriverConfig{
		MessageRequestTable: aws.String(os.Getenv("DDB_TABLE")),
	})
}

func Handler() error {
	reqs, err := dd.FindScheduled(time.Now())
	if err != nil {
		return err
	}
	for _, req := range reqs {
		body := req.Marshal()
		input := &sqs.SendMessageInput{
			QueueUrl:    queueUrl,
			MessageBody: aws.String(string(body)),
		}
		sqc.SendMessage(input)
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}
