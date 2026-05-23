package notify

// NoopSender is a Sender implementation that discards all messages.
// It is useful for testing and for configurations where alerting is disabled.
type NoopSender struct {
	// Messages records every Send call for inspection in tests.
	Messages []SentMessage
}

// SentMessage captures the arguments of a single Send invocation.
type SentMessage struct {
	To      []string
	Subject string
	Body    string
}

// Send records the message without delivering it.
func (n *NoopSender) Send(to []string, subject, body string) error {
	n.Messages = append(n.Messages, SentMessage{
		To:      to,
		Subject: subject,
		Body:    body,
	})
	return nil
}

// Reset clears all recorded messages.
func (n *NoopSender) Reset() {
	n.Messages = nil
}
