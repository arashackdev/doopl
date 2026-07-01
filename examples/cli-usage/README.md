# CLI Usage Examples

These examples show how to use doopl from the command line.

## Setup

Build the CLI:
```bash
task cli:build
# or
go install ./cmd/doopl
```

Set your API key:
```bash
export DEEPL_AUTH_KEY="your-api-key-here"
```

Or configure it once:
```bash
doopl config set your-api-key-here
```

## Quick Start

Translate a single text:
```bash
doopl translate "Hello, world!" DE
```

List supported languages:
```bash
doopl languages translate
```

Check your API quota:
```bash
doopl usage
```

## Output Formats

All commands support three output formats via `--output`:

### text (default)
Plain text output, easy to read:
```bash
doopl translate "Hello" DE --output text
```

### json
Structured JSON, useful for scripting:
```bash
doopl translate "Hello" DE --output json
```

### tui
Rich terminal UI with colors and formatting:
```bash
doopl translate "Hello" DE --output tui
```

## Examples

### basic-workflow.sh
Complete workflow showing:
- List languages
- Check quota
- Translate with options
- Use different output formats

```bash
chmod +x basic-workflow.sh
./basic-workflow.sh
```

### batch-translation.sh
Translate multiple texts from a file:

```bash
chmod +x batch-translation.sh
./batch-translation.sh
```

## Advanced: Document Translation

Translate a document (supports .pdf, .docx, .xlsx, .pptx, .html, .txt):
```bash
doopl document translate --input document.pdf --target DE --output translated.pdf
```

Monitor progress:
```bash
# Get document ID and key
DOC_ID=$(doopl document upload --input file.pdf --target DE | jq -r .document_id)
DOC_KEY=$(doopl document upload --input file.pdf --target DE | jq -r .document_key)

# Check status
doopl document status --document-id $DOC_ID --document-key $DOC_KEY

# Download when done
doopl document download --document-id $DOC_ID --document-key $DOC_KEY --output translated.pdf
```

## Tips

- Use `--output json` in scripts and pipes
- Use `--output tui` for interactive use
- Use `--output text` for simple copy-paste
- Add `--formality more` or `--formality less` for German/Dutch
- Use `--glossary-id ID` to apply a glossary to translations

Run `doopl --help` or `doopl COMMAND --help` for full documentation.
