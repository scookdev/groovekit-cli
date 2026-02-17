# Changelog

All notable changes to the GrooveKit CLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.3.0] - 2026-02-18

### Added
- **DNS Record Monitoring**: Full CRUD operations for DNS record monitors
- Monitor DNS records for unexpected changes across A, AAAA, MX, CNAME, TXT, and NS record types
- Mismatch detection with color-coded status (green = values match, red = values don't match)
- Expected vs current value comparison in `dns show`
- Incident history tracking for DNS change events
- **Domain Expiration Monitoring**: Full CRUD operations for domain expiration monitors
- Track domain registration expiry with configurable alert thresholds
- Color-coded days-until-expiration display (green/yellow/red based on proximity to expiry)
- Three-tier alerting with warning, urgent, and critical thresholds (defaults: 30/14/7 days)
- Registrar and registrar URL display in domain details
- Incident history tracking for domain expiration events

### Commands
- `groovekit dns` - Manage DNS record monitors
  - `list` - View all DNS monitors with mismatch status
  - `show <id>` - Display monitor details including expected and current values
  - `create` - Add new DNS monitor (`--name`, `--domain`, `--type`, `--expected`)
  - `update <id>` - Modify existing DNS monitor
  - `delete <id>` - Remove DNS monitor (with `--force` to skip confirmation)
  - `pause <id>` - Temporarily pause DNS monitoring
  - `resume <id>` - Resume paused DNS monitoring
  - `incidents <id>` - View DNS monitor incident history
- `groovekit domains` - Manage domain expiration monitors
  - `list` - View all domain monitors with days until expiration
  - `show <id>` - Display domain details including registrar and thresholds
  - `create` - Add new domain monitor (`--name`, `--domain`, threshold flags)
  - `update <id>` - Modify existing domain monitor
  - `delete <id>` - Remove domain monitor (with `--force` to skip confirmation)
  - `pause <id>` - Temporarily pause domain monitoring
  - `resume <id>` - Resume paused domain monitoring
  - `incidents <id>` - View domain monitor incident history

### Technical
- Added `DnsMonitor`, `CreateDnsMonitorRequest`, and `UpdateDnsMonitorRequest` types
- Added `DomainMonitor`, `CreateDomainMonitorRequest`, and `UpdateDomainMonitorRequest` types
- Added client methods for DNS and domain CRUD, incidents, pause, and resume
- Short ID prefix matching (Docker-style) for both DNS and domain monitors

## [1.2.0] - 2026-02-12

### Added
- **SSL Certificate Monitoring**: Full CRUD operations for SSL certificate monitors
- Track SSL certificate expiration dates with configurable warning thresholds
- Color-coded certificate expiration display (green/yellow/red based on days remaining)
- Multiple alert threshold levels (warning, urgent, critical)
- Domain and port monitoring support
- Incident history tracking for certificate issues
- Certificate details including issuer, subject, and expiration information

### Commands
- `groovekit certs` - Manage SSL certificate monitors
  - `list` - View all SSL certificate monitors with expiration status
  - `show <id>` - Display detailed certificate information
  - `create` - Add new SSL certificate monitor
  - `update <id>` - Modify existing certificate monitor
  - `delete <id>` - Remove SSL certificate monitor
  - `pause <id>` - Temporarily pause certificate monitoring
  - `resume <id>` - Resume paused certificate monitoring
  - `incidents <id>` - View certificate incident history

### Improved
- Comprehensive test coverage for jobs, monitors, and certs commands
- Added table-driven tests for helper functions
- Flag validation tests for all commands
- Subcommand registration verification tests

### Technical
- Added `SslMonitor` type with full certificate tracking fields
- Implemented `CreateSslMonitorRequest` and `UpdateSslMonitorRequest` types
- Added client methods: `CreateCert`, `UpdateCert`, `DeleteCert`, `ListCertIncidents`
- Port field properly typed as integer for API compatibility

## [1.1.0] - 2026-02-09

### Added
- **CI/CD Integration Support**: Headless authentication via `GROOVEKIT_TOKEN` environment variable
- API key authentication support for long-lived tokens in automated workflows
- Compatible with GitHub Actions, GitLab CI, CircleCI, and all major CI/CD platforms
- Automatic fallback from JWT to API key authentication

### Use Cases
- Block deployments when monitors are down
- Automate monitor creation for new services
- Check monitor status as part of health checks
- View incidents and check history in CI pipelines

### Example
```bash
# Set environment variable (e.g., in GitHub Actions secrets)
export GROOVEKIT_TOKEN="gk_..."

# Use CLI without interactive login
groovekit monitors list
groovekit monitors incidents <id>
groovekit account show
```

See the [CI/CD Integration Guide](https://groovekit.io/dashboard/docs#ci-cd-integration) for detailed setup instructions.

## [1.0.0] - 2026-02-06

### Added
- Initial release of GrooveKit CLI
- Full CRUD operations for cron job monitors (create, list, show, update, delete)
- Full CRUD operations for API monitors (create, list, show, update, delete)
- Pause and resume commands for jobs and monitors
- Check history viewing for both jobs and monitors
- Incident history tracking showing downtime periods
- Account information display with usage bars
- JSON output support for all commands with `--json` flag
- Short ID support (Docker-style prefix matching)
- Loading spinners for better UX
- Authentication with API token storage
- Homebrew installation support

### Commands
- `groovekit auth login` - Authenticate with GrooveKit
- `groovekit auth logout` - Sign out
- `groovekit jobs` - Manage cron job monitors
  - `list`, `show`, `create`, `update`, `delete`
  - `pause`, `resume` - Quick status changes
  - `incidents` - View downtime history
- `groovekit monitors` - Manage API monitors
  - `list`, `show`, `create`, `update`, `delete`
  - `pause`, `resume` - Quick status changes
  - `incidents` - View downtime history
- `groovekit checks` - View check and ping history
- `groovekit account show` - View account details, subscription, and usage
- `groovekit version` - Show CLI version information
