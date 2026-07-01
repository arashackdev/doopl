# Claude Development Notes (Local)

This file is development-only guidance for working with the doopl codebase using Claude. Not shipped with the package.

## Key Design Principles

- **Three-layer architecture:** API wire format → public domain model → CLI display entity, each isolated by generated converters
- **Full parameter parity:** `TranslateText` supports all DeepL v3 options
- **Idiomatic Go:** Functional options, sentinel errors with `errors.Is`/`errors.As`, free/pro endpoint auto-detection
- **Zero hand-written converters:** mappings generated via `goverter` from interface declarations

## Development Stack

| Tool       | Purpose         | Config                                        |
| ---------- | --------------- | --------------------------------------------- |
| Go 1.24+   | Language        | `go.mod`                                      |
| Task       | Task runner     | `Taskfile.yml`                                |
| revive     | Linter          | `revive.toml`                                 |
| goverter   | Code generation | interface declarations in `internal/convert/` |
| urfave/cli | CLI framework   | `cmd/doopl/main.go`                           |

## Quick Commands

```bash
go mod download         # Install dependencies
task ci                 # Full CI suite (fmt, vet, lint, generate, test)
task cli:build          # Build optimized CLI to ./bin/doopl
task test               # Run tests
task generate           # Regenerate converters
task fmt                # Format code
```

## Project Structure

See `docs/ARCHITECTURE.md` for full architecture guide.

```
doopl/
├── client.go                    # Client + options
├── errors.go                    # Sentinel errors
├── translate.go / translate_test.go
├── document.go / document_test.go
├── glossary.go / glossary_test.go
├── languages.go / languages_test.go
├── write.go / write_test.go
├── request.go                   # HTTP helpers
│
├── model/                       # Public types
├── internal/
│   ├── apimodel/               # Wire format
│   ├── convert/                # Generated converters
│   ├── transport/              # HTTP retry/backoff
│   └── config/                 # Config manager
│
├── cmd/doopl/                  # CLI
│   ├── main.go
│   └── internal/entity + convert/
│
├── examples/                   # Runnable examples
├── docs/
│   ├── SETUP.md               # Setup & usage
│   ├── ARCHITECTURE.md        # Design details
│   └── scope-and-checklist.md # Roadmap
│
├── Taskfile.yml
├── revive.toml
└── go.mod
```

## Adding a New Endpoint

Mechanical pattern for each new resource (M2–M6):

1. **Create API wire format:**
   ```go
   // internal/apimodel/languages.go
   type LanguagesResponse struct {
       Languages []Language `json:"languages"`
   }
   ```

2. **Create public domain type:**
   ```go
   // model/languages.go
   type Language struct {
       Code string
       Name string
   }
   ```

3. **Add converter interface:**
   ```go
   // internal/convert/apimodel_to_model.go
   type APIToModel interface {
       LanguagesToModel(src []apimodel.Language) []model.Language
   }
   ```

4. **Run converter generation:**
   ```bash
   task generate
   ```

5. **Add Client method:**
   ```go
   func (c *Client) Languages(ctx context.Context, resource string) ([]model.Language, error)
   ```

6. **Add CLI command** in `cmd/doopl/main.go`

7. **Add tests** in `*_test.go`

8. **Run CI:**
   ```bash
   task ci
   ```

## Code Quality Checks

All must pass:

```bash
task fmt:check          # gofmt check
task vet                # go vet
task lint               # revive linting
task generate:check     # Converters up-to-date
task test               # Unit tests with -race
```

**Linter rules (revive.toml):**
- `exported`: All public symbols must have godoc
- `package-comments`: Package must have a comment
- `var-naming`: snake_case in apimodel (wire format), camelCase elsewhere
- `blank-imports`: No empty imports
- `unreachable-code`: Flagged
- `time-naming`: Duration vars should end in `Duration`

## Testing

Tests use in-process `httptest` server (no real API key needed):

```bash
go test ./... -race -cover
```

Coverage: ~45%. Acceptable for client libraries.

## Common Issues

### `task generate` fails
- Check interface declarations in `internal/convert/apimodel_to_model.go` syntax
- Run `go generate -v ./...` for detailed output

### `revive` linting fails
- Missing godoc? Run `go doc ./...`
- Naming violation? Check `revive.toml` rules

### Tests fail with "connection refused"
- Tests create their own `httptest.Server` (see `translate_test.go`)
- No external API key required

## Environment & Secrets

### Required
- `DEEPL_AUTH_KEY`: Your DeepL API key (free keys end in `:fx`)

### Optional
- `DEEPL_SERVER_URL`: Override API endpoint (for testing with deepl-mock)

## Versioning

**Current:** 0.0.1

### Semantic Versioning
- v0.0.1 is the first public release (M0/M1 complete)
- v0.x.y may have breaking changes before v1.0.0
- v1.0.0 planned after M6 (all features, examples, CI, docs complete)

### To Release
1. All code changes committed
2. `task ci` passes locally
3. Tag: `git tag v0.0.1`
4. Push: `git push origin v0.0.1`
5. GitHub Actions auto-builds and publishes to pkg.go.dev

## Roadmap (M0–M6)

- **M0/M1 (DONE):** Text translation, full parameter parity, three-layer architecture
- **M2:** Languages & usage endpoints + CLI commands
- **M3:** Document translation (upload, poll, download)
- **M4:** Glossaries (create, list, entries, delete)
- **M5:** Write API (rephrase)
- **M6:** Polish, examples, deepl-mock in CI, GoReleaser, v1.0.0

See `docs/scope-and-checklist.md` for detailed status.

## Support

- **Issues:** https://github.com/arashackdev/doopl/issues
- **Discussions:** https://github.com/arashackdev/doopl/discussions
- **Documentation:** [pkg.go.dev](https://pkg.go.dev/github.com/arashackdev/doopl)
