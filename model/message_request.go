package model

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	// IdentifierTypeSMS represents sms identifier
	IdentifierTypeSMS = "sms"

	// IdentifierTypeEmail represents email identifier
	IdentifierTypeEmail = "email"

	// DeliveryStatusScheduled represents message delivery state
	// when MessageRequest is already persisted but not yet queued.
	DeliveryStatusScheduled = "Scheduled"

	// DeliveryStatusSent represents message delivery state
	// when MessageRequest is already persisted, queued and sent
	// via one of supported delivery channels.
	DeliveryStatusSent = "Sent"

	// DeliveryStatusQueued represents message delivery state
	// when MessageRequest is already persisted and queued but
	// not sent yet.
	DeliveryStatusQueued = "Queued"
)

// MessageRequest represents MessageRequest table.
type MessageRequest struct {
	IdentifierType      string `json:"IdentifierType"`
	RecipientIdentifier string `json:"RecipientIdentifier"`
	RequestID           string `json:"RequestId" valid:"readonly"`
	Body                string `json:"Body"`
	SendAt              int64  `json:"SendAt"`
	Subject             string `json:"Subject"`
	DeliveryStatus      string `json:"DeliveryStatus" valid:"readonly"`
}

// NewMessageRequestFromString creates new MessageRequest from string.
func NewMessageRequestFromString(body string) (*MessageRequest, error) {
	var request MessageRequest
	err := json.Unmarshal([]byte(body), &request)
	return &request, err
}

// Marshal serializes the MessageRequest to json byte array.
func (r *MessageRequest) Marshal() []byte {
	b, _ := json.Marshal(r)
	return b
}

// MarshalMp serializes the MessageRequest to map[string]*dynamodb.AttributeValue.
func (r *MessageRequest) MarshalMap() (map[string]*dynamodb.AttributeValue, error) {
	out, err := dynamodbattribute.MarshalMap(r)
	return out, err
}

// Key creates new model.Key from MessageRequest.
func (r *MessageRequest) Key() *Key {
	key := &Key{
		RecipientIdentifier: &r.RecipientIdentifier,
		RequestID:           &r.RequestID,
	}
	return key
}

// MessageRequestList represents a list of MessageRequest items.
type MessageRequestList []*MessageRequest
