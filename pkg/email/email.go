// Package email implements a client for interacting with AWS SES.
package email

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// Client implements methods for interacting with AWS SES.
type Client struct {
	client      *ses.SES
	senderEmail string
}

// New returns an instance of Client.
func New(iamUserAcessKey, iamUserAccessSecret, region, senderEmail string) *Client {
	cfg := &aws.Config{
		MaxRetries:  aws.Int(3),
		Credentials: credentials.NewStaticCredentials(iamUserAcessKey, iamUserAccessSecret, ""),
		Region:      aws.String(region),
	}
	sess := session.Must(session.NewSession(cfg))

	return &Client{
		client:      ses.New(sess, cfg),
		senderEmail: senderEmail,
	}
}

// Params contains details about the email to be sent to the user.
type Params struct {
	To      string
	Subject string
	Body    string
}

// SendEmail sends an email using the provided params to the user using SES.
func (c *Client) SendEmail(email *Params) error {
	emailInput := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(email.To)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Data: aws.String(email.Body),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(email.Subject),
			},
		},
		Source: aws.String(c.senderEmail),
	}
	_, err := c.client.SendEmail(emailInput)
	return err
}
