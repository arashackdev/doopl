# doopl Examples

Complete examples showing how to use doopl as a library, CLI, and MCP server.

## Directories

### library-usage/
Go code examples using doopl as a library in your own applications.

- **simple-translate/** — Basic single-text translation
- **batch-translate/** — Multiple texts in one efficient API call
- **with-options/** — Advanced options (formality, source language, glossaries)

See [library-usage/README.md](./library-usage/README.md) for details.

### cli-usage/
Shell scripts and command-line workflows using the `doopl` CLI.

- **basic-workflow.sh** — Get started: languages, quota, translate, output formats
- **batch-translation.sh** — Translate multiple texts from a file

See [cli-usage/README.md](./cli-usage/README.md) for details.

### mcp-integration/
Setup and usage guide for the doopl MCP server with Claude and other AI clients.

- **README.md** — Full setup, tool specifications, and example use cases

See [mcp-integration/README.md](./mcp-integration/README.md) for details.

## Quick Start

### As a Library
```bash
cd library-usage/simple-translate
export DEEPL_AUTH_KEY="your-api-key"
go run main.go
```

### As a CLI
```bash
export DEEPL_AUTH_KEY="your-api-key"
doopl translate "Hello, world!" DE
```

### With Claude (MCP)
```bash
task mcp:build
# Configure in .claude/settings.json
# Then ask Claude: "Translate 'Hello, world!' to German"
```

## API Documentation

- **Library:** [pkg.go.dev/github.com/arashackdev/doopl](https://pkg.go.dev/github.com/arashackdev/doopl)
- **CLI Help:** `doopl --help` or `doopl COMMAND --help`
- **MCP Tools:** See [mcp-integration/README.md](./mcp-integration/README.md)

## Environment

All examples require the `DEEPL_AUTH_KEY` environment variable or a saved config:

```bash
export DEEPL_AUTH_KEY="your-api-key"
# or
doopl config set your-api-key
```

Get a free API key at [deepl.com](https://www.deepl.com/pro-api).

## License

Same as doopl: MIT + Commons Clause. See LICENSE in the repository root.
