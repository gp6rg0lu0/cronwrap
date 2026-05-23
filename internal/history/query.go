package history

import (
	"fmt"
	"time"
)

// Summary holds aggregated statistics for a job's run history.
type Summary struct {
	JobName      string
	TotalRuns    int
	SuccessCount int
	FailureCount int
	LastRun      time.Time
	LastStatus   string
	AvgDuration  time.Duration
}

// Summarize computes a Summary for the given job from its stored records.
// It returns an error if the job name is empty or no records are found.
func (s *Store) Summarize(jobName string) (Summary, error) {
	if jobName == "" {
		return Summary{}, fmt.Errorf("job name must not be empty")
	}

	records, err := s.List(jobName, 0)
	if err != nil {
		return Summary{}, fmt.Errorf("summarize %q: %w", jobName, err)
	}

	if len(records) == 0 {
		return Summary{JobName: jobName}, nil
	}

	var (
		successes   int
		failures    int
		totalNanos  int64
		latestTime  time.Time
		latestStatus string
	)

	for _, r := range records {
		if r.Success {
			successes++
		} else {
			failures++
		}
		totalNanos += r.Duration.Nanoseconds()

		if r.StartedAt.After(latestTime) {
			latestTime = r.StartedAt
			if r.Success {
				latestStatus = "success"
			} else {
				latestStatus = "failure"
			}
		}
	}

	avg := time.Duration(totalNanos / int64(len(records)))

	return Summary{
		JobName:      jobName,
		TotalRuns:    len(records),
		SuccessCount: successes,
		FailureCount: failures,
		LastRun:      latestTime,
		LastStatus:   latestStatus,
		AvgDuration:  avg,
	}, nil
}
