package dynamodbdriver

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/petrulis/abn-amro-assignment/model"
	"strconv"
	"time"
)

// DynamoDbDriver provides handy operation methods for making requests to
// Amazon DynamoDB.
type DynamoDbDriver struct {
	cfg    *DriverConfig
	client *dynamodb.DynamoDB
}

// New creates a new instance of the DynamoDbDriver with a session and configuration.
func New(sess *session.Session, cfg *DriverConfig) *DynamoDbDriver {
	return &DynamoDbDriver{
		client: dynamodb.New(sess),
		cfg:    cfg,
	}
}

// Update API operation which updates MessageRequest.
// Currently method can only update DeliveryMethod.
func (d *DynamoDbDriver) Update(req *model.MessageRequest) error {
	input := d.newUpdateItemInput(req)
	_, err := d.client.UpdateItem(input)
	return err
}

// Save API opration saves MessageRequest to Amazon DynamoDB
func (d *DynamoDbDriver) Save(req *model.MessageRequest) error {
	item, err := req.MarshalMap()
	if err != nil {
		return nil
	}
	input := &dynamodb.PutItemInput{
		TableName: d.cfg.MessageRequestTable,
		Item:      item,
	}
	_, err = d.client.PutItem(input)
	return err
}

// FindByRecipientIdentifier searches for recipient by RequestIdentifier and exclusiveStartKey.
// The method is handy for paginated responses.
func (d *DynamoDbDriver) FindByRecipientIdentifier(rid *string, exclusiveStartKey *model.Key) (model.MessageRequestList, *model.Key, error) {
	input := d.newFindByRequestIdentifierQueryInput(rid, exclusiveStartKey)
	out, err := d.client.Query(input)
	if err != nil {
		return nil, nil, err
	}
	req, err := d.newMessageRequestList(out.Items)
	if err != nil {
		return nil, nil, err
	}
	key, err := model.NewKeyFromMap(out.LastEvaluatedKey)
	if err != nil {
		return nil, nil, err
	}
	return req, key, err
}

// FindScheduled finds and returns ready to be sent messages.
func (d *DynamoDbDriver) FindScheduled(now time.Time) (model.MessageRequestList, error) {
	input := d.newFindScheduledScanInput(now)
	out, err := d.client.Scan(input)
	if err != nil {
		return nil, err
	}
	return d.newMessageRequestList(out.Items)
}

// newUpdateItemInput creates new update input for MessageRequest
func (d *DynamoDbDriver) newUpdateItemInput(req *model.MessageRequest) *dynamodb.UpdateItemInput {
	input := &dynamodb.UpdateItemInput{
		TableName: d.cfg.MessageRequestTable,
		Key:       req.Key().MarshalMap(),
		ExpressionAttributeNames: map[string]*string{
			"#d": aws.String("DeliveryStatus"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":d": {S: aws.String(model.DeliveryStatusSent)},
		},
		UpdateExpression: aws.String("SET #d = :d"),
	}
	return input
}

// newFindByRequestIdentifierQueryInput creates new QueryInput from MessageRequest.RequestIdentifier
// and exclusiveStartKey.
func (d *DynamoDbDriver) newFindByRequestIdentifierQueryInput(rid *string, exclusiveStartKey *model.Key) *dynamodb.QueryInput {
	input := &dynamodb.QueryInput{
		TableName: d.cfg.MessageRequestTable,
		ExpressionAttributeNames: map[string]*string{
			"#rid": aws.String("RecipientIdentifier"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":rid": {S: rid},
		},
		KeyConditionExpression: aws.String("#rid := rid"),
		ExclusiveStartKey:      exclusiveStartKey.MarshalMap(),
		Limit:                  aws.Int64(10),
	}
	return input
}

// newFindScheduledScanInput creates new ScanInput to search ready to by queued messages.
func (d *DynamoDbDriver) newFindScheduledScanInput(now time.Time) *dynamodb.ScanInput {
	input := &dynamodb.ScanInput{
		TableName: d.cfg.MessageRequestTable,
		ExpressionAttributeNames: map[string]*string{
			"#t": aws.String("SendAt"),
			"#d": aws.String("DeliveryStatus"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {N: aws.String(strconv.FormatInt(now.Unix(), 10))},
			":s": {S: aws.String(model.DeliveryStatusScheduled)},
			":q": {S: aws.String(model.DeliveryStatusQueued)},
		},
		FilterExpression: aws.String("#t <= :t AND #d <> :s AND #d <> :q"),
	}
	return input
}

// newMessageRequestList transforms item map from Amazon DynamoDB to MessageRequestList.
func (d *DynamoDbDriver) newMessageRequestList(items []map[string]*dynamodb.AttributeValue) (model.MessageRequestList, error) {
	var reqs model.MessageRequestList
	err := dynamodbattribute.UnmarshalListOfMaps(items, &reqs)
	return reqs, err
}
