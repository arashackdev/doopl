# Library Usage Examples

These examples show how to use doopl as a Go library in your own applications.

## Setup

Set your DeepL API key:
```bash
export DEEPL_AUTH_KEY="your-api-key-here"
```

## Examples

### simple-translate
Translate a single string to a target language.

```bash
go run simple-translate/main.go
```

### batch-translate
Translate multiple strings in one efficient API call. Results maintain input order.

```bash
go run batch-translate/main.go
```

### with-options
Translate with advanced options like formality level, source language detection, and more.

```bash
go run with-options/main.go
```

## Common Patterns

### Check supported languages
```go
langs, err := client.Languages(ctx, "translate")
for _, lang := range langs {
    fmt.Printf("%s (%s)\n", lang.Code, lang.Name)
}
```

### Check usage and quota
```go
usage, err := client.Usage(ctx)
fmt.Printf("Characters used: %d/%d\n", usage.CharacterCount, usage.CharacterLimit)
```

### Error handling
```go
results, err := client.TranslateText(ctx, texts, "DE")
if err != nil {
    if errors.Is(err, deepl.ErrQuotaExceeded) {
        log.Fatal("API quota exceeded")
    }
    log.Fatalf("Translation failed: %v", err)
}
```

## API Documentation

Full API documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/arashackdev/doopl).
