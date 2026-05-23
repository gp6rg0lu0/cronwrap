// Package notify provides notification backends for cronwrap alerts.
package notify

import (
	"fmt"
	"net/smtp"
	"strings"
)

// SMTPConfig holds the configuration for an SMTP server.
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// Sender is the interface for sending notifications.
type Sender interface {
	Send(to []string, subject, body string) error
}

// SMTPSender sends email notifications via SMTP.
type SMTPSender struct {
	cfg SMTPConfig
}

// NewSMTPSender creates a new SMTPSender with the given config.
func NewSMTPSender(cfg SMTPConfig) *SMTPSender {
	return &SMTPSender{cfg: cfg}
}

// Send delivers an email to the given recipients.
func (s *SMTPSender) Send(to []string, subject, body string) error {
	if len(to) == 0 {
		return fmt.Errorf("notify: no recipients specified")
	}

	header := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n",
		s.cfg.From,
		strings.Join(to, ", "),
		subject,
	)
	message := []byte(header + body)

	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	var auth smtp.Auth
	if s.cfg.Username != "" {
		auth = smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
	}

	return smtp.SendMail(addr, auth, s.cfg.From, to, message)
}
