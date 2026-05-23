package history

import (
	"testing"
	"time"
)

func TestSummarizeNoRuns(t *testing.T) {
	s := tempStore(t)

	sum, err := s.Summarize("no-such-job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum.TotalRuns != 0 {
		t.Errorf("expected 0 runs, got %d", sum.TotalRuns)
	}
	if sum.JobName != "no-such-job" {
		t.Errorf("expected job name 'no-such-job', got %q", sum.JobName)
	}
}

func TestSummarizeSingleSuccess(t *testing.T) {
	s := tempStore(t)

	rec := Record{
		JobName:   "backup",
		StartedAt: time.Now().Add(-5 * time.Minute),
		Duration:  10 * time.Second,
		ExitCode:  0,
		Success:   true,
		Output:    "done",
	}
	if err := s.Save(rec); err != nil {
		t.Fatalf("save: %v", err)
	}

	sum, err := s.Summarize("backup")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum.TotalRuns != 1 {
		t.Errorf("expected 1 run, got %d", sum.TotalRuns)
	}
	if sum.SuccessCount != 1 {
		t.Errorf("expected 1 success, got %d", sum.SuccessCount)
	}
	if sum.FailureCount != 0 {
		t.Errorf("expected 0 failures, got %d", sum.FailureCount)
	}
	if sum.LastStatus != "success" {
		t.Errorf("expected last status 'success', got %q", sum.LastStatus)
	}
	if sum.AvgDuration != 10*time.Second {
		t.Errorf("expected avg duration 10s, got %v", sum.AvgDuration)
	}
}

func TestSummarizeMixedResults(t *testing.T) {
	s := tempStore(t)

	now := time.Now()
	records := []Record{
		{JobName: "sync", StartedAt: now.Add(-10 * time.Minute), Duration: 2 * time.Second, ExitCode: 0, Success: true, Output: "ok"},
		{JobName: "sync", StartedAt: now.Add(-5 * time.Minute), Duration: 4 * time.Second, ExitCode: 1, Success: false, Output: "err"},
		{JobName: "sync", StartedAt: now.Add(-1 * time.Minute), Duration: 6 * time.Second, ExitCode: 0, Success: true, Output: "ok"},
	}
	for _, r := range records {
		if err := s.Save(r); err != nil {
			t.Fatalf("save: %v", err)
		}
	}

	sum, err := s.Summarize("sync")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum.TotalRuns != 3 {
		t.Errorf("expected 3 runs, got %d", sum.TotalRuns)
	}
	if sum.SuccessCount != 2 {
		t.Errorf("expected 2 successes, got %d", sum.SuccessCount)
	}
	if sum.FailureCount != 1 {
		t.Errorf("expected 1 failure, got %d", sum.FailureCount)
	}
	if sum.LastStatus != "success" {
		t.Errorf("expected last status 'success', got %q", sum.LastStatus)
	}
	// avg = (2+4+6)/3 = 4s
	if sum.AvgDuration != 4*time.Second {
		t.Errorf("expected avg duration 4s, got %v", sum.AvgDuration)
	}
}
