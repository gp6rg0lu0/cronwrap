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
