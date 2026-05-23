package history

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// RunRecord represents a single cron job execution record.
type RunRecord struct {
	ID        int64
	JobName   string
	StartedAt time.Time
	EndedAt   time.Time
	ExitCode  int
	Output    string
	Success   bool
}

// Store manages the run history database.
type Store struct {
	db *sql.DB
}

// NewStore opens (or creates) the SQLite database at the given path.
func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

// Close closes the underlying database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// Save persists a RunRecord to the database.
func (s *Store) Save(r RunRecord) (int64, error) {
	res, err := s.db.Exec(
		`INSERT INTO runs (job_name, started_at, ended_at, exit_code, output, success)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		r.JobName, r.StartedAt.UTC(), r.EndedAt.UTC(), r.ExitCode, r.Output, r.Success,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// List returns the most recent n records for the given job name.
func (s *Store) List(jobName string, limit int) ([]RunRecord, error) {
	rows, err := s.db.Query(
		`SELECT id, job_name, started_at, ended_at, exit_code, output, success
		 FROM runs WHERE job_name = ? ORDER BY started_at DESC LIMIT ?`,
		jobName, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []RunRecord
	for rows.Next() {
		var r RunRecord
		if err := rows.Scan(&r.ID, &r.JobName, &r.StartedAt, &r.EndedAt, &r.ExitCode, &r.Output, &r.Success); err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	return records, rows.Err()
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS runs (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			job_name   TEXT NOT NULL,
			started_at DATETIME NOT NULL,
			ended_at   DATETIME NOT NULL,
			exit_code  INTEGER NOT NULL,
			output     TEXT NOT NULL,
			success    BOOLEAN NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_runs_job_name ON runs(job_name);
	`)
	return err
}
