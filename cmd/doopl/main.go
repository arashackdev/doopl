// Command doopl is a thin CLI over the github.com/arashackdev/doopl library.
// It exists to prove the library is genuinely embeddable: every flag here maps
// directly onto a doopl option, and the CLI holds no translation logic of its own
// — only flag parsing and I/O.
//
// The CLI is organized into separate files by concern:
// - auth.go: authentication key resolution
// - output.go: output formatting helpers
// - commands.go: shared command utilities
// - command_*.go: individual command implementations
package main

import (
	"fmt"
	"os"

	"github.com/arashackdev/doopl/cmd/doopl/internal/convert"
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
				Usage: "global output format: json|table|text",
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
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "doopl:", err)
		os.Exit(1)
	}
}
