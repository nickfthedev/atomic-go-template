package mail

import (
	"atomic-go-template/internal/config"
	"errors"
	"os"
)

// TODO: Print to Console Provider

type Service interface {
	Send(to, subject, body string) error
}

func NewMailProvider(c config.Mail) (Service, error) {
	switch c.MailProvider {
	case config.MailProviderResend:
		return NewResendService(os.Getenv("RESEND_API_KEY"))
	case config.MailProviderConsole:
		return NewConsoleService()
	// Add more cases for future providers here
	default:
		return nil, ErrUnsupportedMailProvider
	}
}

// ErrUnsupportedMailProvider is returned when an unsupported mail provider is specified
var ErrUnsupportedMailProvider = errors.New("unsupported mail provider")
