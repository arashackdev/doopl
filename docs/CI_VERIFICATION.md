# CI/CD Verification Checklist

Quick reference to verify all CI/CD improvements are working correctly.

## ✅ Before You Merge

```bash
# 1. Run full CI locally
task ci
# Should see: ✓ fmt:check ✓ vet ✓ lint ✓ generate:check ✓ test ✓ build

# 2. Verify builds work
task cli:build
./bin/doopl --version  # or --help

task mcp:build
./bin/doopl-mcp --version  # or --help

# 3. Check git is clean
git status
```

## ✅ Testing the CI Pipeline (Post-Merge to main)

### Step 1: Monitor GitHub Actions

```bash
# After pushing to main, check GitHub Actions
open https://github.com/arashackdev/doopl/actions

# You should see:
# - "CI" workflow running (tests, lint, build)
# - After success, "tag" job creates next version
# - "Release" workflow triggered by new tag
```

### Step 2: Verify Tagging (CI Automation)

```bash
# Once CI passes and tags created:
git fetch --all --tags

# Check new tags exist
git tag | tail -3

# Example output:
# v0.0.2
# latest
# v0.0.1
```

### Step 3: Verify Release Creation

```bash
# Release should be created on GitHub
open https://github.com/arashackdev/doopl/releases/latest

# You should see:
# - Version tag (v0.0.2)
# - Release notes (auto-generated)
# - Binaries for all platforms (Linux/macOS/Windows × amd64/arm64)
```

## ✅ Testing Local Release Simulation

```bash
# Test that GoReleaser config is correct (doesn't actually release)
task release:check
task release:snapshot

# Check binaries were built
ls -lh dist/

# Clean up
rm -rf dist/
```

## ✅ Key Improvements Verified

| Component | Before | After | Status |
|-----------|--------|-------|--------|
| Version tagging | Manual (error-prone) | Automatic (GitHub Actions) | ✅ |
| Binary builds | CLI only | CLI + MCP server | ✅ |
| Build optimization | None | Static + AMD64v3 | ✅ |
| Module verification | Missing | Added to CI | ✅ |
| Linting paths | Fragile hardcoding | Pattern-based | ✅ |
| Cache consistency | Go version mismatch | Versioned cache key | ✅ |
| Committed binaries | ❌ (in git) | ✅ (gitignored) | ✅ |

## ✅ Testing Release to Public

When ready to release v0.0.2 (or next version):

```bash
# 1. Ensure main is clean and all tests pass
git checkout main
task ci

# 2. Push to GitHub (triggers entire pipeline automatically)
git push origin main

# 3. GitHub Actions will:
#    - Run tests (≈5 min)
#    - Create next version tag (v0.0.2)
#    - Create release with binaries
#    - Update 'latest' tag

# 4. Verify on GitHub
open https://github.com/arashackdev/doopl/releases

# 5. Users can now install:
#    go install github.com/arashackdev/doopl/cmd/doopl@v0.0.2
#    go install github.com/arashackdev/doopl/cmd/doopl-mcp@v0.0.2
```

## ✅ Troubleshooting CI

### CI tests fail but work locally

```bash
# Ensure Go version matches
go version
# Should match matrix in .github/workflows/ci.yml

# Regenerate converters
task generate

# Reinstall tools at exact versions
task tools

# Run full CI
task ci
```

### Tags don't auto-create

```bash
# Check GitHub Actions permissions
# Settings > Actions > General > Workflow permissions
# Should be: "Read and write permissions"

# Check secret access (if needed)
# Settings > Secrets and variables > Repository secrets
# (Our setup doesn't require secrets for tagging)

# Verify CI job completes successfully
# Actions tab > latest CI run > Check "test" job passes
# If "tag" job doesn't appear or fails, check GitHub logs
```

### Release workflow doesn't trigger

```bash
# Ensure tag matches v* pattern
git tag | grep "^v"

# Check release.yml is watching correct tag pattern
cat .github/workflows/release.yml | grep "tags:"
# Should show: - v*

# If still not triggering, verify GoReleaser config
task release:check
```

## ✅ Binaries Verification

After release, verify both binaries work:

```bash
# For doopl CLI
doopl --help
doopl translate "Hello" --target-lang ES

# For MCP server
doopl-mcp serve --help
# (or use in Claude Code via MCP servers setting)
```

---

**Questions?** Check `TESTING_PLAYBOOK.md` for full testing guide, or `REVIEW_SUMMARY.md` for architecture details.
