package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/petrulis/abn-amro-assignment/abnerr"
	"github.com/petrulis/abn-amro-assignment/model"
	"net/http"
	"os"
)

var (
	table *string
	ddb   *dynamodb.DynamoDB
)

func init() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION")),
	}))
	ddb = dynamodb.New(sess)
	table = aws.String(os.Getenv("DDB_TABLE"))
}

func newDynamoDbQueryInputFromQueryString(q map[string]string) (*dynamodb.QueryInput, abnerr.Error) {
	r, ok := q["r"]
	if !ok {
		return nil, errBadRequest
	}
	t, ok := q["t"]
	if !ok {
		return nil, errBadRequest
	}
	tok, ok := q["token"]
	if !ok {
		return nil, errBadRequest
	}
	key, err := model.NewKeyFromString(tok)
	if err != nil {
		return nil, errBadRequest
	}
	rid := t + "_" + r
	input := &dynamodb.QueryInput{
		TableName: table,
		ExpressionAttributeNames: map[string]*string{
			"#rid": aws.String("RecipientIdentifier"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":rid": {S: aws.String(rid)},
		},
		KeyConditionExpression: aws.String("#rid := rid"),
		Limit: aws.Int64(10),
	}
	if key != nil {
		input.ExclusiveStartKey = key.MarshalMap()
	}
	return input, nil
}

// Handler ...
func Handler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	input, aerr := newDynamoDbQueryInputFromQueryString(event.QueryStringParameters)

	if aerr != nil {
		return events.APIGatewayProxyResponse{StatusCode: aerr.GetStatusCode(), Body: aerr.GetJSON()}, nil
	}

	out, err := ddb.Query(input)

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: errInternal.GetStatusCode(), Body: errInternal.GetJSON()}, nil
	}

	var items model.MessageRequestList
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &items)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: errInternal.GetStatusCode(), Body: errInternal.GetJSON()}, nil
	}

	response := Output{
		Items: items,
		Count: out.Count,
	}

	if out.LastEvaluatedKey != nil {
		lastKey, err := model.NewKeyFromMap(out.LastEvaluatedKey)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: errInternal.GetStatusCode(), Body: errInternal.GetJSON()}, nil
		}
		nextToken := lastKey.EncodeBase64()
		response.NextToken = aws.String(nextToken)
	}

	return events.APIGatewayProxyResponse{Body: response.String(), StatusCode: http.StatusOK}, nil
}

func main() {
	lambda.Start(Handler)
}
