# Changelog

All notable changes to the GrooveKit CLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
