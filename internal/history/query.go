package history

import (
	"database/sql"
	"time"
)

// Summary holds aggregated statistics for a single job.
type Summary struct {
	JobName    string
	TotalRuns  int
	Successes  int
	Failures   int
	LastRun    time.Time
	LastStatus string
}

// Summarize returns aggregated run statistics for the given job name.
// If no runs exist for the job, it returns a zero-value Summary and no error.
func (s *Store) Summarize(jobName string) (Summary, error) {
	rows, err := s.db.Query(`
		SELECT status, started_at
		FROM runs
		WHERE job_name = ?
		ORDER BY started_at DESC
	`, jobName)
	if err != nil {
		return Summary{}, err
	}
	defer rows.Close()

	var sum Summary
	sum.JobName = jobName

	for rows.Next() {
		var status string
		var startedAt time.Time

		if err := rows.Scan(&status, &startedAt); err != nil {
			return Summary{}, err
		}

		sum.TotalRuns++

		if sum.TotalRuns == 1 {
			sum.LastRun = startedAt
			sum.LastStatus = status
		}

		switch status {
		case "success":
			sum.Successes++
		case "failure":
			sum.Failures++
		}
	}

	if err := rows.Err(); err != nil {
		return Summary{}, err
	}

	return sum, nil
}

// SummarizeAll returns aggregated statistics for every job that has at least
// one recorded run.
func (s *Store) SummarizeAll() ([]Summary, error) {
	rows, err := s.db.Query(`SELECT DISTINCT job_name FROM runs ORDER BY job_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name sql.NullString
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name.String)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	summaries := make([]Summary, 0, len(names))
	for _, name := range names {
		sum, err := s.Summarize(name)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, sum)
	}

	return summaries, nil
}
