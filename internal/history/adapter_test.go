package history_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourorg/cronwrap/internal/history"
	"github.com/yourorg/cronwrap/internal/runner"
)

func TestRecordResult(t *testing.T) {
	s := tempStore(t)
	ctx := context.Background()

	res := &runner.Result{
		JobName:   "backup",
		Command:   "pg_dump",
		Args:      []string{"-Fc", "mydb"},
		Stdout:    "dump ok",
		Stderr:    "",
		ExitCode:  0,
		StartedAt: time.Now().UTC(),
		Duration:  2 * time.Second,
		Success:   true,
	}

	if err := history.RecordResult(ctx, s, res); err != nil {
		t.Fatalf("RecordResult: %v", err)
	}

	runs, err := s.List(ctx, "backup", 10)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(runs))
	}

	got := runs[0]
	if got.JobName != "backup" {
		t.Errorf("job name: want backup, got %q", got.JobName)
	}
	if got.ExitCode != 0 {
		t.Errorf("exit code: want 0, got %d", got.ExitCode)
	}
	if !got.Success {
		t.Errorf("expected success")
	}
	if got.Stdout != "dump ok" {
		t.Errorf("stdout: want 'dump ok', got %q", got.Stdout)
	}
}

func TestRecordResultFailure(t *testing.T) {
	s := tempStore(t)
	ctx := context.Background()

	res := &runner.Result{
		JobName:   "cleanup",
		Command:   "rm",
		ExitCode:  1,
		StartedAt: time.Now().UTC(),
		Duration:  100 * time.Millisecond,
		Success:   false,
		Stderr:    "permission denied",
	}

	if err := history.RecordResult(ctx, s, res); err != nil {
		t.Fatalf("RecordResult: %v", err)
	}

	runs, err := s.List(ctx, "cleanup", 5)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(runs) != 1 || runs[0].Success {
		t.Errorf("expected one failed run")
	}
}
