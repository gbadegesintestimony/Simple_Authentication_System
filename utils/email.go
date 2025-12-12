package utils

import (
	"crypto/tls"
	"net/smtp"
	"os"
)

func SendEmail(to, subject, body string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpPass := os.Getenv("SMTP_PASS")
	smtpEmail := os.Getenv("SMTP_EMAIL")

	msg := []byte(
		"From: " + smtpEmail + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n" +
			body,
	)

	auth := smtp.PlainAuth("", smtpEmail, smtpPass, smtpHost)

	// 1. Connect
	client, err := smtp.Dial(smtpHost + ":" + smtpPort)
	if err != nil {
		return err
	}

	// 2. Upgrade to TLS (THIS IS THE FIX)
	tlsConfig := &tls.Config{
		ServerName: smtpHost,
	}

	if err = client.StartTLS(tlsConfig); err != nil {
		return err
	}

	// 3. Authenticate
	if err = client.Auth(auth); err != nil {
		return err
	}

	// 4. Send email
	if err = client.Mail(smtpEmail); err != nil {
		return err
	}
	if err = client.Rcpt(to); err != nil {
		return err
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}

	_, err = writer.Write(msg)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	return client.Quit()
}
