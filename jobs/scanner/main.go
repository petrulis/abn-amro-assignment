package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"time"
	"strconv"
	"github.com/petrulis/abn-amro-assignment/model"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sqs"
	"encoding/json"
)

var ddb *dynamodb.DynamoDB
var sqc *sqs.SQS
var tbl *string
var defaultRegion string
var queueUrl *string

func init() {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(os.Getenv("REGION"))}))
	ddb = dynamodb.New(sess)
	tbl = aws.String(os.Getenv("DDB_TABLE"))
	sqc = sqs.New(sess)
	queueUrl = aws.String(os.Getenv("QUEUE_URL"))
}

func Handler() error {
	now := time.Now().Unix()
	input := &dynamodb.ScanInput{
		TableName: tbl,
		ExpressionAttributeNames: map[string]*string{
			"#t": aws.String("SendAt"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {N: aws.String(strconv.FormatInt(now, 10))},
		},
		FilterExpression: aws.String("#t <= :t"),
	}
	out, err := ddb.Scan(input)
	if err != nil {
		return err
	}
	var requests model.MessageRequestList
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &requests)
	if err != nil {
		return err
	}
	for _, request := range requests {
		b, _ := json.Marshal(&request)
		input := &sqs.SendMessageInput{
			QueueUrl: queueUrl,
			MessageBody: aws.String(string(b)),
		}
		sqc.SendMessage(input)
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}