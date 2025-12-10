package utils

import (
	"net/smtp"
	"os"
)

func SendEmail(to string, subject string, body string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpPass := os.Getenv("SMTP_PASS")
	smtpemail := os.Getenv("SMTP_EMAIL")

	auth := smtp.PlainAuth("", smtpemail, smtpPass, smtpHost)
	msg := "From: " + smtpemail + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body + "\r\n"

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpemail, []string{to}, []byte(msg))
}
