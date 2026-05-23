package history

import (
	"context"

	"github.com/yourorg/cronwrap/internal/runner"
)

// RecordResult converts a runner.Result into a Run and persists it via Store.
func RecordResult(ctx context.Context, s *Store, res *runner.Result) error {
	run := Run{
		JobName:   res.JobName,
		Command:   res.Command,
		StartedAt: res.StartedAt,
		Duration:  res.Duration,
		ExitCode:  res.ExitCode,
		Success:   res.Success,
		Stdout:    res.Stdout,
		Stderr:    res.Stderr,
	}
	return s.Save(ctx, run)
}
