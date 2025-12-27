package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/resend/resend-go/v2"
)

func SendEmail(to, subject, body string) error {
	apiKey := os.Getenv("RESEND_API_KEY")

	if apiKey == "" {
		log.Printf("RESEND_API_KEY not set in environment")
		return fmt.Errorf("RESEND_API_KEY not set")
	}

	log.Printf("Sending email to: %s via Resend", to)

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "Password Reset <onboarding@resend.dev>",
		To:      []string{to},
		Subject: subject,
		Text:    body,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		log.Printf("Failed to send email via Resend: %v", err)
		return err
	}

	log.Printf(" Email sent successfully!")
	log.Printf(" Email ID: %s", sent.Id)
	return nil
}
