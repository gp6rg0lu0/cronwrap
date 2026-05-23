package alert_test

import (
	"net/smtp"
	"strings"
	"testing"

	"github.com/yourorg/cronwrap/internal/alert"
)

func TestNotifySkipsWhenNoRecipients(t *testing.T) {
	n := alert.New(alert.Config{
		SMTPHost: "localhost",
		SMTPPort: 25,
		From:     "cronwrap@example.com",
		To:       nil,
	})

	// Should return nil without attempting to send.
	if err := n.Notify("backup", 1, "some output"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNotifySendsEmail(t *testing.T) {
	var capturedAddr string
	var capturedFrom string
	var capturedTo []string
	var capturedMsg string

	n := alert.New(alert.Config{
		SMTPHost: "smtp.example.com",
		SMTPPort: 587,
		From:     "cronwrap@example.com",
		To:       []string{"ops@example.com"},
	})

	// Inject a fake send function via the exported test hook.
	n.SetSendFn(func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		capturedAddr = addr
		capturedFrom = from
		capturedTo = to
		capturedMsg = string(msg)
		return nil
	})

	if err := n.Notify("db-backup", 2, "disk full"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedAddr != "smtp.example.com:587" {
		t.Errorf("unexpected addr: %s", capturedAddr)
	}
	if capturedFrom != "cronwrap@example.com" {
		t.Errorf("unexpected from: %s", capturedFrom)
	}
	if len(capturedTo) != 1 || capturedTo[0] != "ops@example.com" {
		t.Errorf("unexpected to: %v", capturedTo)
	}
	if !strings.Contains(capturedMsg, "db-backup") {
		t.Error("message should contain job name")
	}
	if !strings.Contains(capturedMsg, "disk full") {
		t.Error("message should contain output")
	}
	if !strings.Contains(capturedMsg, "Exit Code: 2") {
		t.Error("message should contain exit code")
	}
}

func TestNotifyEmptyOutputFallback(t *testing.T) {
	var capturedMsg string

	n := alert.New(alert.Config{
		SMTPHost: "localhost",
		SMTPPort: 25,
		From:     "a@b.com",
		To:       []string{"c@d.com"},
	})
	n.SetSendFn(func(_ string, _ smtp.Auth, _ string, _ []string, msg []byte) error {
		capturedMsg = string(msg)
		return nil
	})

	_ = n.Notify("cleanup", 1, "   ")
	if !strings.Contains(capturedMsg, "(no output)") {
		t.Error("expected (no output) fallback text")
	}
}
