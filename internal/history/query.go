package history

import (
	"database/sql"
	"fmt"
	"time"
)

// Summary holds aggregated statistics for a single job.
type Summary struct {
	JobName     string
	TotalRuns   int
	SuccessRuns int
	FailureRuns int
	LastRun     time.Time
	LastStatus  string
	AvgDuration float64 // seconds
}

// Summarize returns aggregated run statistics for the named job.
// If the job has no recorded runs, it returns an error.
func (s *Store) Summarize(jobName string) (*Summary, error) {
	query := `
		SELECT
			COUNT(*) AS total,
			SUM(CASE WHEN exit_code = 0 THEN 1 ELSE 0 END) AS successes,
			SUM(CASE WHEN exit_code != 0 THEN 1 ELSE 0 END) AS failures,
			MAX(started_at) AS last_run,
			AVG(duration_ms) / 1000.0 AS avg_duration_sec
		FROM runs
		WHERE job_name = ?
	`

	row := s.db.QueryRow(query, jobName)

	var total, successes, failures int
	var lastRun time.Time
	var avgDuration float64

	err := row.Scan(&total, &successes, &failures, &lastRun, &avgDuration)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("history: no runs found for job %q", jobName)
		}
		return nil, fmt.Errorf("history: summarize query failed: %w", err)
	}

	if total == 0 {
		return nil, fmt.Errorf("history: no runs found for job %q", jobName)
	}

	// Fetch the status of the most recent run.
	var lastStatus string
	statusRow := s.db.QueryRow(
		`SELECT exit_code FROM runs WHERE job_name = ? ORDER BY started_at DESC LIMIT 1`,
		jobName,
	)
	var lastCode int
	if err := statusRow.Scan(&lastCode); err == nil {
		if lastCode == 0 {
			lastStatus = "success"
		} else {
			lastStatus = "failure"
		}
	}

	return &Summary{
		JobName:     jobName,
		TotalRuns:   total,
		SuccessRuns: successes,
		FailureRuns: failures,
		LastRun:     lastRun,
		LastStatus:  lastStatus,
		AvgDuration: avgDuration,
	}, nil
}
