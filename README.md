# cronwrap

A lightweight wrapper for cron jobs that captures output, tracks run history, and sends failure alerts.

## Installation

```bash
go install github.com/yourusername/cronwrap@latest
```

## Usage

Wrap any cron job command by prefixing it with `cronwrap`:

```bash
# In your crontab
* * * * * cronwrap --name "backup-job" --alert email@example.com /usr/local/bin/backup.sh
```

### Common Flags

| Flag | Description |
|------|-------------|
| `--name` | Human-readable name for the job |
| `--alert` | Email address to notify on failure |
| `--history` | Number of past runs to retain (default: 50) |
| `--timeout` | Max execution time before job is killed |
| `--log` | Path to write captured output |

### Example

```bash
cronwrap --name "db-cleanup" \
         --alert ops@example.com \
         --timeout 5m \
         --log /var/log/cronwrap/db-cleanup.log \
         /opt/scripts/cleanup.sh
```

Run history and captured output are stored in `~/.cronwrap/history/` by default.

### View Run History

```bash
cronwrap history --name "db-cleanup"
```

## Requirements

- Go 1.21 or later

## License

MIT © [yourusername](https://github.com/yourusername)