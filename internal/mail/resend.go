package mail

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

type ResendService struct {
	client *resend.Client
}

func NewResendService(apiKey string) (Service, error) {
	client := resend.NewClient(apiKey)
	return &ResendService{client: client}, nil
}

func (r *ResendService) Send(to, subject, body string) error {

	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", os.Getenv("APP_NAME"), os.Getenv("RESEND_FROM_EMAIL")),
		To:      []string{to},
		Html:    body,
		Subject: subject,
		// Cc:      []string{"cc@example.com"},
		// Bcc:     []string{"bcc@example.com"},
		// ReplyTo: "replyto@example.com",
	}
	sent, err := r.client.Emails.Send(params)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("Email sent with ID:", sent.Id)

	return nil
}
