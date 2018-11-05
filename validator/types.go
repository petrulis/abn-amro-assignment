package validator

import (
	"encoding/json"
	"github.com/petrulis/abn-amro-assignment/model"
	"github.com/ttacon/libphonenumber"
	"time"
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

func (v *MessageRequestValidator) Validate(request *model.MessageRequest) bool {
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

func (v *MessageRequestValidator) validateIdentifier(request *model.MessageRequest) bool {
	if request.IdentifierType == model.IdentifierTypeSMS {
		_, err := libphonenumber.Parse(request.RecipientIdentifier, v.defaultRegion)
		return err == nil
	} else if request.IdentifierType == model.IdentifierTypeEmail {
		return true
	} else {
		return false
	}
}

func (v *MessageRequestValidator) Errors() ValidationErrors {
	return v.errors
}
