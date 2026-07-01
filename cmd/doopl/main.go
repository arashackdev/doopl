// Command doopl is a command-line interface to the DeepL translation API v3.
// It is built on top of the github.com/arashackdev/doopl library and serves
// as a practical example of library usage. The CLI is designed to be both useful
// as a standalone tool and educational as a reference implementation.
//
// Each doopl subcommand maps directly to a Client method, with CLI flags mapping
// to functional options. No translation logic lives in the CLI — it is purely
// flag parsing, I/O, and output formatting.
//
// # Usage
//
//	doopl translate "Hello, world!" --target-lang DE
//	doopl languages --resource translate
//	doopl usage
//	doopl glossary create --name tech --source EN --target DE --entries "API:Schnittstelle"
//
// # Global Flags
//
// - --auth-key: DeepL API key (env: DEEPL_AUTH_KEY)
// - --output: Output format: "text" (plain), "tui" (rich), "json" (structured)
// - --server-url: Override API endpoint (for testing or custom deployments)
//
// # Available Commands
//
// - translate: Translate text into a target language
// - document: Upload, check, and download document translations
// - glossary: Create, list, and manage translation glossaries
// - languages: List supported languages for translation/document/glossary/write
// - usage: Check your API quota and usage
// - rephrase: Rephrase text using the Write API
//
// # File Organization
//
// - main.go: App initialization and command registration
// - auth.go: Authentication key resolution and validation
// - output.go: Output formatting (text, TUI, JSON)
// - commands.go: Shared command utilities and error handling
// - command_*.go: Individual command implementations
//
// For more details, see the library documentation at pkg.go.dev/github.com/arashackdev/doopl
package main

import (
	"fmt"
	"os"

	"github.com/arashackdev/doopl/cmd/doopl/internal/convert"
	_ "github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
)

func main() {
	modelToEntity := &convert.ModelToEntityImpl{}

	app := &cli.App{
		Name:  "doopl",
		Usage: "translate text from the command line using the DeepL API",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "auth-key",
				Usage:   "DeepL API auth key (precedence: flag > env DEEPL_AUTH_KEY > config)",
				EnvVars: []string{"DEEPL_AUTH_KEY"},
			},
			&cli.StringFlag{
				Name:  "output",
				Usage: "output format: text (plain)|tui (rich terminal)|json (structured)",
				Value: "text",
			},
		},
		Commands: []*cli.Command{
			configCommand(),
			translateCommand(modelToEntity),
			documentCommand(),
			glossaryCommand(),
			rephraseCommand(),
			languagesCommand(),
			usageCommand(),
			doctorCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "doopl:", err)
		os.Exit(1)
	}
}
