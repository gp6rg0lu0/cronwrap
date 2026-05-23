// Package alert provides failure notification capabilities for cronwrap.
package alert

import (
	"fmt"
	"net/smtp"
	"strings"
)

// Config holds configuration for sending alert emails.
type Config struct {
	SMTPHost string
	SMTPPort int
	From     string
	To       []string
	Username string
	Password string
}

// Notifier sends alerts when cron jobs fail.
type Notifier struct {
	cfg    Config
	sendFn func(addr string, a smtp.Auth, from string, to []string, msg []byte) error
}

// New creates a new Notifier with the given configuration.
func New(cfg Config) *Notifier {
	return &Notifier{
		cfg:    cfg,
		sendFn: smtp.SendMail,
	}
}

// Notify sends a failure alert email for the given job.
func (n *Notifier) Notify(jobName string, exitCode int, output string) error {
	if len(n.cfg.To) == 0 {
		return nil
	}

	subject := fmt.Sprintf("[cronwrap] FAILED: %s (exit %d)", jobName, exitCode)
	body := buildBody(jobName, exitCode, output)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		n.cfg.From,
		strings.Join(n.cfg.To, ", "),
		subject,
		body,
	))

	addr := fmt.Sprintf("%s:%d", n.cfg.SMTPHost, n.cfg.SMTPPort)
	var auth smtp.Auth
	if n.cfg.Username != "" {
		auth = smtp.PlainAuth("", n.cfg.Username, n.cfg.Password, n.cfg.SMTPHost)
	}

	return n.sendFn(addr, auth, n.cfg.From, n.cfg.To, msg)
}

func buildBody(jobName string, exitCode int, output string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Job:       %s\n", jobName))
	sb.WriteString(fmt.Sprintf("Exit Code: %d\n", exitCode))
	sb.WriteString("\n--- Output ---\n")
	if strings.TrimSpace(output) == "" {
		sb.WriteString("(no output)\n")
	} else {
		sb.WriteString(output)
	}
	return sb.String()
}
