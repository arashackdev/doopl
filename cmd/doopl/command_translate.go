// Package main provides the doopl CLI.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/arashackdev/doopl/cmd/doopl/internal/convert"
	"github.com/arashackdev/doopl/cmd/doopl/internal/output"
	deepl "github.com/arashackdev/doopl/pkg/deepl"
	"github.com/urfave/cli/v2"
)

// translateCommand returns the translate command for translating text into a target language.
// The entity converter is used to convert model results to CLI display entities.
func translateCommand(entity *convert.ModelToEntityImpl) *cli.Command {
	return &cli.Command{
		Name:      "translate",
		Usage:     "translate text into a target language",
		ArgsUsage: "TEXT...",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "to", Required: true, Usage: "target language code, e.g. DE"},
			&cli.StringFlag{Name: "from", Usage: "source language code (auto-detected if omitted)"},
			&cli.StringFlag{Name: "formality", Usage: "default|more|less|prefer_more|prefer_less"},
			&cli.StringFlag{Name: "glossary-id", Usage: "glossary ID to apply (requires --from)"},
			&cli.StringFlag{Name: "context", Usage: "domain context for translation"},
			&cli.StringFlag{Name: "split-sentences", Value: "1", Usage: "0|1|nonewlines"},
			&cli.BoolFlag{Name: "preserve-formatting", Usage: "keep source formatting"},
			&cli.StringFlag{Name: "model-type", Usage: "quality_optimized|latency_optimized|prefer_quality_optimized"},
			&cli.StringFlag{Name: "tag-handling", Usage: "xml|html"},
			&cli.StringFlag{Name: "tag-handling-version", Usage: "v1|v2"},
			&cli.StringSliceFlag{Name: "custom-instructions", Usage: "style hints (repeatable)"},
			&cli.StringFlag{Name: "style-id", Usage: "style preset ID"},
			&cli.StringFlag{Name: "translation-memory-id", Usage: "translation memory ID"},
			&cli.Float64Flag{Name: "translation-memory-threshold", Usage: "TM match threshold (0.0-1.0)"},
			&cli.StringFlag{Name: "input-file", Usage: "read text from file instead of args"},
			&cli.StringFlag{Name: "output-file", Usage: "write results to file instead of stdout"},
		},
		Action: func(c *cli.Context) error {
			var text string

			// Read input from file or args
			if inputFile := c.String("input-file"); inputFile != "" {
				data, err := os.ReadFile(inputFile)
				if err != nil {
					return fmt.Errorf("cannot read input file: %w", err)
				}
				text = strings.TrimSpace(string(data))
			} else {
				text = strings.Join(c.Args().Slice(), " ")
			}

			if text == "" {
				return errors.New("no text provided")
			}

			authKey, err := getAuthKey(c)
			if err != nil {
				return err
			}

			client, err := deepl.New(authKey, deepl.WithAppInfo("doopl-cli", deepl.Version))
			if err != nil {
				return err
			}

			var opts []deepl.TranslateTextOption
			if v := c.String("from"); v != "" {
				opts = append(opts, deepl.WithSourceLang(v))
			}
			if v := c.String("formality"); v != "" {
				opts = append(opts, deepl.WithFormality(deepl.Formality(v)))
			}
			if v := c.String("glossary-id"); v != "" {
				opts = append(opts, deepl.WithGlossaryID(v))
			}
			if v := c.String("context"); v != "" {
				opts = append(opts, deepl.WithTranslationContext(v))
			}
			if v := c.String("split-sentences"); v != "" && v != "1" {
				opts = append(opts, deepl.WithSplitSentences(deepl.SplitSentences(v)))
			}
			if c.Bool("preserve-formatting") {
				opts = append(opts, deepl.WithPreserveFormatting(true))
			}
			if v := c.String("model-type"); v != "" {
				opts = append(opts, deepl.WithModelType(deepl.ModelType(v)))
			}
			if v := c.String("tag-handling"); v != "" {
				opts = append(opts, deepl.WithTagHandling(deepl.TagHandling(v)))
			}
			if v := c.String("tag-handling-version"); v != "" {
				opts = append(opts, deepl.WithTagHandlingVersion(v))
			}
			if instrs := c.StringSlice("custom-instructions"); len(instrs) > 0 {
				opts = append(opts, deepl.WithCustomInstructions(instrs))
			}
			if v := c.String("style-id"); v != "" {
				opts = append(opts, deepl.WithStyleID(v))
			}
			if v := c.String("translation-memory-id"); v != "" {
				opts = append(opts, deepl.WithTranslationMemoryID(v))
			}
			if c.IsSet("translation-memory-threshold") {
				threshold := c.Float64("translation-memory-threshold")
				opts = append(opts, deepl.WithTranslationMemoryThreshold(threshold))
			}

			results, err := client.TranslateText(context.Background(), []string{text}, c.String("to"), opts...)
			if err != nil {
				return err
			}

			rows := entity.TranslationRows(results)
			formatter := output.NewFormatter(c.String("output"))
			formatted := formatter.FormatTranslations(rows)

			if outputFile := c.String("output-file"); outputFile != "" {
				return os.WriteFile(outputFile, []byte(formatted), 0o644)
			}

			fmt.Print(formatted)
			return nil
		},
	}
}
