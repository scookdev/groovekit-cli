# GrooveKit CLI

Official command-line interface for [GrooveKit](https://groovekit.io) - Monitor your cron jobs and APIs with confidence.

## Installation

### macOS / Linux (Homebrew)

```bash
brew install scookdev/groovekit/groovekit
```

### Windows / Manual Installation

Download the latest release for your platform from the [releases page](https://github.com/scookdev/groovekit-cli/releases/latest).

**Windows:**
1. Download `groovekit_*_Windows_x86_64.zip`
2. Extract and add to your PATH

**Linux:**
```bash
# Download and extract (adjust version and architecture as needed)
curl -L https://github.com/scookdev/groovekit-cli/releases/download/v1.0.0/groovekit_1.0.0_Linux_x86_64.tar.gz | tar xz

# Move to PATH
sudo mv groovekit /usr/local/bin/
```

## Getting Started

### Authentication

```bash
groovekit auth login
```

Enter your GrooveKit email and password. Your credentials are stored securely in `~/.groovekit/config.json`.

### View Account Info

```bash
groovekit account show
```

## Usage

### Cron Job Monitoring

```bash
# List all jobs
groovekit jobs list

# Create a new job
groovekit jobs create --name "Daily Backup" --interval 1440 --grace-period 5

# Show job details
groovekit jobs show <job-id>

# Update a job
groovekit jobs update <job-id> --name "Updated Name" --interval 720

# Pause/resume a job
groovekit jobs pause <job-id>
groovekit jobs resume <job-id>

# View incident history
groovekit jobs incidents <job-id>

# Delete a job
groovekit jobs delete <job-id>
```

**Job intervals are in minutes.** Example: `--interval 1440` = check every 24 hours.

### API Monitoring

```bash
# List all monitors
groovekit monitors list

# Create a new monitor
groovekit monitors create \
  --name "Production API" \
  --url https://api.example.com/health \
  --interval 60 \
  --method GET

# Show monitor details
groovekit monitors show <monitor-id>

# Update a monitor
groovekit monitors update <monitor-id> --interval 30 --timeout 10

# Pause/resume a monitor
groovekit monitors pause <monitor-id>
groovekit monitors resume <monitor-id>

# View incident history
groovekit monitors incidents <monitor-id>

# Delete a monitor
groovekit monitors delete <monitor-id>
```

**Monitor intervals are in minutes.** Example: `--interval 60` = check every hour.

### Check History

```bash
# View recent checks for a monitor
groovekit checks list --monitor <monitor-id>

# View recent pings for a job
groovekit checks list --job <job-id>
```

### JSON Output

All commands support `--json` flag for machine-readable output:

```bash
groovekit jobs list --json
groovekit account show --json
```

## Features

- **Cron Job Monitoring**: Heartbeat ping monitoring with configurable intervals and grace periods
- **API Monitoring**: HTTP endpoint health checks with response time tracking and status code validation
- **Incident Tracking**: View downtime history and recovery times
- **Check History**: Review recent pings and health check results
- **Short IDs**: Docker-style ID prefix matching (use `abc123` instead of full UUID)
- **Account Management**: View subscription details, usage limits, and current usage
- **JSON Output**: Machine-readable output for automation and scripting

## Documentation

For more information about GrooveKit, visit:
- **Website**: https://groovekit.io
- **Features**: https://groovekit.io/features
- **Pricing**: https://groovekit.io/pricing
- **API Documentation**: https://groovekit.io/dashboard/docs

## Support

- **Contact**: https://groovekit.io/contact
- **Issues**: https://github.com/scookdev/groovekit-cli/issues

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

MIT License - see [LICENSE](LICENSE) for details.
