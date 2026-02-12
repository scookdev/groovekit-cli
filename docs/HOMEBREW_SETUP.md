# Homebrew Setup Guide

## One-Time Setup

### 1. Create the Homebrew Tap Repository

On GitHub, create a new repository named `homebrew-groovekit` (must start with `homebrew-`).

```bash
# Initialize it as an empty repo
gh repo create scookdev/homebrew-groovekit --public
```

### 2. Set Up GitHub Token for Homebrew

GoReleaser needs a GitHub token to update your Homebrew tap automatically.

1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Click "Generate new token (classic)"
3. Name it: "GoReleaser Homebrew Tap"
4. Select scopes:
   - ✅ `repo` (full control)
5. Generate token and copy it

### 3. Add Token to Repository Secrets

In your `groovekit-cli` repo:

1. Go to Settings → Secrets and variables → Actions
2. Click "New repository secret"
3. Name: `HOMEBREW_TAP_GITHUB_TOKEN`
4. Value: paste the token from step 2
5. Save

### 4. Commit and Push the GoReleaser Config

```bash
cd /Users/stevecook/src/groovekit-cli
git add .goreleaser.yaml .github/workflows/release.yml main.go cmd/version.go
git commit -m "Add GoReleaser and Homebrew tap support"
git push
```

## Creating a Release

Now, whenever you want to release a new version:

```bash
# 1. Tag the version
git tag v1.0.0

# 2. Push the tag (this triggers the workflow)
git push origin v1.0.0
```

**That's it!** The GitHub Action will automatically:
- Build binaries for Mac (Intel + ARM) and Linux
- Create a GitHub release with downloadable binaries
- Update your Homebrew tap with the new formula
- Generate release notes from commits

## Users Can Then Install Via:

```bash
# One-time tap
brew tap scookdev/groovekit

# Install
brew install groovekit

# Or in one command
brew install scookdev/groovekit/groovekit
```

## Testing Locally (Before Release)

To test GoReleaser locally without publishing:

```bash
# Install goreleaser
brew install goreleaser

# Test the build (creates snapshot without publishing)
goreleaser release --snapshot --clean

# Test output will be in ./dist/
```

## Version Numbering

Use semantic versioning:
- `v1.0.0` - Major release
- `v1.1.0` - New features (backwards compatible)
- `v1.0.1` - Bug fixes

## Updating After Release

GoReleaser automatically updates the formula, but if you need to manually update:

```bash
# Clone your tap
git clone https://github.com/scookdev/homebrew-groovekit

# Edit Formula/groovekit.rb
# Change url and sha256

# Commit and push
git commit -am "Update to v1.0.1"
git push
```
