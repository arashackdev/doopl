# doopl Testing Playbook

Complete guide to testing the doopl library and CLI with your own DeepL API token.

## Prerequisites

1. **DeepL API Key**: Get one at https://www.deepl.com/pro-api (free tier available)
   - Free keys end in `:fx`
   - Pro keys don't have this suffix
   
2. **Go 1.24+**:
   ```bash
   go version  # Should be >= go1.24
   ```

3. **Task runner**:
   ```bash
   task --version  # Already required by the project
   ```

---

## Quick Start (5 minutes)

### 1. Set your API key

```bash
# Option A: Export as environment variable
export DEEPL_AUTH_KEY="your-deepl-api-key-here"

# Option B: Create a .env.local file (gitignored)
echo "DEEPL_AUTH_KEY=your-deepl-api-key-here" > .env.local
source .env.local
```

### 2. Run the full test suite

```bash
# Run all local tests (unit + integration, no real API calls)
task test

# Run with verbose output
go test ./... -v -race

# Run specific package
go test ./pkg/deepl -v
```

### 3. Build and test the CLI

```bash
# Build CLI
task cli:build

# Verify it works (uses your DEEPL_AUTH_KEY)
./bin/doopl --help

# Quick translate test
./bin/doopl translate "Hello, world!" --target-lang DE
```

---

## Full Testing Matrix

### Phase 1: Unit Tests (No API Required)

These tests use mocked HTTP responses — your real API key is not needed.

```bash
# All unit tests
task test

# With coverage report
go test ./... -cover

# With detailed stats
go test ./... -race -cover -v

# Check for race conditions (important!)
go test -race ./...

# Specific test only
go test ./pkg/deepl -run TestTranslateText -v
```

**What's tested:**
- Request/response marshaling
- Error handling and sentinel errors
- Language detection
- Functional options
- Glossary operations

---

### Phase 2: CLI Tests (With Your API)

These use your real DeepL account. We recommend starting with a free tier.

#### 2.1: Authentication

```bash
# Test that your key is detected correctly
./bin/doopl doctor

# Output shows:
# - Auth key status (✓ configured)
# - Free vs Pro (based on key suffix)
# - Endpoint detected
```

#### 2.2: Basic Translation

```bash
# Translate a single phrase
./bin/doopl translate "Good morning" --target-lang ES

# Translate multiple phrases (from stdin)
echo -e "Hello\nWorld" | ./bin/doopl translate --target-lang FR

# Translate from file
cat > /tmp/test.txt << 'EOF'
The quick brown fox
Jumps over the lazy dog
EOF

./bin/doopl translate --input-file /tmp/test.txt --target-lang DE

# Specify source language (auto-detection is default)
./bin/doopl translate "Hola" --target-lang EN --source-lang ES

# Test formality
./bin/doopl translate "Hello" --target-lang DE --formality formal
./bin/doopl translate "Hello" --target-lang DE --formality informal
```

#### 2.3: Languages & Usage

```bash
# List all available target languages
./bin/doopl languages --resource translate

# List source languages (for translation)
./bin/doopl languages --resource translate --source

# Check your API usage
./bin/doopl usage

# Output shows:
# - Characters used this month
# - Monthly character limit
# - Remaining quota
```

#### 2.4: JSON Output (for CI/CD integration)

```bash
# Get structured JSON output
./bin/doopl translate "Hello" --target-lang ES --output json

# Useful for parsing in scripts
./bin/doopl usage --output json | jq '.character_count'
```

---

### Phase 3: MCP Server Tests

The MCP server exposes doopl to Claude and other AI clients.

#### 3.1: Build and Verify

```bash
# Build the MCP server
task mcp:build

# Test the binary exists and runs
./bin/doopl-mcp --help
```

#### 3.2: Add to Claude Code (Optional)

If you have Claude Code set up:

1. Build the MCP server:
   ```bash
   task mcp:build
   ```

2. Add to `.claude/settings.json`:
   ```json
   {
     "mcpServers": {
       "doopl": {
         "command": "/full/path/to/doopl/bin/doopl-mcp",
         "args": ["serve"],
         "env": { "DEEPL_AUTH_KEY": "your-key" }
       }
     }
   }
   ```

3. Restart Claude Code, new tools appear

4. In Claude Code, try:
   ```
   You: "Translate 'Hello' to German, Spanish, and French"
   Claude: [calls translate tool, returns results]
   ```

---

### Phase 4: Advanced Testing

#### 4.1: Glossary Operations

```bash
# Create a glossary
./bin/doopl glossary create \
  --name "my-tech-terms" \
  --source EN \
  --target DE \
  --entries "API:Schnittstelle,backend:Backend,frontend:Frontend"

# This prints a glossary ID. Keep it for the next command.

# List your glossaries
./bin/doopl glossary list

# Translate using a glossary
./bin/doopl translate "The API provides a backend interface" \
  --target-lang DE \
  --glossary-id <your-glossary-id>
  
# Delete when done
./bin/doopl glossary delete --glossary-id <your-glossary-id>
```

