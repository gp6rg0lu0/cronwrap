// Package history provides persistent storage for cron job run records.
//
// It uses an embedded SQLite database to store execution history including
// job name, start/end times, exit code, captured output, and success status.
//
// Basic usage:
//
//	store, err := history.NewStore("/var/lib/cronwrap/history.db")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer store.Close()
//
//	// Record a completed run
//	_, err = store.Save(history.RunRecord{
//		JobName:   "daily-backup",
//		StartedAt: startTime,
//		EndedAt:   time.Now(),
//		ExitCode:  0,
//		Output:    combinedOutput,
//		Success:   true,
//	})
//
//	// Retrieve the last 20 runs for a job
//	records, err := store.List("daily-backup", 20)
package history
