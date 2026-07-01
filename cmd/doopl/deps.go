package main

// Dependencies used in cmd/doopl:
//
// - github.com/arashackdev/doopl/pkg/deepl: Main library client
//   Used by: All command handlers to call API methods
//   Key types: Client, TranslateTextOption, TranslateDocumentOption, etc.
//   See: ../../../pkg/deepl/client.go
//
// - github.com/arashackdev/doopl/cmd/doopl/internal/convert: Goverter-generated
//   converter from model → CLI entity types (ModelToEntity interface)
//   Used by: Command handlers for formatting output
//   See: ./internal/convert/model_to_entity.go
//
// - github.com/arashackdev/doopl/cmd/doopl/internal/entity: CLI display types
//   (TranslationRow, LanguageRow, UsageRow)
//   Used by: Output formatting (JSON, table, text)
//   See: ./internal/entity/
//
// - github.com/arashackdev/doopl/internal/config: Config file management
//   Used by: getAuthKey() to load saved auth keys from ~/.doopl/config.toml
//   See: ../../../internal/config/config.go
//
// - github.com/urfave/cli/v2: CLI framework
//   Used by: main() for command parsing, flags, subcommands
//
// - github.com/arashackdev/doopl/pkg/model: Public domain types
//   Used by: Command return types and conversion
//   See: ../../../pkg/model/
//
// See: main.go (entry point and command definitions)
