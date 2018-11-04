package model

import (
	"github.com/ttacon/libphonenumber"
	"encoding/json"
)

const IdentifierTypeSMS = "sms"
const IdentifierTypeEmail = "email"

type ValidationErrors []*ValidationError

func (e *ValidationErrors) Marshal() []byte {
	out, _ := json.Marshal(e)
	return out
}

type ValidationError struct {

}

type MessageRequestValidator struct {
	defaultRegion string
	errors ValidationErrors
}

func NewMessageRequestValidator(defaultRegion string) *MessageRequestValidator {
	return &MessageRequestValidator{}
}

func (v *MessageRequestValidator) Validate(request *MessageRequest) bool {
	if request.IdentifierType == IdentifierTypeSMS {
		_, err := libphonenumber.Parse(request.RecipientIdentifier, v.defaultRegion)
		return err == nil
	} else if request.IdentifierType == IdentifierTypeEmail {
		return true
	}
	return false
}

func (v *MessageRequestValidator) Errors() ValidationErrors {
	return v.errors
}

// MessageRequest represents MessageRequest table
type MessageRequest struct {
	IdentifierType      string `json:"identifierType"`
	RecipientIdentifier string `json:"recipientIdentifier"`
	RequestID           string `json:"requestId"`
	Body                string `json:"body"`
}

// MessageRequestList represents a list of MessageRequest items
type MessageRequestList []*MessageRequest