package runner_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourorg/cronwrap/internal/runner"
)

func TestRunSuccess(t *testing.T) {
	r := runner.New(0)
	res, err := r.Run(context.Background(), "echo-job", "echo", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Success {
		t.Errorf("expected success, got failure")
	}
	if res.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", res.ExitCode)
	}
	if res.Stdout != "hello\n" {
		t.Errorf("expected stdout 'hello\\n', got %q", res.Stdout)
	}
	if res.JobName != "echo-job" {
		t.Errorf("expected job name 'echo-job', got %q", res.JobName)
	}
}

func TestRunFailure(t *testing.T) {
	r := runner.New(0)
	res, err := r.Run(context.Background(), "fail-job", "false")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Success {
		t.Errorf("expected failure, got success")
	}
	if res.ExitCode == 0 {
		t.Errorf("expected non-zero exit code")
	}
}

func TestRunTimeout(t *testing.T) {
	r := runner.New(50 * time.Millisecond)
	res, err := r.Run(context.Background(), "sleep-job", "sleep", "5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Success {
		t.Errorf("expected timeout failure, got success")
	}
}

func TestRunDurationTracked(t *testing.T) {
	r := runner.New(0)
	res, err := r.Run(context.Background(), "dur-job", "echo", "hi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Duration <= 0 {
		t.Errorf("expected positive duration, got %v", res.Duration)
	}
}

func TestRunStderr(t *testing.T) {
	r := runner.New(0)
	res, err := r.Run(context.Background(), "stderr-job", "sh", "-c", "echo error >&2; exit 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Stderr != "error\n" {
		t.Errorf("expected stderr 'error\\n', got %q", res.Stderr)
	}
}
