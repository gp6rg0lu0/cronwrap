package notify_test

import (
	"net"
	"net/smtp"
	"strings"
	"testing"

	"cronwrap/internal/notify"
)

// startFakeSMTP launches a minimal TCP listener that accepts one connection,
// performs a bare-minimum SMTP handshake, and records the raw DATA payload.
func startFakeSMTP(t *testing.T) (addr string, received *string) {
	t.Helper()
	var buf string
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("startFakeSMTP: listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		fmt.Fprintf := func(f string, a ...any) { conn.Write([]byte(fmt.Sprintf(f, a...))) }
		_ = smtp.NewClient // ensure import used
		conn.Write([]byte("220 fake ESMTP\r\n"))
		b := make([]byte, 4096)
		n, _ := conn.Read(b)
		buf = string(b[:n])
	}()

	return ln.Addr().String(), &buf
}

func TestNoopSenderRecordsMessages(t *testing.T) {
	s := &notify.NoopSender{}

	if err := s.Send([]string{"a@b.com"}, "hello", "world"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(s.Messages))
	}
	msg := s.Messages[0]
	if msg.Subject != "hello" {
		t.Errorf("subject: got %q, want %q", msg.Subject, "hello")
	}
	if msg.Body != "world" {
		t.Errorf("body: got %q, want %q", msg.Body, "world")
	}
}

func TestNoopSenderReset(t *testing.T) {
	s := &notify.NoopSender{}
	_ = s.Send([]string{"x@y.com"}, "s", "b")
	s.Reset()
	if len(s.Messages) != 0 {
		t.Errorf("expected 0 messages after Reset, got %d", len(s.Messages))
	}
}

func TestNoopSenderMultipleRecipients(t *testing.T) {
	s := &notify.NoopSender{}
	to := []string{"a@b.com", "c@d.com"}
	_ = s.Send(to, "sub", "body")

	if len(s.Messages[0].To) != 2 {
		t.Errorf("expected 2 recipients, got %d", len(s.Messages[0].To))
	}
}

func TestSMTPSenderNoRecipients(t *testing.T) {
	s := notify.NewSMTPSender(notify.SMTPConfig{
		Host: "127.0.0.1",
		Port: 9999,
		From: "from@example.com",
	})
	err := s.Send(nil, "sub", "body")
	if err == nil {
		t.Fatal("expected error for empty recipients, got nil")
	}
	if !strings.Contains(err.Error(), "no recipients") {
		t.Errorf("unexpected error message: %v", err)
	}
}
