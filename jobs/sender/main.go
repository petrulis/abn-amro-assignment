package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
)

var ddb *dynamodb.DynamoDB
var tbl *string

func init() {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(os.Getenv("REGION"))}))
	ddb = dynamodb.New(sess)
	tbl = aws.String(os.Getenv("DDB_TABLE"))
}

func Handler(ctx context.Context, event events.SQSEvent) error {
	fmt.Println(ctx)
	fmt.Println(event)
	return nil
}

func main() {
	lambda.Start(Handler)
}