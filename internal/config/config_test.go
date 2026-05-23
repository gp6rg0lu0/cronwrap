package config_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/yourorg/cronwrap/internal/config"
)

func writeConfig(t *testing.T, v any) string {
	t.Helper()
	f, err := os.CreateTemp("", "cronwrap-cfg-*.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	if err := json.NewEncoder(f).Encode(v); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoadValid(t *testing.T) {
	path := writeConfig(t, map[string]any{
		"smtp_host": "mail.example.com",
		"smtp_port": 587,
		"smtp_from": "alerts@example.com",
		"db_path":   "/tmp/cw.db",
		"jobs": []map[string]any{
			{"name": "backup", "command": "/usr/bin/backup.sh", "timeout": time.Minute},
		},
	})
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.SMTPHost != "mail.example.com" {
		t.Errorf("smtp_host = %q", cfg.SMTPHost)
	}
	if len(cfg.Jobs) != 1 || cfg.Jobs[0].Name != "backup" {
		t.Errorf("unexpected jobs: %+v", cfg.Jobs)
	}
}

func TestLoadDefaults(t *testing.T) {
	path := writeConfig(t, map[string]any{
		"jobs": []map[string]any{
			{"name": "ping", "command": "ping -c1 localhost"},
		},
	})
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DBPath != "cronwrap.db" {
		t.Errorf("default db_path not set, got %q", cfg.DBPath)
	}
	if cfg.SMTPPort != 25 {
		t.Errorf("default smtp_port not set, got %d", cfg.SMTPPort)
	}
}

func TestLoadMissingCommand(t *testing.T) {
	path := writeConfig(t, map[string]any{
		"jobs": []map[string]any{{"name": "bad"}},
	})
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for missing command")
	}
}

func TestLoadDuplicateName(t *testing.T) {
	path := writeConfig(t, map[string]any{
		"jobs": []map[string]any{
			{"name": "dup", "command": "echo a"},
			{"name": "dup", "command": "echo b"},
		},
	})
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for duplicate job name")
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/cronwrap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
