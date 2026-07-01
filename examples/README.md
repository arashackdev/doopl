# doopl Examples

Complete, runnable examples showing how to use doopl as a Go library, CLI tool, and MCP server for Claude AI integration.

## 🚀 Quick Start

### Prerequisites
- Go 1.24 or later
- DeepL API key (free: https://www.deepl.com/pro-api)

### Set Your API Key
```bash
export DEEPL_AUTH_KEY="your-deepl-api-key"
```

## 📚 Library Examples

Using doopl as a package in your own Go code.

### Basic Translation
Translate a single text with default options:
```bash
go run examples/translate.go
```

### With Formality & Options
Translate with advanced options (coming soon):
```bash
# Reference implementations in examples/
# - glossary.go: Create glossaries, translate with custom terms
# - document.go: Upload, poll, download document translations
```

### Batch Translation
Translate multiple texts efficiently in one call:
```bash
cd examples/library-usage/batch-translate
export DEEPL_AUTH_KEY="your-api-key"
go run main.go
```

### Advanced Options
Translate with formality, glossaries, context, and more:
```bash
cd examples/library-usage/with-options
export DEEPL_AUTH_KEY="your-api-key"
go run main.go
```

## 🗂️ Example Directory Structure

### `examples/*.go` — Standalone Reference Examples
Simple, self-contained files you can run directly with `go run`:

- **`translate.go`** — Basic text translation
  ```bash
  go run examples/translate.go
  ```
  Demonstrates: Multiple texts, formality, detected language

- **`glossary.go`** — Glossary management
  ```bash
  go run examples/glossary.go
  ```
  Demonstrates: Create, list, use in translation, delete

- **`document.go`** — Document translation workflow
  ```bash
  go run examples/document.go
  ```
  Demonstrates: Upload, poll status, download (reference implementation)

### `examples/library-usage/` — Structured Library Examples
Professional project examples for embedding doopl in your application.

#### `simple-translate/`
Minimal example: Create a client, translate text.
```bash
cd examples/library-usage/simple-translate
go run main.go
# Output: Hallo, Welt!
```

#### `batch-translate/`
Efficient batch translation of multiple texts:
```bash
cd examples/library-usage/batch-translate
go run main.go
# Output: 3 translations in a single API call
```

#### `with-options/`
Advanced translation options (formality, glossaries, context):
```bash
cd examples/library-usage/with-options
go run main.go
# Output: Translations with custom options applied
```

### `examples/cli-usage/` — Command-Line Workflows
Shell scripts demonstrating the doopl CLI tool.

#### `basic-workflow.sh`
Get started with the CLI: check languages, usage, translate, format output.
```bash
chmod +x examples/cli-usage/basic-workflow.sh
export DEEPL_AUTH_KEY="your-api-key"
bash examples/cli-usage/basic-workflow.sh
```

#### `batch-translation.sh`
Translate texts from a file using the CLI:
```bash
chmod +x examples/cli-usage/batch-translation.sh
export DEEPL_AUTH_KEY="your-api-key"
bash examples/cli-usage/batch-translation.sh
```

### `examples/mcp-integration/` — Claude AI Integration
Setup and usage guide for the doopl MCP server with Claude and Claude Desktop.

See [mcp-integration/README.md](./mcp-integration/README.md) for:
- Setup instructions for Claude Code and Claude Desktop
- Tool specifications (translate, languages, usage)
- Example use cases and prompts

## 🛠️ Running Examples

### Environment Setup

All examples read `DEEPL_AUTH_KEY` from your environment:
```bash
export DEEPL_AUTH_KEY="your-api-key"
```

Or, with the doopl CLI, save it to config:
```bash
doopl config set your-api-key
```

### From the Root Directory

Run any example from the repository root:
```bash
go run examples/translate.go
go run examples/glossary.go
go run examples/document.go
```

Or run examples in subdirectories:
```bash
go run examples/library-usage/simple-translate/main.go
go run examples/library-usage/batch-translate/main.go
go run examples/library-usage/with-options/main.go
```

### Individual Example Directories

Enter an example directory and run locally:
```bash
cd examples/library-usage/simple-translate
go run main.go
```

### With Different Target Languages

Most examples are configurable. Edit the `.go` file or pass environment variables:
```bash
# Example: Translate to Spanish instead of German
DEEPL_TARGET_LANG=ES go run examples/translate.go
```

## 📖 Example Patterns

### Pattern 1: Create Client & Translate
```go
import "github.com/arashackdev/doopl/pkg/deepl"

client, err := deepl.New(os.Getenv("DEEPL_AUTH_KEY"))
results, err := client.TranslateText(ctx, []string{"Hello"}, "DE")
```

### Pattern 2: With Functional Options
```go
results, err := client.TranslateText(ctx, texts, "FR",
    deepl.WithFormality(deepl.FormalityMore),
    deepl.WithSourceLang("EN"),
)
```

### Pattern 3: Error Handling
```go
import "errors"

results, err := client.TranslateText(ctx, texts, lang)
if err != nil {
    if errors.Is(err, deepl.ErrQuotaExceeded) {
        // Handle quota exhaustion
    }
}
```

### Pattern 4: Concurrent Translation
```go
// Client is safe for concurrent use
var wg sync.WaitGroup
for _, text := range texts {
    wg.Add(1)
    go func(t string) {
        defer wg.Done()
        client.TranslateText(ctx, []string{t}, "DE")
    }(text)
}
wg.Wait()
```

## 🎯 Common Tasks

### Check Supported Languages
```bash
cd examples/library-usage/simple-translate
go run main.go
# or with CLI:
doopl languages --resource translate
```

### Translate & Save to File
```go
results, _ := client.TranslateText(ctx, []string{text}, lang)
os.WriteFile("output.txt", []byte(results[0].Text), 0644)
```

### Monitor API Usage
```bash
doopl usage
# or in code:
usage, _ := client.Usage(ctx)
fmt.Printf("Used %d / %d chars\n", usage.CharacterCount, usage.CharacterLimit)
```

### Translate Documents
See `examples/document.go` for the full workflow:
1. Upload with `DocumentUpload()`
2. Poll with `DocumentStatus()`
3. Download with `DocumentDownload()`

## 🔗 API Documentation

- **Library API:** [pkg.go.dev/github.com/arashackdev/doopl/pkg/deepl](https://pkg.go.dev/github.com/arashackdev/doopl/pkg/deepl)
- **Model Types:** [pkg.go.dev/github.com/arashackdev/doopl/pkg/model](https://pkg.go.dev/github.com/arashackdev/doopl/pkg/model)
- **CLI Reference:** `doopl --help` or `doopl COMMAND --help`
- **MCP Setup:** [examples/mcp-integration/README.md](./mcp-integration/README.md)

## 🐛 Troubleshooting

### "DEEPL_AUTH_KEY not set"
```bash
export DEEPL_AUTH_KEY="your-api-key"
go run examples/translate.go
```

### API Key Invalid
- Check your key at https://www.deepl.com/account/keys
- Free API keys end in `:fx`
- Pro API keys don't have a suffix

### "Too many requests" or Rate Limiting
The client retries automatically with exponential backoff (default: 5 retries).
Customize with `deepl.WithMaxRetries(n)` or space out requests manually.

### Document Upload Fails
Supported formats: PDF, DOCX, PPTX, XLSX, TXT, HTML, HTM, JPG, JPEG, PNG.
Ensure the file exists and is not corrupted.

## 📝 Creating Your Own Example

1. Copy a template (e.g., `library-usage/simple-translate/`)
2. Modify `main.go` for your use case
3. Run with `go run main.go`

Template:
```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/arashackdev/doopl/pkg/deepl"
)

func main() {
	authKey := os.Getenv("DEEPL_AUTH_KEY")
	if authKey == "" {
		log.Fatal("DEEPL_AUTH_KEY not set")
	}

	client, err := deepl.New(authKey)
	if err != nil {
		log.Fatal(err)
	}

	// Your code here
}
```

## 📦 Dependencies

All examples use the same dependencies as doopl:
- `github.com/arashackdev/doopl` — The client library
- Standard Go packages only for examples (no external dependencies in examples)

## 📄 License

MIT. See [LICENSE](../LICENSE) in the repository root for details.

---

**Questions?** See [CLAUDE.md](../CLAUDE.md) for development setup, or open an issue at https://github.com/arashackdev/doopl/issues.
