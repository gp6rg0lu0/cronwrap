package history

import (
	"testing"
	"time"
)

func TestSummarizeNoRuns(t *testing.T) {
	st := tempStore(t)

	sum, err := st.Summarize("nonexistent-job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sum.TotalRuns != 0 {
		t.Errorf("expected 0 runs, got %d", sum.TotalRuns)
	}
	if sum.JobName != "nonexistent-job" {
		t.Errorf("expected job name to be preserved, got %q", sum.JobName)
	}
}

func TestSummarizeSingleSuccess(t *testing.T) {
	st := tempStore(t)

	now := time.Now().UTC()
	err := st.Save(RunRecord{
		JobName:   "backup",
		Status:    "success",
		StartedAt: now,
		Duration:  2500,
	})
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	sum, err := st.Summarize("backup")
	if err != nil {
		t.Fatalf("Summarize: %v", err)
	}

	if sum.TotalRuns != 1 {
		t.Errorf("expected 1 run, got %d", sum.TotalRuns)
	}
	if sum.Successes != 1 {
		t.Errorf("expected 1 success, got %d", sum.Successes)
	}
	if sum.Failures != 0 {
		t.Errorf("expected 0 failures, got %d", sum.Failures)
	}
	if sum.LastStatus != "success" {
		t.Errorf("expected last status 'success', got %q", sum.LastStatus)
	}
}

func TestSummarizeMixedResults(t *testing.T) {
	st := tempStore(t)

	base := time.Now().UTC()
	records := []RunRecord{
		{JobName: "sync", Status: "success", StartedAt: base.Add(-2 * time.Hour), Duration: 100},
		{JobName: "sync", Status: "failure", StartedAt: base.Add(-1 * time.Hour), Duration: 50},
		{JobName: "sync", Status: "success", StartedAt: base, Duration: 120},
	}

	for _, r := range records {
		if err := st.Save(r); err != nil {
			t.Fatalf("Save: %v", err)
		}
	}

	sum, err := st.Summarize("sync")
	if err != nil {
		t.Fatalf("Summarize: %v", err)
	}

	if sum.TotalRuns != 3 {
		t.Errorf("expected 3 runs, got %d", sum.TotalRuns)
	}
	if sum.Successes != 2 {
		t.Errorf("expected 2 successes, got %d", sum.Successes)
	}
	if sum.Failures != 1 {
		t.Errorf("expected 1 failure, got %d", sum.Failures)
	}
	if sum.LastStatus != "success" {
		t.Errorf("expected last status 'success', got %q", sum.LastStatus)
	}

	summaries, err := st.SummarizeAll()
	if err != nil {
		t.Fatalf("SummarizeAll: %v", err)
	}
	if len(summaries) != 1 {
		t.Errorf("expected 1 summary, got %d", len(summaries))
	}
}
