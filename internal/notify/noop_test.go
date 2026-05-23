package notify_test

import (
	"testing"

	"cronwrap/internal/notify"
)

func TestNoopImplementsSender(t *testing.T) {
	// Compile-time check: *NoopSender must satisfy Sender.
	var _ notify.Sender = (*notify.NoopSender)(nil)
}

func TestNoopSenderEmptyByDefault(t *testing.T) {
	s := &notify.NoopSender{}
	if len(s.Messages) != 0 {
		t.Errorf("expected no messages on fresh NoopSender, got %d", len(s.Messages))
	}
}

func TestNoopSenderAccumulatesMessages(t *testing.T) {
	s := &notify.NoopSender{}
	for i := 0; i < 3; i++ {
		if err := s.Send([]string{"op@example.com"}, "alert", "body"); err != nil {
			t.Fatalf("Send[%d] returned error: %v", i, err)
		}
	}
	if len(s.Messages) != 3 {
		t.Errorf("expected 3 messages, got %d", len(s.Messages))
	}
}

func TestNoopSenderResetClearsSlice(t *testing.T) {
	s := &notify.NoopSender{}
	_ = s.Send([]string{"a@b.com"}, "x", "y")
	s.Reset()
	if s.Messages != nil {
		t.Errorf("expected nil slice after Reset, got %v", s.Messages)
	}
}

func TestNoopSenderRecordsFields(t *testing.T) {
	s := &notify.NoopSender{}
	recipients := []string{"alice@example.com", "bob@example.com"}
	subject := "test subject"
	body := "test body"

	if err := s.Send(recipients, subject, body); err != nil {
		t.Fatalf("Send returned unexpected error: %v", err)
	}

	if len(s.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(s.Messages))
	}

	msg := s.Messages[0]
	if msg.Subject != subject {
		t.Errorf("expected subject %q, got %q", subject, msg.Subject)
	}
	if msg.Body != body {
		t.Errorf("expected body %q, got %q", body, msg.Body)
	}
	if len(msg.To) != len(recipients) {
		t.Errorf("expected %d recipients, got %d", len(recipients), len(msg.To))
	}
}
