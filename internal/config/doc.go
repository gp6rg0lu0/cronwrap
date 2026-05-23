// Package config provides types and helpers for loading cronwrap's JSON
// configuration file.
//
// # Configuration File
//
// cronwrap expects a JSON file (default: cronwrap.json) with the following
// top-level fields:
//
//	{
//	  "smtp_host": "mail.example.com",   // SMTP server hostname
//	  "smtp_port": 587,                  // SMTP port (default: 25)
//	  "smtp_from": "alerts@example.com", // envelope From address
//	  "db_path":   "cronwrap.db",        // SQLite history database path
//	  "jobs": [
//	    {
//	      "name":        "backup",             // unique job identifier
//	      "command":     "/usr/bin/backup.sh", // shell command to execute
//	      "timeout":     3600000000000,        // nanoseconds (1 h)
//	      "alert_email": ["ops@example.com"]   // recipients on failure
//	    }
//	  ]
//	}
//
// Use [Load] to read and validate a configuration file.
package config
