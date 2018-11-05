package model

import (
	"encoding/json"
	"github.com/ttacon/libphonenumber"
	"time"
)

const (
	IdentifierTypeSMS = "sms"
	IdentifierTypeEmail = "email"
	DeliveryStatusScheduled = "Scheduled"
)

type ValidationErrors []*ValidationError

func (e *ValidationErrors) Marshal() []byte {
	out, _ := json.Marshal(e)
	return out
}

type ValidationError struct {
}

type MessageRequestValidator struct {
	defaultRegion string
	errors        ValidationErrors
}

func NewMessageRequestValidator(defaultRegion string) *MessageRequestValidator {
	return &MessageRequestValidator{defaultRegion: defaultRegion}
}

func (v *MessageRequestValidator) Validate(request *MessageRequest) bool {
	if v.validateIdentifier(request) == false {
		return false
	}
	now := time.Now().UTC()
	if now.Unix() > request.SendAt {
		return false
	}
	if request.DeliveryStatus != "" {
		return false
	}
	return true
}

func (v *MessageRequestValidator) validateIdentifier(request *MessageRequest) bool {
	if request.IdentifierType == IdentifierTypeSMS {
		_, err := libphonenumber.Parse(request.RecipientIdentifier, v.defaultRegion)
		return err == nil
	} else if request.IdentifierType == IdentifierTypeEmail {
		return true
	} else {
		return false
	}
}

func (v *MessageRequestValidator) Errors() ValidationErrors {
	return v.errors
}

// MessageRequest represents MessageRequest table
type MessageRequest struct {
	IdentifierType      string `json:"IdentifierType"`
	RecipientIdentifier string `json:"RecipientIdentifier"`
	RequestID           string `json:"RequestId"`
	Body                string `json:"Body"`
	SendAt              int64  `json:"SendAt"`
	DeliveryStatus      string `json:"DeliveryStatus"`
}

// Marshal serializes the MessageRequest to json byte array
func (r *MessageRequest) Marshal() []byte {
	b, _ := json.Marshal(r)
	return b
}

// MessageRequestList represents a list of MessageRequest items
type MessageRequestList []*MessageRequest
