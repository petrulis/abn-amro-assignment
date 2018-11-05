package sender

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/petrulis/abn-amro-assignment/model"
)

type SMSSender struct {
	client *sns.SNS
}

// NewSMSSender ...
func NewSMSSender(sess *session.Session) *SMSSender {
	return &SMSSender{
		client: sns.New(sess),
	}
}

func (s *SMSSender) Send(request *model.MessageRequest) error {
	input := &sns.PublishInput{
		PhoneNumber: aws.String(request.RecipientIdentifier),
		Message:     aws.String(request.Body),
	}
	_, err := s.client.Publish(input)
	return err
}

type EmailSenderConfig struct {
	CharSet    *string
	SenderName *string
}

type EmailSender struct {
	client *ses.SES
	cfg    *EmailSenderConfig
}

// NewEmailSender ...
func NewEmailSender(sess *session.Session, cfg *EmailSenderConfig) *EmailSender {
	return &EmailSender{
		client: ses.New(sess),
		cfg:    cfg,
	}
}

// Send ...
func (s *EmailSender) Send(request *model.MessageRequest) error {
	input := s.newSendEmailInput(request)
	_, err := s.client.SendEmail(input)
	return err
}

func (s *EmailSender) newSendEmailInput(request *model.MessageRequest) *ses.SendEmailInput {
	html := &ses.Content{
		Charset: s.cfg.CharSet,
		Data:    aws.String(request.Body),
	}
	body := &ses.Body{Html: html}
	subject := &ses.Content{
		Charset: s.cfg.CharSet,
		Data:    aws.String(""),
	}
	message := &ses.Message{
		Body: body, Subject: subject,
	}
	destination := &ses.Destination{
		ToAddresses: []*string{aws.String(request.RecipientIdentifier)},
	}
	input := &ses.SendEmailInput{
		Destination: destination,
		Message:     message,
		Source:      s.cfg.SenderName,
	}
	return input
}