#### 4.2: Document Translation

```bash
# Create a test document
cat > /tmp/sample.txt << 'EOF'
This is a sample document.
It will be translated to Spanish.
EOF

# Upload for translation
./bin/doopl document upload \
  --input-file /tmp/sample.txt \
  --target-lang ES

# Output includes document_id and key. Use these to check status:
./bin/doopl document status \
  --document-id <id> \
  --document-key <key>

# Once done=true, download
./bin/doopl document download \
  --document-id <id> \
  --document-key <key> \
  --output-file /tmp/sample_es.txt

# Clean up
./bin/doopl document delete \
  --document-id <id> \
  --document-key <key>
```

#### 4.3: Rephrase (Write API)

```bash
# Rephrase text to a formal tone
./bin/doopl rephrase "Hey, how's it going?" \
  --target-lang EN \
  --tone formal

# Try different tones
./bin/doopl rephrase "The project is finished" --target-lang EN --tone friendly
./bin/doopl rephrase "The project is finished" --target-lang EN --tone formal
./bin/doopl rephrase "The project is finished" --target-lang EN --tone casual
```

---

## CI/CD Pipeline Testing

### Phase 5: Verify CI Will Pass

```bash
# Run the exact CI checks locally
task ci

# This runs:
# 1. Format check (gofmt -l)
# 2. go vet
# 3. revive linter
# 4. Generate check (converters up-to-date)
# 5. Unit tests with race detector
# 6. Build verification
```

### Phase 6: Test Release Process

```bash
# Simulate GoReleaser (doesn't push anywhere)
task release:snapshot

# This creates a local release in ./dist/
# with binaries for all platforms:
# - Linux (amd64, arm64)
# - macOS (amd64, arm64)
# - Windows (amd64)

# Verify outputs exist
ls -lh dist/

# Clean up
rm -rf dist/
```

---

## Quota & Cost Management

### Monitor Your Usage

```bash
# Check remaining quota before and after tests
./bin/doopl usage

# For free tier:
# - 500,000 characters/month
# - If you hit limit, wait until next month or upgrade
```

### Cost Estimation

Testing the main features (translate, languages, usage):
- ~1,000–5,000 characters per full test run
- Free tier: 500,000 chars/month → plenty for development

Document & glossary operations use more quota:
- Each document upload/download = character count of file
- Glossaries are free (no character cost)

---

## Troubleshooting

### "API key must not be empty"

```bash
# Check your environment variable
echo $DEEPL_AUTH_KEY

# If empty, set it:
export DEEPL_AUTH_KEY="your-key"

# Or verify .env.local exists
cat .env.local
```

### "401 Unauthorized"

- Invalid or expired API key
- Verify at https://www.deepl.com/pro-api → API key section
- Keys can be regenerated if compromised

### "429 Too Many Requests"

- You've hit the rate limit
- The client automatically retries with backoff
- If persistent, wait a minute and retry

### "Quota Exceeded"

- You've used all monthly characters (free tier: 500k)
- Wait until next calendar month
- Or upgrade to pro tier

### Tests pass locally but fail in CI

```bash
# Check exact CI environment
task tools  # Install pinned tool versions

# Run exact CI task
task ci

# If it fails, check:
# 1. Go version mismatch: go version
# 2. Generated files stale: task generate
# 3. Formatting: task fmt
# 4. Dependencies: go mod tidy
```

---

## Test Coverage

View test coverage for each package:

```bash
# Coverage report
go test ./... -cover

# Detailed HTML report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

**Current targets:**
- `pkg/deepl`: ~85% (main client, all methods)
- `pkg/model`: 100% (simple types, no logic)
- `cmd/doopl`: ~70% (CLI output formatting is hard to test)

---

## Integration with Your Workflow

### Option 1: Local Development

```bash
# Watch mode (if you use entr or similar)
ls *.go cmd/doopl/*.go | entr task test
```

### Option 2: Pre-commit Hook

Create `.git/hooks/pre-commit`:

```bash
#!/bin/bash
set -e

echo "Running tests..."
task test

echo "Running linter..."
task lint

echo "✓ Pre-commit checks passed"
```

### Option 3: GitHub Pull Request

Merges to main automatically:
1. Run `task ci` (all tests)
2. Tag with next version (if main branch)
3. Build release binaries
4. Create GitHub release

---

## Next Steps

- [ ] Set DEEPL_AUTH_KEY environment variable
- [ ] Run `task test` to verify setup
- [ ] Build CLI: `task cli:build`
- [ ] Test a translation: `./bin/doopl translate "Hello" --target-lang ES`
- [ ] Check usage: `./bin/doopl usage`
- [ ] Run full CI: `task ci`
- [ ] Explore other commands: `./bin/doopl --help`

---

## Questions?

- **DeepL API Docs**: https://developers.deepl.com/docs/
- **doopl Godoc**: https://pkg.go.dev/github.com/arashackdev/doopl
- **GitHub Issues**: https://github.com/arashackdev/doopl/issues

Happy testing! 🚀
