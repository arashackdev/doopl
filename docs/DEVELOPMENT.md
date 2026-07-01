# doopl Development Guide

Comprehensive guide for contributing to doopl — the idiomatic Go client for DeepL API v3.

## Quick Start

```bash
# Install dependencies
go mod download

# Run full CI suite
task ci

# Build the CLI
task cli:build

# Run tests
task test

# Generate converters (after editing interfaces)
task generate

# Format, vet, lint
task fmt && task vet && task lint
```

## Architecture

doopl enforces a **three-layer separation** to keep concerns isolated and make the library genuinely embeddable:

1. **Wire Format Layer** (`v3/apimodel/`) — exact DeepL API shape
   - JSON types that mirror API requests/responses (snake_case, raw string enums)
   - Never exposed to library consumers
   - Updated when DeepL API changes (via OpenAPI spec drift detection)
   - Example: `TranslateRequest`, `TranslateResponse`

2. **Domain Model Layer** (`pkg/model/`) — idiomatic Go types
   - Public API surface; what library consumers use
   - Go conventions: CamelCase enums, zero values are safe
   - Converters auto-generated from interface declarations
   - Example: `TextResult`, `Formality`, `Language`

3. **CLI Entity Layer** (`cmd/doopl/internal/entity/`) — display types
   - Terminal/JSON output only; never used by library code
   - Shaped for readability and batch operations
   - Separate from model so CLI concerns don't pressure public API
   - Example: `TranslationRow`, `LanguageRow`, `DoctorReport`

**Why this pattern?**
- Changes to the API don't cascade to consumer code
- CLI concerns (column ordering, display formatting) stay isolated
- Converters are generated, never hand-written, so they stay in sync

## Converters

Generated converters map types between layers using `goverter`:

- **Library converters** (`internal/convert/apimodel_to_model.go` → `converter_gen.go`)
  - Wire format → domain model
  - Regenerated with `task generate`

- **CLI converters** (`cmd/doopl/internal/convert/model_to_entity.go` → `converter_gen.go`)
  - Domain model → CLI display entity
  - Same pattern, one layer further out

**After editing converter interfaces, always run `task generate`.**

## Output System

All CLI commands support three output modes via the `--output` flag:

- **text** (default) — Plain text for readability
- **tui** — Rich terminal UI using lipgloss (colors, boxes, formatting)
- **json** — Structured output matching entity/model definitions

The `Formatter` interface in `cmd/doopl/internal/output/formatter.go` provides:
- `FormatTranslations()`, `FormatLanguages()`, `FormatUsage()`, etc.
- Three implementations: `TextFormatter`, `TUIFormatter`, `JSONFormatter`

Commands use `NewFormatter(c.String("output"))` to get the appropriate formatter.

## Adding a New Endpoint

Follow this pattern for each new DeepL API resource (M2–M6 from roadmap):

### 1. Create API wire format
```go
// v3/apimodel/mystuff.go
type MyResourceRequest struct {
    Param string `json:"param"`
}

type MyResourceResponse struct {
    Result string `json:"result"`
}
```

### 2. Create public domain type
```go
// pkg/model/mystuff.go
// MyResource represents ...
type MyResource struct {
    Param string
    Result string
}
```

### 3. Declare converter interface
```go
// internal/convert/apimodel_to_model.go — add method
type APIToModel interface {
    // existing methods...
    MyResourceToModel(src apimodel.MyResourceResponse) model.MyResource
}
```

### 4. Add Client method
```go
// pkg/deepl/mystuff.go
func (c *Client) MyResource(ctx context.Context, options ...Option) (*model.MyResource, error) {
    // Use apiToModel converter for wire format → model conversion
    // Implementation details
}
```

### 5. Add CLI command and entity
```go
// cmd/doopl/internal/entity/myresource.go
type MyResourceRow struct {
    Param string `json:"param"`
    Result string `json:"result"`
}

// cmd/doopl/command_myresource.go
func myResourceCommand() *cli.Command {
    return &cli.Command{
        Name: "myresource",
        Action: func(c *cli.Context) error {
            // fetch, convert, format, output
            formatter := output.NewFormatter(c.String("output"))
            fmt.Print(formatter.FormatMyResources(...))
        },
    }
}

// cmd/doopl/main.go — register in commands slice
```

### 6. Test end-to-end
```bash
task generate  # Regenerate converters
task test      # Verify everything compiles and works
task ci        # Full suite
```

## Testing

### Test Coverage

Tests use in-process `httptest` servers that mock DeepL responses. Current coverage: **48.1%** of statements.

Run tests with:
```bash
go test ./... -v
go test ./... -race           # Check for race conditions
go test ./... -cover          # Full coverage report
go test ./pkg/deepl -v -run TestName  # Single test
```

### Mocking Strategy

- Each endpoint has an `httptest.Server` that mocks the DeepL API response
- No real API key needed for tests
- Example: `translate_test.go` runs against a mock server

### Testing Best Practices

- All public functions should have at least one test
- Test both success and failure paths
- Use `errors.Is()` and `errors.As()` for error checking
- Verify retry behavior, backoff, and context cancellation

## Code Standards

### Godoc

All exported symbols (functions, types, methods, constants) **must** have godoc comments:

```go
// TranslateText translates text into the target language.
// It returns results for each input text, in order.
// If ctx is canceled, TranslateText returns the context error.
func (c *Client) TranslateText(ctx context.Context, texts []string, targetLang string, opts ...TranslateTextOption) ([]model.TextResult, error)
```

- One-line summary describing what it is/does
- Additional detail on inputs, outputs, and notable errors
- No need to repeat what the signature already says

