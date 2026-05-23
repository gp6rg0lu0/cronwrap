// Package alert implements failure notification for cronwrap jobs.
//
// When a monitored cron job exits with a non-zero status, the alert package
// can send an email to one or more recipients describing the failure,
// including the job name, exit code, and captured output.
//
// Basic usage:
//
//	n := alert.New(alert.Config{
//		SMTPHost: "smtp.example.com",
//		SMTPPort: 587,
//		From:     "cronwrap@example.com",
//		To:       []string{"ops@example.com"},
//		Username: "user",
//		Password: "secret",
//	})
//
//	if result.ExitCode != 0 {
//		_ = n.Notify(jobName, result.ExitCode, result.Output)
//	}
//
// If the To field is empty, Notify is a no-op and returns nil.
package alert
