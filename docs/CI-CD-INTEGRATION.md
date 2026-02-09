# CI/CD Integration Guide

GrooveKit CLI supports headless authentication via environment variables, making it perfect for CI/CD pipelines.

## Quick Start

### 1. Get Your Access Token

After logging in locally with `groovekit login`, find your token in the config file:

```bash
cat ~/.groovekit/config.json
```

Look for the `access_token` field. Copy this value.

### 2. Add Token to CI/CD Secrets

#### GitHub Actions
1. Go to your repository → **Settings** → **Secrets and variables** → **Actions**
2. Click **New repository secret**
3. Name: `GROOVEKIT_TOKEN`
4. Value: Paste your access token
5. Click **Add secret**

#### GitLab CI
Add to your project's CI/CD variables:
```
Settings → CI/CD → Variables → Add variable
Key: GROOVEKIT_TOKEN
Value: <your-token>
Flags: ✅ Masked, ✅ Protected
```

#### CircleCI
```
Project Settings → Environment Variables → Add Variable
Name: GROOVEKIT_TOKEN
Value: <your-token>
```

### 3. Use in Your Pipeline

The CLI automatically detects the `GROOVEKIT_TOKEN` environment variable:

```yaml
# .github/workflows/deploy.yml
- name: Check Monitor Health
  env:
    GROOVEKIT_TOKEN: ${{ secrets.GROOVEKIT_TOKEN }}
  run: |
    groovekit monitors list
```

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `GROOVEKIT_TOKEN` | Your access token for authentication | Yes (in CI/CD) |
| `GROOVEKIT_API_URL` | Custom API endpoint (default: `https://api.groovekit.io`) | No |

**Precedence:** Environment variables take precedence over config file values.

## Common CI/CD Use Cases

### Pre-Deployment Health Checks

Block deployments if monitors are experiencing issues:

```yaml
- name: Pre-Deploy Health Check
  env:
    GROOVEKIT_TOKEN: ${{ secrets.GROOVEKIT_TOKEN }}
  run: |
    if groovekit monitors list --json | jq -e '.api_monitors[] | select(.down == true)'; then
      echo "🚨 Monitors are down! Blocking deployment."
      exit 1
    fi
```

### Automated Monitor Creation

Create monitors when deploying new services:

```yaml
- name: Setup Production Monitoring
  env:
    GROOVEKIT_TOKEN: ${{ secrets.GROOVEKIT_TOKEN }}
  run: |
    groovekit monitors create \
      --name "Production API - ${{ github.sha }}" \
      --url https://api.example.com/health \
      --interval 5 \
      --expected-status-codes 200,201
```

### Monitor Management as Code

Store monitor configurations in your repo:

```bash
# scripts/setup-monitors.sh
#!/bin/bash
set -e

groovekit monitors create \
  --name "API Gateway" \
  --url https://api.example.com/health \
  --interval 5

groovekit monitors create \
  --name "User Service" \
  --url https://users.example.com/health \
  --interval 5

groovekit monitors create \
  --name "Payment Service" \
  --url https://payments.example.com/health \
  --interval 3 \
  --timeout 10
```

### Incident Reporting

Get incident summaries in your pipeline:

```yaml
- name: Check Recent Incidents
  env:
    GROOVEKIT_TOKEN: ${{ secrets.GROOVEKIT_TOKEN }}
  run: |
    # Get monitor ID from previous step or hardcode
    MONITOR_ID="abc123..."

    echo "## Recent Incidents" >> $GITHUB_STEP_SUMMARY
    groovekit monitors incidents $MONITOR_ID >> $GITHUB_STEP_SUMMARY
```

## Complete GitHub Actions Example

See [example-ci-monitoring.yml](../.github/workflows/example-ci-monitoring.yml) for a complete working example.

## Security Best Practices

### ✅ DO:
- Store tokens in CI/CD secrets (never commit to repo)
- Use masked/protected secrets in GitLab
- Rotate tokens periodically
- Use different tokens for different environments (if managing multiple accounts)
- Limit token access with read-only permissions when possible (future feature)

### ❌ DON'T:
- Commit tokens to version control
- Share tokens in plain text
- Use the same token across different organizations
- Log tokens in CI output (they should be automatically masked)

## Troubleshooting

### Authentication Errors

```
Error: not authenticated. Please run 'groovekit login' first
```

**Solution:** Ensure `GROOVEKIT_TOKEN` is set and contains a valid token.

```yaml
env:
  GROOVEKIT_TOKEN: ${{ secrets.GROOVEKIT_TOKEN }}
```

### Token Not Found in Secrets

**Solution:** Verify the secret is added to your repository/organization and the name matches exactly (case-sensitive).

### API Connection Issues

```
Error: failed to connect to API
```

**Solution:** Check if you need to set a custom API URL:

```yaml
env:
  GROOVEKIT_TOKEN: ${{ secrets.GROOVEKIT_TOKEN }}
  GROOVEKIT_API_URL: https://your-custom-api.example.com
```

## Installation in CI/CD

### Using Go
```yaml
- name: Install GrooveKit
  run: |
    go install github.com/scookdev/groovekit-cli@latest
    # Or build from source:
    # git clone https://github.com/scookdev/groovekit-cli
    # cd groovekit-cli && go build -o /usr/local/bin/groovekit
```

### Using Pre-built Binaries (Future)
```yaml
- name: Install GrooveKit
  run: |
    curl -sSL https://install.groovekit.io | sh
```

## Need Help?

- 📚 [Main Documentation](../README.md)
- 🐛 [Report Issues](https://github.com/scookdev/groovekit-cli/issues)
- 💬 [Discussions](https://github.com/scookdev/groovekit-cli/discussions)
