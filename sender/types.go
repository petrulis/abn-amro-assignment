package sender

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/petrulis/abn-amro-assignment/model"
)

// SMSSender represents SMS delivery method.
type SMSSender struct {
	client *sns.SNS
}

// NewSMSSender creates new SMSSender form session.
func NewSMSSender(sess *session.Session) *SMSSender {
	return &SMSSender{
		client: sns.New(sess),
	}
}

// Send publishes new message to Amazon SNS.
func (s *SMSSender) Send(req *model.MessageRequest) error {
	input := &sns.PublishInput{
		PhoneNumber: aws.String(req.RecipientIdentifier),
		Message:     aws.String(req.Body),
	}
	_, err := s.client.Publish(input)
	fmt.Println(err)
	return err
}

// EmailSenderConfig represents EmailSender condfiguration.
type EmailSenderConfig struct {
	// CharSet represents message and subject encoding.
	CharSet *string

	// SenderName represents email from which new message can be sent.
	SenderName *string
}

// EmailSender represents Email delivery method.
type EmailSender struct {
	client *ses.SES
	cfg    *EmailSenderConfig
}

// NewEmailSender creates new EmailSender from session and configuration.
func NewEmailSender(sess *session.Session, cfg *EmailSenderConfig) *EmailSender {
	return &EmailSender{
		client: ses.New(sess),
		cfg:    cfg,
	}
}

// Send executes Amazon SES SenEmail operation.
func (s *EmailSender) Send(request *model.MessageRequest) error {
	input := s.newSendEmailInput(request)
	_, err := s.client.SendEmail(input)
	fmt.Println(err)
	return err
}

// newSendEmailInput creates new Amazon SES SendEmailInput from MessageRequest.
func (s *EmailSender) newSendEmailInput(req *model.MessageRequest) *ses.SendEmailInput {
	html := &ses.Content{
		Charset: s.cfg.CharSet,
		Data:    aws.String(req.Body),
	}
	body := &ses.Body{Html: html}
	subject := &ses.Content{
		Charset: s.cfg.CharSet,
		Data:    aws.String(req.Subject),
	}
	message := &ses.Message{
		Body: body, Subject: subject,
	}
	destination := &ses.Destination{
		ToAddresses: []*string{aws.String(req.RecipientIdentifier)},
	}
	input := &ses.SendEmailInput{
		Destination: destination,
		Message:     message,
		Source:      s.cfg.SenderName,
	}
	return input
}
