// Package notify provides pluggable notification backends for cronwrap.
//
// # Overview
//
// cronwrap can alert operators when a job fails. The notify package
// abstracts the delivery mechanism behind the Sender interface so that
// different backends (SMTP, webhooks, etc.) can be swapped without
// changing the alert logic.
//
// # Backends
//
// SMTPSender — delivers messages via a standard SMTP relay. Configure
// it with an SMTPConfig that specifies the host, port, and optional
// credentials.
//
// NoopSender — silently discards every message. Intended for use in
// tests and in environments where alerting is explicitly disabled.
//
// # Usage
//
//	var s notify.Sender = notify.NewSMTPSender(notify.SMTPConfig{
//		Host: "smtp.example.com",
//		Port: 587,
//		From: "cronwrap@example.com",
//	})
//	_ = s.Send([]string{"ops@example.com"}, "Job failed", body)
package notify
