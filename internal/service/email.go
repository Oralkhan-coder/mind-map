package service

import (
	"fmt"
	"net/smtp"

	"github.com/Oralkhan-coder/mind-map/config"
)

type EmailService struct {
	ec *config.EmailConfig
}

func NewEmailService(ec *config.EmailConfig) *EmailService {
	return &EmailService{ec: ec}
}

func (es *EmailService) SendHTML(subject string, htmlBody string, to string) error {
	toSlice := []string{to}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"%s"+
		"%s\r\n", es.ec.From, to, subject, mime, htmlBody)

	auth := smtp.PlainAuth("", es.ec.From, es.ec.Password, es.ec.SmtpHost)
	addr := es.ec.SmtpHost + ":" + es.ec.SmtpPort

	err := smtp.SendMail(addr, auth, es.ec.From, toSlice, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send HTML email: %w", err)
	}

	return nil
}
