// Command cronwrap is the entry point for the cronwrap job runner.
// It loads configuration, executes the specified job, records history,
// and sends alerts on failure.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/example/cronwrap/internal/alert"
	"github.com/example/cronwrap/internal/config"
	"github.com/example/cronwrap/internal/history"
	"github.com/example/cronwrap/internal/runner"
)

func main() {
	cfgPath := flag.String("config", "cronwrap.toml", "path to configuration file")
	jobName := flag.String("job", "", "name of the job to run (required)")
	flag.Parse()

	if *jobName == "" {
		fmt.Fprintln(os.Stderr, "error: --job flag is required")
		flag.Usage()
		os.Exit(2)
	}

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("cronwrap: failed to load config: %v", err)
	}

	jobCfg, ok := cfg.Jobs[*jobName]
	if !ok {
		log.Fatalf("cronwrap: job %q not found in config", *jobName)
	}

	store, err := history.NewStore(cfg.HistoryDB)
	if err != nil {
		log.Fatalf("cronwrap: failed to open history store: %v", err)
	}
	defer store.Close()

	alerter := alert.New(cfg.SMTP)

	r := runner.New(jobCfg)
	result := r.Run()

	if err := history.RecordResult(store, *jobName, result); err != nil {
		log.Printf("cronwrap: warning: failed to record history: %v", err)
	}

	if !result.Success {
		if err := alerter.Notify(*jobName, result); err != nil {
			log.Printf("cronwrap: warning: failed to send alert: %v", err)
		}
		os.Exit(1)
	}
}
