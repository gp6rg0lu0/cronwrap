// Package config handles loading and validating cronwrap configuration.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Job describes a single cron job managed by cronwrap.
type Job struct {
	Name       string        `json:"name"`
	Command    string        `json:"command"`
	Timeout    time.Duration `json:"timeout"`
	AlertEmail []string      `json:"alert_email"`
}

// Config is the top-level cronwrap configuration.
type Config struct {
	SMTPHost string `json:"smtp_host"`
	SMTPPort int    `json:"smtp_port"`
	SMTPFrom string `json:"smtp_from"`
	DBPath   string `json:"db_path"`
	Jobs     []Job  `json:"jobs"`
}

// Load reads and parses a JSON config file from the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// validate checks that required fields are present.
func (c *Config) validate() error {
	if c.DBPath == "" {
		c.DBPath = "cronwrap.db"
	}
	if c.SMTPPort == 0 {
		c.SMTPPort = 25
	}
	seen := make(map[string]bool)
	for i, j := range c.Jobs {
		if j.Name == "" {
			return fmt.Errorf("config: job[%d] missing name", i)
		}
		if j.Command == "" {
			return fmt.Errorf("config: job %q missing command", j.Name)
		}
		if seen[j.Name] {
			return fmt.Errorf("config: duplicate job name %q", j.Name)
		}
		seen[j.Name] = true
	}
	return nil
}
