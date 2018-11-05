package model

import (
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Key represents MessageRequest partition key
type Key struct {
	RecipientIdentifier *string `json:"RecipientIdentifier"`
	RequestID           *string `json:"RequestId"`
}

// NewKeyFromString creates new Key from base64 encoded string.
func NewKeyFromString(s string) (*Key, error) {
	if s == "" {
		return nil, nil
	}
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	var key Key
	err = json.Unmarshal(b, &key)
	return &key, err
}

// NewKeyFromMap creates new Key from a map of dynamodb.AttributeValue.
func NewKeyFromMap(m map[string]*dynamodb.AttributeValue) (*Key, error) {
	var key Key
	err := dynamodbattribute.UnmarshalMap(m, &key)
	return &key, err
}

// MarshalMap marshals Key into a value that is suitable for
// Amazon DynamoDB operations.
func (k *Key) MarshalMap() map[string]*dynamodb.AttributeValue {
	v, _ := dynamodbattribute.MarshalMap(k)
	return v
}

// EncodeBase64 marshals and encodes Key to base64 string.
func (k *Key) EncodeBase64() string {
	b, _ := json.Marshal(k)
	return base64.StdEncoding.EncodeToString(b)
}
