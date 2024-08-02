package mail

import (
	"fmt"
)

type ConsoleService struct {
}

func NewConsoleService() (Service, error) {
	return &ConsoleService{}, nil
}

func (r *ConsoleService) Send(to, subject, body string) error {

	fmt.Println("Email sent to:", to)
	fmt.Println("Email subject:", subject)
	fmt.Println("Email body:", body)

	return nil
}
