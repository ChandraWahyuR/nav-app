package mailer

import (
	"fmt"
	"net/smtp"
	"proyek1/config"
	"strings"
)

type MailInterface interface {
	SendMail(emailReceiver string, subject, templatePath string, data any) error
}

type mailer struct {
	c config.SMTP
}

func NewMail(c config.SMTP) mailer {
	return mailer{
		c: c,
	}
}

func (m *mailer) SendMail(emailReceiver string, subject, templatePath string, data any) error {
	msg := "MIME-version: 1.0;\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\n" +
		"From: " + m.c.SMTP_EMAIL_ADDRESS + "\n" +
		"To: " + emailReceiver + "\n" +
		"Subject: " + subject + "\n\n" +
		templatePath

	auth := smtp.PlainAuth("", m.c.SMTP_EMAIL_ADDRESS, m.c.SMTP_TOKEN_EMAIL, m.c.SMTP_HOST)
	smtpAddr := fmt.Sprintf("%s:%d", m.c.SMTP_HOST, m.c.SMTP_PORT)

	err := smtp.SendMail(smtpAddr, auth, m.c.SMTP_EMAIL_ADDRESS, strings.Split(emailReceiver, ","), []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
