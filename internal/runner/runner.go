// Package runner executes shell commands and captures their output.
package runner

import (
	"bytes"
	"context"
	"os/exec"
	"time"
)

// Result holds the outcome of a command execution.
type Result struct {
	JobName   string
	Command   string
	Args      []string
	Stdout    string
	Stderr    string
	ExitCode  int
	StartedAt time.Time
	Duration  time.Duration
	Success   bool
}

// Runner executes commands and returns a Result.
type Runner struct {
	timeout time.Duration
}

// New creates a Runner with an optional timeout (0 means no timeout).
func New(timeout time.Duration) *Runner {
	return &Runner{timeout: timeout}
}

// Run executes the given command and captures stdout/stderr.
func (r *Runner) Run(ctx context.Context, jobName, command string, args ...string) (*Result, error) {
	if r.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.timeout)
		defer cancel()
	}

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	result := &Result{
		JobName:   jobName,
		Command:   command,
		Args:      args,
		Stdout:    stdout.String(),
		Stderr:    stderr.String(),
		StartedAt: start,
		Duration:  duration,
		Success:   err == nil,
	}

	if cmd.ProcessState != nil {
		result.ExitCode = cmd.ProcessState.ExitCode()
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
			return result, nil
		}
		return result, err
	}

	return result, nil
}
