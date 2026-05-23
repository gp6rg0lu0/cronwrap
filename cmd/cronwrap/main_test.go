package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// buildBinary compiles the cronwrap binary into a temp directory and
// returns its path. The test is skipped if the build fails.
func buildBinary(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "cronwrap")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Skipf("skipping integration test: build failed: %v\n%s", err, out)
	}
	return bin
}

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "cronwrap-*.toml")
	if err != nil {
		t.Fatalf("create temp config: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp config: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestMainMissingJobFlag(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit when --job is missing")
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 2 {
			t.Errorf("expected exit code 2, got %d", exitErr.ExitCode())
		}
	}
	_ = out
}

func TestMainSuccessfulJob(t *testing.T) {
	bin := buildBinary(t)
	dbPath := filepath.Join(t.TempDir(), "history.db")
	cfg := writeTempConfig(t, `
history_db = "`+dbPath+`"

[jobs.hello]
command = "echo hello"
timeout = "5s"
`)
	cmd := exec.Command(bin, "--config", cfg, "--job", "hello")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected success, got error: %v\noutput: %s", err, out)
	}
}

func TestMainUnknownJob(t *testing.T) {
	bin := buildBinary(t)
	dbPath := filepath.Join(t.TempDir(), "history.db")
	cfg := writeTempConfig(t, `
history_db = "`+dbPath+`"

[jobs.hello]
command = "echo hello"
timeout = "5s"
`)
	cmd := exec.Command(bin, "--config", cfg, "--job", "nonexistent")
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit for unknown job")
	}
}
