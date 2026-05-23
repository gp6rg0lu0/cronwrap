// Package runner provides functionality for executing shell commands
// as part of cron job wrapping.
//
// It captures stdout and stderr independently, records timing information,
// and normalises exit codes so callers can decide how to handle failures
// without parsing raw errors.
//
// Basic usage:
//
//	r := runner.New(30 * time.Second) // zero timeout means no limit
//	res, err := r.Run(ctx, "my-job", "pg_dump", "-Fc", "mydb")
//	if err != nil {
//		// hard error — command could not be started
//	}
//	if !res.Success {
//		// command ran but exited non-zero
//	}
package runner
