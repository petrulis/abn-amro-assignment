package model

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	IdentifierTypeSMS       = "sms"
	IdentifierTypeEmail     = "email"
	DeliveryStatusScheduled = "Scheduled"
	DeliveryStatusSent      = "Sent"
	DeliveryStatusQueued = "Queued"
)

// MessageRequest represents MessageRequest table
type MessageRequest struct {
	IdentifierType      string `json:"IdentifierType"`
	RecipientIdentifier string `json:"RecipientIdentifier"`
	RequestID           string `json:"RequestId"`
	Body                string `json:"Body"`
	SendAt              int64  `json:"SendAt"`
	DeliveryStatus      string `json:"DeliveryStatus"`
}

// NewMessageRequestFromString ...
func NewMessageRequestFromString(body string) (*MessageRequest, error) {
	var request MessageRequest
	err := json.Unmarshal([]byte(body), &request)
	return &request, err
}

// Marshal serializes the MessageRequest to json byte array
func (r *MessageRequest) Marshal() []byte {
	b, _ := json.Marshal(r)
	return b
}

// MarshalMp serializes the MessageRequest to map[string]*dynamodb.AttributeValue
func (r *MessageRequest) MarshalMap() (map[string]*dynamodb.AttributeValue, error) {
	out, err := dynamodbattribute.MarshalMap(r)
	return out, err
}

func (r *MessageRequest) Key() *Key {
	key := &Key{
		RecipientIdentifier: &r.RecipientIdentifier,
		RequestID: &r.RequestID,
	}
	return key
}

// MessageRequestList represents a list of MessageRequest items
type MessageRequestList []*MessageRequest