package validator

import (
	"encoding/json"
	"fmt"
	"github.com/petrulis/abn-amro-assignment/model"
	"github.com/ttacon/libphonenumber"
	"regexp"
	"time"
)

const (
	emailRegex = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
)

// ValidationErrors represents multiple ValidationError.
type ValidationErrors []*ValidationError

func (e *ValidationErrors) Marshal() []byte {
	out, _ := json.Marshal(e)
	return out
}

// ValidatorError represents MessageRquest validation error.
type ValidationError struct {
	// Message represents validation message.
	Message string

	// Code represents error code.
	Code    string
}

// MessageRequestValidator represents validator which validates
// MessageRequest.
type MessageRequestValidator struct {
	defaultRegion string
	errors        ValidationErrors
}

// NewMessageRequestValidator creates new MessageRequestValidator.
func NewMessageRequestValidator(defaultRegion string) *MessageRequestValidator {
	return &MessageRequestValidator{defaultRegion: defaultRegion}
}

// Validate validates MessageRequest.
func (v *MessageRequestValidator) Validate(req *model.MessageRequest) bool {
	v.validateIdentifier(req)
	if req.RequestID != "" {
		v.errors = append(v.errors, &ValidationError{
			Message: fmt.Sprint("Trying to update read-only attribute"),
			Code:    "InvalidRequestId",
		})
	}
	now := time.Now().UTC()
	if now.Unix() > req.SendAt {
		v.errors = append(v.errors, &ValidationError{
			Message: fmt.Sprint("Given attribute `SendAt` is not valid timestamp"),
			Code:    "InvalidSendAt",
		})
	}
	if req.DeliveryStatus != "" {
		v.errors = append(v.errors, &ValidationError{
			Message: fmt.Sprint("Trying to update read-only attribute"),
			Code:    "InvalidDeliveryStatus",
		})
	}
	if req.IdentifierType == model.IdentifierTypeEmail && req.Subject == "" {
		v.errors = append(v.errors, &ValidationError{
			Message: fmt.Sprint("Attribute `Subject` is required when `IdentifierType` is `email`"),
			Code:    "InvalidSubject",
		})
	}
	if req.IdentifierType == model.IdentifierTypeEmail && len(req.Subject) < 5 {
		v.errors = append(v.errors, &ValidationError{
			Message: fmt.Sprint("Attribute `Subject` must be at least 5 characters long"),
			Code:    "InvalidSubject",
		})
	}
	return len(v.errors) == 0
}

// validateIdentifier validates message request identifier.
func (v *MessageRequestValidator) validateIdentifier(req *model.MessageRequest) bool {
	if req.IdentifierType == model.IdentifierTypeSMS {
		_, err := libphonenumber.Parse(req.RecipientIdentifier, v.defaultRegion)
		if err != nil {
			v.errors = append(v.errors, &ValidationError{
				Message: fmt.Sprint("Given `RecipientIdentifier` is not valid phone number"),
				Code:    "InvalidPhoneNumber",
			})
		}
		return err == nil
	} else if req.IdentifierType == model.IdentifierTypeEmail {
		rgx := regexp.MustCompile(emailRegex)
		ok := rgx.MatchString(req.RecipientIdentifier)
		if !ok {
			v.errors = append(v.errors, &ValidationError{
				Message: fmt.Sprint("Given `RecipientIdentifier` is not valid email address"),
				Code:    "InvalidEmailAddress",
			})
		}
		return ok
	} else {
		v.errors = append(v.errors, &ValidationError{
			Message: fmt.Sprint("Given `RecipientIdentifier` is not valid identifier. Valid values: sms, email"),
			Code:    "InvalidRequestIdentifier",
		})
		return false
	}
}

// Errors returns ValidationErrors.
func (v *MessageRequestValidator) Errors() ValidationErrors {
	return v.errors
}
