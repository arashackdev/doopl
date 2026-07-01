# doopl — Claude Setup & Development Guide

## AI Integration (Start Here!)

**Using doopl with Claude or other AI clients?** Start with **[AI Integration Guide](#ai-integration-guide-for-mcp-servers)** below, then return here for development.

## Project Overview

**doopl** is an idiomatic Go client for the DeepL API v3, published as a **library** (`github.com/arashackdev/doopl`), a **CLI** (`cmd/doopl`), and an **MCP server** (`cmd/doopl-mcp`) for AI clients. This is version 0.0.1, with full text translation support and complete godoc documentation for all exported APIs.

**Repository:** https://github.com/arashackdev/doopl

### Key Design Principles

- **Three-layer architecture:** API wire format → public domain model → CLI display entity, each isolated by generated converters
- **Full parameter parity:** `TranslateText` supports all DeepL v3 options (formality, glossary, context, style, tag handling, etc.)
- **Idiomatic Go:** functional options, sentinel errors with `errors.Is`/`errors.As`, free/pro endpoint auto-detection
- **Zero hand-written converters:** mappings generated via `goverter` from interface declarations

## Development Resources

For comprehensive development guide, see **[`docs/DEVELOPMENT.md`](./docs/DEVELOPMENT.md)**. It covers:
- Architecture and three-layer separation
- Adding new endpoints
- Testing strategy
- Code standards and godoc conventions
- Release process

**Additional Documentation (in `docs/` directory):**
- `TESTING_PLAYBOOK.md` — Complete testing guide with your own DeepL API key
- `CI_VERIFICATION.md` — CI/CD pipeline verification and troubleshooting
- `REVIEW_SUMMARY.md` — Code review, architecture assessment, improvements
- `QUICK_TEST.sh` — One-command validation script

> **📌 Documentation Policy:** All docs must go in `docs/` directory. Never scatter `.md` or `.sh` files in the project root. See `.gitignore` for enforcement patterns.

| Tool | Purpose | Config |
|------|---------|--------|
| Go 1.24+ | Language | `go.mod` |
| Task | Task runner | `Taskfile.yml` |
| revive | Linter | `revive.toml` |
| goverter | Code generation | interface declarations in `internal/convert/` + `cmd/doopl/internal/convert/` |
| urfave/cli | CLI framework | `cmd/doopl/main.go` |
| lipgloss | TUI output | `cmd/doopl/internal/output/formatter.go` |

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

# Generate converters (after editing interface declarations)
task generate

# Format, vet, lint
task fmt && task vet && task lint
```

## Project Structure

```
doopl/
├── client.go                    # Client with functional options
├── errors.go                    # Sentinel errors (ErrQuotaExceeded, etc.)
├── translate.go                 # TranslateText + all options
├── translate_test.go            # Tests vs. httptest mock
├── request.go                   # Shared HTTP transport helpers
│
├── model/
│   └── translate.go             # model.TextResult (public domain type)
│
├── internal/
│   ├── apimodel/
│   │   └── translate.go         # apimodel.TranslateRequest/Response (wire format)
│   ├── convert/
│   │   ├── apimodel_to_model.go # goverter interface: APIToModel
│   │   └── converter_gen.go     # GENERATED — DO NOT EDIT
│   └── transport/
│       └── transport.go         # Retry + backoff + auth
│
├── cmd/doopl/
│   ├── main.go                  # CLI entry point, translate subcommand
│   └── internal/
│       ├── entity/
│       │   └── translation.go   # entity.TranslationRow (CLI display type)
│       └── convert/
│           ├── model_to_entity.go # goverter interface: ModelToEntity
│           └── converter_gen.go  # GENERATED — DO NOT EDIT
│
├── docs/
│   └── scope-and-checklist.md   # Current status + M2–M6 roadmap
│
├── Taskfile.yml
├── revive.toml
├── go.mod
├── go.sum
└── README.md
```

## Godoc & API Documentation

All public APIs (functions, types, methods) carry godoc comments. After publication to pkg.go.dev, documentation is available at:

- [doopl library](https://pkg.go.dev/github.com/arashackdev/doopl)
- [doopl CLI](https://pkg.go.dev/github.com/arashackdev/doopl/cmd/doopl)

Before v0.0.1 release, verify:

```bash
go doc ./...  # Verify all exported symbols have godoc
```

Or run `task lint` — revive enforces the `exported` rule, which fails if any exported symbol lacks a comment.

## Adding New Features

The repo follows a mechanical pattern for each new resource (M2–M6 from the roadmap):

### 1. Create the API wire format
```go
// internal/apimodel/languages.go
type LanguagesResponse struct {
    Languages []Language `json:"languages"`
}

type Language struct {
    Language string `json:"language"`
    Name     string `json:"name"`
}
```

### 2. Create the public domain type
```go
// model/languages.go
// Language represents a supported language for translation.
type Language struct {
    Code string
    Name string
}
```

### 3. Declare the converter interface
```go
// internal/convert/apimodel_to_model.go
//go:generate go run github.com/jmattheis/goverter/cmd/goverter

type APIToModel interface {
    // ... existing conversions ...
    LanguagesToModel(src []apimodel.Language) []model.Language
}
```

### 4. Add the Client method
```go
// client.go
// Languages returns the list of supported languages for the given resource.
func (c *Client) Languages(ctx context.Context, resource string) ([]model.Language, error) {
    // implementation
}
```

### 5. Add CLI entity and command
```go
// cmd/doopl/internal/entity/language.go
type LanguageRow struct {
    Code string `json:"code"`
    Name  string `json:"name"`
}

// cmd/doopl/main.go — add Languages command
```

### 6. Test end-to-end
```bash
task generate  # Regenerate converters
task test      # Verify everything compiles
```

## Environment & Secrets

### Required
- `DEEPL_AUTH_KEY`: Your DeepL API key (free keys end in `:fx`)

### Optional
- `DEEPL_SERVER_URL`: Override API endpoint (useful for testing with `deepl-mock`)

### Testing
Tests use an in-process `httptest` server that mocks DeepL responses. No real API key needed for tests.

## Code Quality Checks

All checks must pass before committing:

```bash
task ci  # Runs all of the below
```

Individual checks:
- `task fmt` — gofmt check
- `task vet` — go vet
- `task lint` — revive (fast, focused linter)
- `task generate:check` — verify converters are up-to-date
- `task test` — unit tests with race detector

Linter config is in `revive.toml`. Current rules:

- `exported`: all public symbols must have godoc
- `package-comments`: package must have a comment
- `var-naming`: snake_case in internal/apimodel (wire format), camelCase elsewhere
- `blank-imports`: no empty imports
- `unreachable-code`: flagged
- `time-naming`: `time.Duration` vars should end in `Duration`

## AI Integration Guide for MCP Servers

doopl provides an **MCP (Model Context Protocol) server** that exposes DeepL translation capabilities to Claude, Claude Desktop, and any MCP-compatible AI client.

### What This Means

Instead of writing Go code to call doopl, Claude can **directly call translation tools** in your conversations:

```
You: "Translate 'Hello, world!' to German and French"
Claude: [calls translate tool] "Hallo, Welt!" / "Bonjour, le monde!"
```

### Quick Setup

1. **Build the MCP server:**
   ```bash
   cd /path/to/doopl
   task mcp:build
   # Creates: ./bin/doopl-mcp
   ```

2. **Add to Claude Code settings** (`.claude/settings.json`):
   ```json
   {
     "mcpServers": {
       "doopl": {
         "command": "/full/path/to/bin/doopl-mcp",
         "args": ["serve"],
         "env": { "DEEPL_AUTH_KEY": "your-deepl-key" }
       }
     }
   }
   ```

3. **Restart Claude** — new tools will appear.

### Available Tools

| Tool | Capability |
|------|-----------|
| `translate` | Translate text; supports formality, glossaries, context, style, tag handling |
| `languages` | List supported languages for translate/document/glossary/write |
| `usage` | Check API quota and current usage |

### Use Cases

1. **Document localization** — translate README.md to multiple languages
2. **Multilingual support** — reply to customers in their language
3. **Language availability checks** — verify DeepL supports a language
4. **Quota-aware workflows** — check usage before batch operations

### Complete Reference

See **[`.claude/AI-INTEGRATION.md`](./.claude/AI-INTEGRATION.md)** for:
- Detailed input/output specs for each tool
- Full examples and integration patterns
- Advanced setup (config files, custom endpoints)
- Troubleshooting

---

### For Claude Code Development

- Full codebase access; no restrictions
- Godoc-driven API discoverability
- Use `/understand-anything:understand` to generate architecture graphs
- Memory persists across sessions at `.claude/projects/.../memory/`

### For Other Editors (Copilot, etc.)

- `.gitignore` automatically hides vendor/, test caches, build artifacts
- `.editorconfig` enforces consistent formatting
- Godoc comments serve as inline documentation
- Generated files (`*_gen.go`) are read-only; never suggest edits to them

## Version & Release

**Current:** 0.0.1

### Tooling Versions

Dev tools (goverter, revive) are pinned in two places for reproducibility:
- `go.mod`: Listed as `// +build tools` dependencies (tracked by `tools.go`)
- `Taskfile.yml`: `REVIVE_VERSION` and `GOVERTER_VERSION` variables

To update a tool:
1. Edit `tools.go` to import the new version
2. Run `go mod tidy` to update `go.mod`
3. Extract the version from `go.mod` and update `Taskfile.yml` vars
4. Run `task ci` to verify everything works

### To release

1. **All code changes committed and tests passing:**
   ```bash
   task ci  # Local verification
   ```

2. **Tag the release:**
   ```bash
   git tag v0.0.1
   git push origin v0.0.1
   ```

3. **GitHub Actions does the rest:**
   - `release.yml` workflow automatically:
     - Builds CLI and MCP server for Linux (amd64, arm64), macOS (amd64, arm64), Windows (amd64)
     - Creates a GitHub Release with all binaries
     - Triggers pkg.go.dev indexing (automatic; no action needed)

4. **Users can then install via:**
   ```bash
   go get github.com/arashackdev/doopl@v0.0.1
   go install github.com/arashackdev/doopl/cmd/doopl@v0.0.1
   ```

### pkg.go.dev Publishing

No setup required. pkg.go.dev automatically:
- Polls GitHub for new tags matching `v*`
- Indexes the release within seconds
- Makes godoc available at `https://pkg.go.dev/github.com/arashackdev/doopl@v0.0.1`

### Versioning
- Follows semantic versioning (MAJOR.MINOR.PATCH)
- v0.0.1 is the first public release; v0.x.y may have breaking changes before v1.0.0
- v1.0.0 planned after M6 (glossaries, write API, docs complete)

## Common Tasks

### Add a new endpoint
1. Define `apimodel` types in `internal/apimodel/`
2. Define public `model` types in `model/`
3. Add converter interface in `internal/convert/apimodel_to_model.go`
4. Run `task generate` (regenerates `*_gen.go`)
5. Add Client method in `client.go`
6. Add CLI command in `cmd/doopl/main.go`
7. Add tests in `*_test.go`
8. Run `task ci` to verify

### Debug a failing test
```bash
go test ./... -v -run TestName
go test ./... -race  # check for race conditions
go test ./... -cover # coverage report
```

### Check for API drift
```bash
# After DeepL API updates:
# Compare internal/apimodel against the official OpenAPI spec
# Manual process; automated drift detection is M6
```

### Generate converters after editing interface
```bash
task generate
```

This runs `go generate ./...`, which invokes `goverter` on all interfaces tagged with `//go:generate`.

## Documentation Standards

### Documentation Location Policy (ENFORCED)

**All documentation goes in `docs/` directory — NEVER scatter .md or .sh files in the project root.**

Exceptions (project root only):
- `README.md` — Project overview
- `CLAUDE.md` — This file (Claude setup guide)
- `LICENSE.md` — License text

**Enforcement:**
- `.gitignore` patterns prevent scattered `*.md` and `*.sh` in root from being committed
- Pattern: `/*.md` and `/*.sh` (root-level only, allowed files are exceptions via `!` rules)
- CI will fail if scattered docs appear in git tree

**Examples:**
- ❌ TESTING_PLAYBOOK.md (root) → ✅ docs/TESTING_PLAYBOOK.md
- ❌ QUICK_TEST.sh (root) → ✅ docs/QUICK_TEST.sh
- ❌ REVIEW_SUMMARY.md (root) → ✅ docs/REVIEW_SUMMARY.md
- ✅ README.md (root) — exception, stay here
- ✅ CLAUDE.md (root) — exception, stay here

**When creating new docs or scripts:**
1. **Create in docs/**: `docs/MY_GUIDE.md` or `docs/my_script.sh`
2. **Update docs/DEVELOPMENT.md**: Add entry to the documentation index
3. **Update code comments**: If the doc explains a design decision, add a comment to the relevant code pointing to `docs/MY_GUIDE.md`
4. **Never run `git add *.md` or `git add *.sh` at root level** — `.gitignore` patterns will prevent commit
5. **Local pre-commit hook** (optional, for strict enforcement):
   ```bash
   # Create in .git/hooks/pre-commit:
   #!/bin/bash
   if git diff --cached --name-only | grep -qE '^[^/]+\.(md|sh)$'; then
     echo "❌ Scattered docs/scripts in root are prohibited. Move to docs/ directory."
     exit 1
   fi
   ```

### Godoc comments
- One-line summary at the start, no punctuation (reads as "X is ...")
- For functions: describe inputs, outputs, and notable errors
- For types: describe the semantic meaning, not the fields
- Example:

```go
// TextResult represents a translated text segment.
type TextResult struct {
    Text string
    // ...
}

// TranslateText translates text into the target language.
// It returns translation results for each input text, in order.
// If ctx is canceled, TranslateText returns the context error.
func (c *Client) TranslateText(ctx context.Context, texts []string, targetLang string, opts ...TranslateTextOption) ([]model.TextResult, error)
```

### Code Comments (Link to Docs, Don't Duplicate)
- If the comment explains *why* a decision was made, add a code comment
- If the comment is a *how-to guide or reference*, put it in `docs/` and link from code
- Comments are for *why* (constraints, non-obvious decisions), not *what* (code already shows that)

**Example:**
```go
// For glossary naming conventions, see docs/DESIGN_DECISIONS.md#glossary-ids
// We use base32 to avoid URL encoding issues in API responses.
const GlossaryIDFormat = "base32"
```

### No hand-written comments on logic
- If the code is unclear, improve the code (rename variables, break into smaller functions)
- Comments are for *why* (constraints, non-obvious decisions), not *what* (code already shows that)

## Roadmap (M1–M6)

- **M0/M1 (done):** Text translation, full parameter parity, three-layer architecture
- **M2:** Languages & usage endpoints + CLI commands
- **M3:** Document translation (upload, poll, download)
- **M4:** Glossaries (create, list, entries, delete)
- **M5:** Write API (rephrase)
- **M6:** Polish, examples, deepl-mock in CI, GoReleaser, v1.0.0

See `docs/scope-and-checklist.md` for detailed status on each milestone.

## Troubleshooting

### `task generate` fails
- Check that interface declarations in `internal/convert/apimodel_to_model.go` and `cmd/doopl/internal/convert/model_to_entity.go` are syntactically valid
- Run `go generate -v ./...` for detailed output

### `revive` linting fails
- Missing godoc? Run `go doc ./...` to see which symbols lack comments
- Naming violation? Check `var-naming` rule in `revive.toml` — wire-format types use snake_case, everything else uses camelCase

### Tests fail with "connection refused"
- Ensure tests create their own `httptest.Server` (they should; see `translate_test.go`)
- No external API key should be required for tests

### CLI doesn't recognize a flag
- Check `cmd/doopl/main.go` — each flag must be registered in the command definition
- Flag name must match the option name (e.g., `--formality` → `WithFormality`)

## License

MIT. See [LICENSE](./LICENSE) for details.

## Support

- **Issues:** https://github.com/arashackdev/doopl/issues
- **Discussions:** https://github.com/arashackdev/doopl/discussions
- **Documentation:** [pkg.go.dev](https://pkg.go.dev/github.com/arashackdev/doopl)
