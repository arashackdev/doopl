# DoopL — DeepL API v3 Client for Go

![Go Version](https://img.shields.io/badge/Go-1.24%2B-00ADD8?style=flat-square)
![Test Coverage](https://img.shields.io/badge/coverage-48.1%25-yellow?style=flat-square)
![Code Quality](https://img.shields.io/badge/linter-revive-green?style=flat-square)
![License](https://img.shields.io/badge/license-MIT-blue?style=flat-square)

An idiomatic Go client for the [DeepL API v3](https://developers.deepl.com), usable as a **library**, **CLI**, or **MCP server for AI clients**.

## Use with Claude & AI Clients (MCP)

**Fastest way to get started:** Install DoopL as an MCP server for Claude Code, Claude Desktop, or any MCP-compatible AI client.

```bash
# Build the MCP server
go build -o doopl-mcp ./cmd/doopl-mcp

# Add to Claude Code settings (.claude/settings.json)
{
  "mcpServers": {
    "doopl": {
      "command": "/path/to/doopl-mcp",
      "args": ["serve"],
      "env": { "DEEPL_AUTH_KEY": "your-key-here" }
    }
  }
}
```

Now Claude can translate, list languages, and check API usage directly. See **[AI Integration Guide](#ai-integration-guide)** for setup, examples, and use cases.

## Install

```sh
go get github.com/arashackdev/doopl
```

## Library usage

### Text Translation

```go
client, err := doopl.New(os.Getenv("DEEPL_AUTH_KEY"))
if err != nil {
    log.Fatal(err)
}

results, err := client.TranslateText(ctx, []string{"Hello, world!"}, "DE",
    doopl.WithFormality(doopl.FormalityMore))
if err != nil {
    var apiErr *doopl.Error
    if errors.As(err, &apiErr) && errors.Is(err, doopl.ErrQuotaExceeded) {
        // handle quota specifically
    }
    log.Fatal(err)
}
fmt.Println(results[0].Text) // "Hallo, Welt!"
```

### Document Translation

```go
file, err := os.Open("document.pdf")
defer file.Close()
output, err := os.Create("document_de.pdf")
defer output.Close()

err = client.TranslateDocument(ctx, file, "document.pdf", "DE", output)
```

### Glossaries

```go
entries := model.GlossaryEntries{
    "API":    "Schnittstelle",
    "server": "Server",
}
glos, err := client.CreateGlossary(ctx, "tech-terms", "EN", "DE", entries)

results, err := client.TranslateText(ctx, []string{"The API is fast"}, "DE",
    doopl.WithGlossaryID(glos.GlossaryID))
```

### Rephrase (Write API)

```go
results, err := client.Rephrase(ctx, []string{"Hello!"}, "EN",
    doopl.WithTone(model.WriteToneFormal))
fmt.Println(results[0].Text) // "Greetings!"
```

### Languages & Usage

```go
langs, err := client.SourceLanguages(ctx)
usage, err := client.Usage(ctx)
```

Free vs. Pro API endpoint is detected automatically from your key (free keys end in `:fx`). Override with `doopl.WithServerURL(...)` — also how you point the client at [`deepl-mock`](https://github.com/DeepLcom/deepl-mock) for testing.

## AI Integration Guide

### What You Can Do

With doopl's MCP server, Claude (or any AI client) can:

- **Translate text** — supports all DeepL options (formality, glossaries, context, style, tag handling)
- **List languages** — query supported languages for any resource type (translate, document, glossary, write)
- **Check quota** — inspect character and document usage before large operations
- **Multi-language workflows** — translate documents, test UI localization, route support tickets

### Quick Example: Claude + doopl

```
Human: "Translate this support ticket from German to English and draft a friendly reply"

Claude:
1. Calls translate(text=<ticket>, target_lang=EN)
2. Gets English text: "The login page is broken"
3. Drafts reply: "Thank you for reporting. We're investigating..."
4. Calls translate(text=<reply>, target_lang=DE)  
5. Returns German reply to send to customer
```

### Setup

1. **Build the server:**
   ```bash
   task mcp:build  # Creates ./bin/doopl-mcp
   ```

2. **Configure Claude Code** (`.claude/settings.json`):
   ```json
   {
     "mcpServers": {
       "doopl": {
         "command": "/full/path/to/bin/doopl-mcp",
         "args": ["serve"],
         "env": { "DEEPL_AUTH_KEY": "your-key" }
       }
     }
   }
   ```

3. **Restart Claude** — it will detect and load the new tools.

### Available Tools

| Tool        | What It Does                                                                                  |
| ----------- | --------------------------------------------------------------------------------------------- |
| `translate` | Translate text to a target language, optionally with source language, formality, and glossary |
| `languages` | List supported languages for translate/document/glossary/write resources                      |
| `usage`     | Check character and document quota usage                                                      |

**Details:** See the full [AI Integration Guide](https://github.com/arashackdev/doopl/blob/main/.claude/AI-INTEGRATION.md) for input/output specs, examples, and advanced setup.

## CLI Usage

```sh
go install github.com/arashackdev/doopl/cmd/doopl@latest

export DEEPL_AUTH_KEY=your-key-here

# Translate text
doopl translate --to DE "Hello, world!"

# List supported languages
doopl languages

# Check API usage
doopl usage

# Check API health (light and verbose diagnostics)
doopl doctor
doopl doctor --verbose

# Different output formats (text, tui, json)
doopl translate --to DE --output tui "Hello"
doopl languages --output json
```

The CLI is a thin `urfave/cli` wrapper — every flag maps directly to a library option. No translation logic lives in `cmd/doopl`; it's proof the library is genuinely embeddable, not a parallel implementation.

**Output Modes:**
- **text** (default) — Plain text for scripts and terminals
- **tui** — Rich terminal UI with colors and formatting (lipgloss)
- **json** — Structured output for programmatic use

## Features

### Full DeepL v3 API Coverage

- ✅ **Text Translation** — all 8 options (formality, glossary, context, style, tag handling, custom instructions, translation memory, model type)
- ✅ **Document Translation** — upload, poll, download with automatic backoff
- ✅ **Glossaries** — create, list, get, delete, fetch entries
- ✅ **Languages** — list for translate/document/glossary/write resources
- ✅ **Usage** — character and document quotas
- ✅ **Write API** — rephrase with tone and emoji control
- ✅ **Health Check** — `doctor` command to verify connectivity and quota

### Production Ready

- **Robust Networking** — exponential backoff + jitter on rate limits; Retry-After header support
- **Context Support** — cancellation throughout; proper deadline handling
- **Error Handling** — sentinel errors with `errors.Is()`/`errors.As()`
- **Concurrency Safe** — Client safe for concurrent use (no locks needed)
- **Auto-Detection** — free/pro endpoint detected from API key
- **Fully Documented** — godoc on all public APIs; zero public API without docs

## Architecture

The project enforces a **three-layer separation** to keep concerns isolated:

- **Wire Format** (`v3/apimodel/`) — exact DeepL API shape, never exposed publicly
- **Domain Model** (`pkg/model/`) — idiomatic Go types, public API for library consumers
- **CLI Entity** (`cmd/doopl/internal/entity/`) — display types for terminal/JSON output

Each layer is isolated via **generated converters** (`goverter`). Converters are never hand-edited — they're automatically regenerated from interface declarations. This ensures:
- API changes don't cascade to consumer code
- No silent field mismatches
- CLI concerns stay separated from library concerns

See [`docs/DEVELOPMENT.md`](./docs/DEVELOPMENT.md) for detailed architecture guide.

## Development

```sh
task fmt           # format code
task vet           # run go vet
task generate      # regenerate converters
task lint          # lint with revive (enforces 100% godoc coverage)
task test          # run tests (48.1% coverage, race detector enabled)
task ci            # full local CI suite
task cli:build     # build the CLI
task mcp:build     # build MCP server for AI clients
```

### Code Quality

- **Test Coverage:** 48.1% of statements in library code
- **Linting:** revive with strict rules (godoc required for all exports)
- **Type Safety:** Full godoc on all public APIs; zero hand-written converters
- **Concurrency:** Race detector on all tests; Client is concurrent-safe
- **Minimum Go:** 1.24+

### Documentation

- **Development Guide:** See [`docs/DEVELOPMENT.md`](./docs/DEVELOPMENT.md) for architecture, testing, and contributing
- **API Docs:** All public symbols documented with godoc at [pkg.go.dev](https://pkg.go.dev/github.com/arashackdev/doopl)
- **MCP Setup:** See [`CLAUDE.md`](./CLAUDE.md) for AI client integration

## Minimum Go version

Go 1.24+

## Examples

See [`_examples/`](./_examples/) for runnable examples:

- `translate.go` — basic text translation
- `document.go` — document upload/poll/download round-trip
- `glossary.go` — glossary CRUD and usage

Build and run with `go run -tags=ignore _examples/translate.go`.
