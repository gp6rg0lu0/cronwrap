package history

import (
	"os"
	"testing"
	"time"
)

func tempStore(t *testing.T) (*Store, func()) {
	t.Helper()
	f, err := os.CreateTemp("", "cronwrap-*.db")
	if err != nil {
		t.Fatalf("create temp db: %v", err)
	}
	f.Close()

	store, err := NewStore(f.Name())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return store, func() {
		store.Close()
		os.Remove(f.Name())
	}
}

func TestSaveAndList(t *testing.T) {
	store, cleanup := tempStore(t)
	defer cleanup()

	now := time.Now().UTC().Truncate(time.Second)
	rec := RunRecord{
		JobName:   "backup",
		StartedAt: now,
		EndedAt:   now.Add(2 * time.Second),
		ExitCode:  0,
		Output:    "done",
		Success:   true,
	}

	id, err := store.Save(rec)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if id <= 0 {
		t.Errorf("expected positive id, got %d", id)
	}

	records, err := store.List("backup", 10)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].JobName != "backup" {
		t.Errorf("unexpected job name: %s", records[0].JobName)
	}
	if !records[0].Success {
		t.Error("expected success=true")
	}
}

func TestListLimit(t *testing.T) {
	store, cleanup := tempStore(t)
	defer cleanup()

	now := time.Now().UTC()
	for i := 0; i < 5; i++ {
		_, err := store.Save(RunRecord{
			JobName:   "ping",
			StartedAt: now,
			EndedAt:   now.Add(time.Second),
			ExitCode:  0,
			Output:    "",
			Success:   true,
		})
		if err != nil {
			t.Fatalf("Save: %v", err)
		}
	}

	records, err := store.List("ping", 3)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(records) != 3 {
		t.Errorf("expected 3 records with limit, got %d", len(records))
	}
}

func TestListUnknownJob(t *testing.T) {
	store, cleanup := tempStore(t)
	defer cleanup()

	records, err := store.List("nonexistent", 10)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("expected 0 records, got %d", len(records))
	}
}
