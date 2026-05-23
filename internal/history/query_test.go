package history

import (
	"testing"
	"time"
)

func TestSummarizeNoRuns(t *testing.T) {
	s := tempStore(t)

	_, err := s.Summarize("ghost-job")
	if err == nil {
		t.Fatal("expected error for unknown job, got nil")
	}
}

func TestSummarizeSingleSuccess(t *testing.T) {
	s := tempStore(t)

	err := s.Save(Run{
		JobName:    "backup",
		StartedAt:  time.Now(),
		DurationMs: 1200,
		ExitCode:   0,
		Output:     "done",
	})
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	sum, err := s.Summarize("backup")
	if err != nil {
		t.Fatalf("Summarize: %v", err)
	}

	if sum.TotalRuns != 1 {
		t.Errorf("TotalRuns: want 1, got %d", sum.TotalRuns)
	}
	if sum.SuccessRuns != 1 {
		t.Errorf("SuccessRuns: want 1, got %d", sum.SuccessRuns)
	}
	if sum.FailureRuns != 0 {
		t.Errorf("FailureRuns: want 0, got %d", sum.FailureRuns)
	}
	if sum.LastStatus != "success" {
		t.Errorf("LastStatus: want success, got %s", sum.LastStatus)
	}
	if sum.JobName != "backup" {
		t.Errorf("JobName: want backup, got %s", sum.JobName)
	}
}

func TestSummarizeMixedResults(t *testing.T) {
	s := tempStore(t)

	runs := []Run{
		{JobName: "sync", StartedAt: time.Now().Add(-3 * time.Minute), DurationMs: 500, ExitCode: 0, Output: "ok"},
		{JobName: "sync", StartedAt: time.Now().Add(-2 * time.Minute), DurationMs: 700, ExitCode: 1, Output: "err"},
		{JobName: "sync", StartedAt: time.Now().Add(-1 * time.Minute), DurationMs: 600, ExitCode: 0, Output: "ok"},
	}
	for _, r := range runs {
		if err := s.Save(r); err != nil {
			t.Fatalf("Save: %v", err)
		}
	}

	sum, err := s.Summarize("sync")
	if err != nil {
		t.Fatalf("Summarize: %v", err)
	}

	if sum.TotalRuns != 3 {
		t.Errorf("TotalRuns: want 3, got %d", sum.TotalRuns)
	}
	if sum.SuccessRuns != 2 {
		t.Errorf("SuccessRuns: want 2, got %d", sum.SuccessRuns)
	}
	if sum.FailureRuns != 1 {
		t.Errorf("FailureRuns: want 1, got %d", sum.FailureRuns)
	}
	if sum.LastStatus != "success" {
		t.Errorf("LastStatus: want success, got %s", sum.LastStatus)
	}

	wantAvg := (500.0 + 700.0 + 600.0) / 3.0 / 1000.0
	if diff := sum.AvgDuration - wantAvg; diff > 0.001 || diff < -0.001 {
		t.Errorf("AvgDuration: want %.4f, got %.4f", wantAvg, sum.AvgDuration)
	}
}