Run `go doc ./...` to verify all exports are documented. The linter `revive` enforces this with the `exported` rule.

### Naming

- Library layer (`pkg/deepl`, `pkg/model`): **CamelCase** everywhere
- Wire format layer (`v3/apimodel`): **snake_case** in JSON tags, reflects API exactly
- CLI layer: CamelCase (Go convention)
- `Functional options` pattern for Client methods: `WithOption()` functions

### Comments

Write comments for **why**, not **what**:
- ✅ `// Retry with exponential backoff per RFC 6797 to handle rate limits` (non-obvious constraint)
- ✗ `// Loop over results` (code already shows this)

If code is hard to understand, refactor instead of commenting. Good names are better than comments.

### No Hand-Written Converters

All type conversions between layers are generated by `goverter`. This ensures:
- Converters stay in sync when types change
- No silent field mismatches
- Edit the interface, run `task generate`, done

## Error Handling

Use sentinel errors with `errors.Is()` and `errors.As()`:

```go
import "errors"

// Define sentinel errors in errors.go
var ErrQuotaExceeded = errors.New("quota exceeded")

// Check with errors.Is()
if errors.Is(err, doopl.ErrQuotaExceeded) {
    // Handle quota specifically
}

// Check API errors with errors.As()
var apiErr *doopl.Error
if errors.As(err, &apiErr) {
    // Handle with apiErr.Code, apiErr.Message
}
```

Never wrap errors carelessly — `errors.Is()` won't match wrapped errors unless you use `errors.Wrap()` or `fmt.Errorf("%w", err)`.

## Building

### CLI
```bash
go build -o doopl ./cmd/doopl
```

### MCP Server (for AI clients)
```bash
go build -o doopl-mcp ./cmd/doopl-mcp
```

### Cross-Platform Release (via GitHub Actions)
GitHub Actions (`release.yml`) automatically builds binaries for:
- Linux: amd64, arm64
- macOS: amd64, arm64
- Windows: amd64

Tag a release and push: `git tag v0.0.1 && git push origin v0.0.1`

## Linting & Formatting

```bash
# Format Go code
go fmt ./...

# Vet (catch suspicious constructs)
go vet ./...

# Lint with revive (fast, focused)
revive ./...

# All together
task ci
```

### Linter Rules (`revive.toml`)

- `exported` — all public symbols must have godoc
- `package-comments` — package must have a comment
- `var-naming` — snake_case in wire format, camelCase elsewhere
- `blank-imports` — no empty imports
- `unreachable-code` — flagged
- `time-naming` — `time.Duration` vars should end in `Duration`

## Release Process

### Version Numbers

Follow semantic versioning:
- **v0.0.1** — initial release
- **v0.x.y** — may have breaking changes (pre-v1.0.0)
- **v1.0.0** — stable API (after M6 complete)

### To Release

1. **Verify all changes are committed and tests pass:**
   ```bash
   task ci  # Local verification
   ```

2. **Tag the release:**
   ```bash
   git tag v0.0.1
   git push origin v0.0.1
   ```

3. **GitHub Actions does the rest:**
   - Builds CLI and MCP server for all platforms
   - Creates a GitHub Release with binaries
   - `pkg.go.dev` automatically indexes within seconds

4. **Users can install via:**
   ```bash
   go get github.com/arashackdev/doopl@v0.0.1
   go install github.com/arashackdev/doopl/cmd/doopl@v0.0.1
   ```

## Troubleshooting

### `task generate` fails
- Check interface declarations in `internal/convert/apimodel_to_model.go` for syntax errors
- Run `go generate -v ./...` for detailed output
- Ensure all goverter directives are correct

### Linting fails
- Missing godoc? `go doc ./...` shows which symbols lack comments
- Naming violation? Check `var-naming` rule in `revive.toml`
- Wire format types should be snake_case; everything else CamelCase

### Tests fail with "connection refused"
- Ensure tests create their own `httptest.Server` (they should)
- No external API key should be required for tests

### CLI flag not recognized
- Check `cmd/doopl/main.go` — flag must be registered in command definition
- Flag name must match the option name (e.g., `--formality` → `WithFormality`)

## Project Files

- **Taskfile.yml** — Task definitions (build, test, lint, generate, release)
- **revive.toml** — Linter configuration
- **go.mod**, **go.sum** — Dependencies (pinned, reproducible)
- **tools.go** — Dev tool versions (goverter, revive)
- **.goreleaser.yml** — Release configuration for multi-platform builds
- **CLAUDE.md** — High-level overview and MCP setup

## Roadmap (M1–M6)

- **M1** (✅ done): Text translation, full parameter parity, three-layer architecture
- **M2** (in progress): Languages & usage endpoints + CLI commands
- **M3**: Document translation (upload, poll, download)
- **M4**: Glossaries (create, list, entries, delete)
- **M5**: Write API (rephrase)
- **M6**: Polish, examples, deepl-mock in CI, GoReleaser, v1.0.0

## Key Principles

- **Idiomatic Go** — functional options, sentinel errors, simple interfaces
- **Full parameter parity** — every DeepL v3 option is available
- **Three-layer isolation** — API changes don't cascade to consumers
- **Genuinely embeddable** — CLI is proof; holds no translation logic
- **Zero hand-written converters** — generated from interfaces, always in sync
- **Production ready** — backoff, retry, context cancellation, concurrent-safe

## Support

- **Bugs**: https://github.com/arashackdev/doopl/issues
- **Discussions**: https://github.com/arashackdev/doopl/discussions
- **API Docs**: https://pkg.go.dev/github.com/arashackdev/doopl
- **CLI Help**: `doopl --help` or `doopl COMMAND --help`
