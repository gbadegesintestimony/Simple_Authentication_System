package utils

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"os"
)

func SendEmail(to, subject, body string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPass := os.Getenv("SMTP_PASS")

	if smtpHost == "" || smtpPort == "" || smtpEmail == "" || smtpPass == "" {
		return fmt.Errorf("SMTP environment variables not set")
	}

	addr := net.JoinHostPort(smtpHost, smtpPort)

	tlsConfig := &tls.Config{
		ServerName: smtpHost,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return err
	}
	defer client.Close()

	auth := smtp.PlainAuth("", smtpEmail, smtpPass, smtpHost)
	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(smtpEmail); err != nil {
		return err
	}
	if err = client.Rcpt(to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	msg := []byte(
		"From: " + smtpEmail + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n\r\n" +
			body,
	)

	if _, err = w.Write(msg); err != nil {
		return err
	}

	if err = w.Close(); err != nil {
		return err
	}

	return client.Quit()
}
