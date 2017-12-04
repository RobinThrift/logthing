package main

import (
	"fmt"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"os"
	"strings"
	"time"
)

type SendGridError struct {
	response *rest.Response
}

func (sg *SendGridError) Error() string {
	return sg.response.Body
}

func sendGridSender(sender string, recipient string, body string) error {
	if len(body) == 0 {
		return nil
	}

	now := time.Now().Format("Monday, 02-Jan-06 15:04:05 MST")
	from := mail.NewEmail("Logthing", sender)
	subject := "Logs from " + now
	to := mail.NewEmail(recipient, recipient)
	htmlBody := strings.Replace(body, "\n", "<br />", -1)
	message := mail.NewSingleEmail(from, subject, to, body, htmlBody)
	client := sendgrid.NewSendClient(os.Getenv("LT_SENDGRID_API_KEY"))
	response, err := client.Send(message)

	if err != nil {
		return err
	}

	if response.StatusCode != 202 {
		return &SendGridError{response}
	}

	fmt.Println(now, "- Email sent.")

	return nil
}
