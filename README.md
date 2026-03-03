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

### Man Pages

Man pages are included with the Homebrew installation and available for all commands and subcommands:

```bash
man groovekit
man groovekit-jobs
man groovekit-apis-create
# etc.
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
# List all job monitors
groovekit jobs list

# Create a new job monitor
groovekit jobs create --name "Daily Backup" --interval 1440 --grace-period 5

# Show job monitor details
groovekit jobs show <job-id>

# Update a job monitor
groovekit jobs update <job-id> --name "Updated Name" --interval 720

# Pause/resume a job monitor
groovekit jobs pause <job-id>
groovekit jobs resume <job-id>

# View incident history
groovekit jobs incidents <job-id>

# Delete a job monitor
groovekit jobs delete <job-id>
```

**Job intervals are in minutes.** Example: `--interval 1440` = check every 24 hours.

### API Monitoring

```bash
# List all api monitors
groovekit apis list

# Create a new api monitor
groovekit apis create \
  --name "Production API" \
  --url https://api.example.com/health \
  --interval 60 \
  --method GET

# Show api monitor details
groovekit apis show <monitor-id>

# Update an api monitor
groovekit apis update <monitor-id> --interval 30 --timeout 10

# Pause/resume an api monitor
groovekit apis pause <monitor-id>
groovekit apis resume <monitor-id>

# View incident history
groovekit apis incidents <monitor-id>

# Delete a monitor
groovekit apis delete <monitor-id>
```

**API Monitor intervals are in minutes.** Example: `--interval 60` = check every hour.

### SSL Certificate Monitoring

```bash
# List all SSL certificate monitors
groovekit certs list

# Create a new SSL certificate monitor
groovekit certs create --name "example.com SSL" --domain example.com --port 443

# Show certificate details
groovekit certs show <cert-id>

# Update a certificate monitor
groovekit certs update <cert-id> --warning-threshold 45 --critical-threshold 14

# Pause/resume a certificate monitor
groovekit certs pause <cert-id>
groovekit certs resume <cert-id>

# View incident history
groovekit certs incidents <cert-id>

# Delete a certificate monitor
groovekit certs delete <cert-id>
```

### Domain Expiration Monitoring

```bash
# List all domain monitors
groovekit domains list

# Create a new domain monitor
groovekit domains create --name "example.com" --domain example.com

# Show domain details
groovekit domains show <domain-id>

# Update a domain monitor
groovekit domains update <domain-id> --warning-threshold 45

# Pause/resume a domain monitor
groovekit domains pause <domain-id>
groovekit domains resume <domain-id>

# View incident history
groovekit domains incidents <domain-id>

# Delete a domain monitor
groovekit domains delete <domain-id>
```

### DNS Record Monitoring

```bash
# List all DNS monitors
groovekit dns list

# Create a new DNS monitor
groovekit dns create \
  --name "Example MX" \
  --domain example.com \
  --type MX \
  --expected mail.example.com

# Show DNS monitor details (including expected vs current values)
groovekit dns show <dns-id>

# Update a DNS monitor
groovekit dns update <dns-id> --expected "new-value.example.com"

# Pause/resume a DNS monitor
groovekit dns pause <dns-id>
groovekit dns resume <dns-id>

# View incident history
groovekit dns incidents <dns-id>

# Delete a DNS monitor
groovekit dns delete <dns-id>
```

Supported DNS record types: `A`, `AAAA`, `MX`, `CNAME`, `TXT`, `NS`

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
- **SSL Certificate Monitoring**: Track certificate expiration with color-coded days remaining and multi-tier alert thresholds
- **Domain Expiration Monitoring**: Monitor domain registration expiry with configurable warning, urgent, and critical thresholds
- **DNS Record Monitoring**: Detect unexpected DNS changes across A, AAAA, MX, CNAME, TXT, and NS record types
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
