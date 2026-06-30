# doopl

An idiomatic Go client for the [DeepL API](https://developers.deepl.com), usable as a **library**, **CLI**, or **MCP server for AI clients**.

> **Version:** 0.0.1 — initial release. Supports text translation, document translation, glossaries, languages, usage, and text rephrasing (Write API). All public APIs are fully documented with godoc available on [pkg.go.dev](https://pkg.go.dev/github.com/arashackdev/doopl).

## Use with Claude & AI Clients (MCP)

**Fastest way to get started:** Install doopl as an MCP server for Claude Code, Claude Desktop, or any MCP-compatible AI client.

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

| Tool | What It Does |
|------|-------------|
| `translate` | Translate text to a target language, optionally with source language, formality, and glossary |
| `languages` | List supported languages for translate/document/glossary/write resources |
| `usage` | Check character and document quota usage |

**Details:** See the full [AI Integration Guide](https://github.com/arashackdev/doopl/blob/main/.claude/AI-INTEGRATION.md) for input/output specs, examples, and advanced setup.

## CLI usage

```sh
go install github.com/arashackdev/doopl/cmd/doopl@latest

export DEEPL_AUTH_KEY=your-key-here

# Translate text
doopl translate --to DE "Hello, world!"

# List supported languages
doopl languages

# Check API usage
doopl usage
```

The CLI is a thin `urfave/cli` wrapper — every flag maps directly to a library option. No translation logic lives in `cmd/doopl`; it's proof the library is genuinely embeddable, not a parallel implementation.

## Features

### Full DeepL v3 API Coverage

- ✅ **Text Translation** — all 8 options (formality, glossary, context, style, tag handling, custom instructions, translation memory, model type)
- ✅ **Document Translation** — upload, poll, download with automatic backoff
- ✅ **Glossaries** — create, list, get, delete, fetch entries
- ✅ **Languages** — list for translate/document/glossary/write resources
- ✅ **Usage** — character and document quotas
- ✅ **Write API** — rephrase with tone and emoji control

### Production Ready

- Free/Pro endpoint auto-detection from API key
- Exponential backoff + jitter on rate limits and errors
- Retry-After header support
- Context cancellation throughout
- Sentinel errors with errors.Is/errors.As
- Concurrent-safe Client
- Full godoc on all APIs

## Architecture

See [`docs/scope-and-checklist.md`](./docs/scope-and-checklist.md) for rationale. The project enforces a three-layer separation:

- **apimodel** (internal): wire-format types matching the DeepL API exactly
- **model**: public, Go-idiomatic domain types
- **entity** (CLI): display types for table/JSON output

Each layer is isolated via generated converters (`goverter`), so changes to the API surface don't cascade to consumer code. Converters are never hand-edited — they're automatically regenerated from interface declarations.

## Development

```sh
task fmt           # format code
task vet           # run go vet
task generate      # regenerate converters
task lint          # lint with revive
task test          # run tests
task ci            # full local CI suite
task cli:build     # build the CLI
task mcp:build     # build MCP server for AI clients
```

Tests run against an in-process `httptest` server with 48% coverage. All exported symbols are documented with godoc.

See [`CLAUDE.md`](./CLAUDE.md) for development guide and [`MILESTONES.md`](./MILESTONES.md) for feature completion status.

## Minimum Go version

Go 1.23+

## Examples

See [`_examples/`](./_examples/) for runnable examples:

- `translate.go` — basic text translation
- `document.go` — document upload/poll/download round-trip
- `glossary.go` — glossary CRUD and usage

Build and run with `go run -tags=ignore _examples/translate.go`.

## License

Apache License 2.0. See [LICENSE](./LICENSE) file.
