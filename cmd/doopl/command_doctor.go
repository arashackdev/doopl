// Package main provides the doopl CLI.
package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/arashackdev/doopl/cmd/doopl/internal/entity"
	"github.com/arashackdev/doopl/cmd/doopl/internal/output"
	deepl "github.com/arashackdev/doopl/pkg/deepl"
	"github.com/urfave/cli/v2"
)

// doctorCommand returns the doctor command for checking API health.
func doctorCommand() *cli.Command {
	return &cli.Command{
		Name:  "doctor",
		Usage: "check API connectivity and health",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "run full diagnostic suite (languages, glossaries, rephrase)",
			},
		},
		Action: func(c *cli.Context) error {
			authKey, err := getAuthKey(c)
			if err != nil {
				return err
			}

			formatter := output.NewFormatter(c.String("output"))
			report := runDoctorChecks(context.Background(), authKey, c.Bool("verbose"))
			fmt.Println(formatter.FormatDoctorReport(report))
			return nil
		},
	}
}

// runDoctorChecks performs health checks on the API.
func runDoctorChecks(ctx context.Context, authKey string, verbose bool) entity.DoctorReport {
	report := entity.DoctorReport{
		Verbose: verbose,
	}

	client, err := deepl.New(authKey, deepl.WithAppInfo("doopl-cli", deepl.Version))
	if err != nil {
		report.Connected = false
		report.ConnectError = err.Error()
		return report
	}

	// Test 1: Basic connectivity with a simple translation
	start := time.Now()
	results, err := client.TranslateText(ctx, []string{"Hello"}, "DE")
	report.ConnectLatencyMs = int64(time.Since(start).Milliseconds())

	if err != nil {
		report.Connected = false
		report.ConnectError = err.Error()
		return report
	}

	report.Connected = true
	report.TranslationWorks = true

	if len(results) > 0 {
		report.DetectedLanguage = results[0].DetectedSourceLang
	}

	// Get latency for translation (already measured above)
	report.TranslationLatencyMs = report.ConnectLatencyMs

	// Test 2: Check quota
	start = time.Now()
	usage, err := client.Usage(ctx)
	_ = time.Since(start).Milliseconds() // Could add usage_latency_ms if needed

	if err == nil {
		report.CharacterCount = usage.CharacterCount
		report.CharacterLimit = usage.CharacterLimit
		report.DocumentCount = usage.DocumentCount
		report.DocumentLimit = usage.DocumentLimit
	}

	// Verbose mode: run additional checks
	if verbose {
		runVerboseChecks(ctx, client, &report)
	}

	return report
}

// runVerboseChecks performs additional health checks for verbose mode.
func runVerboseChecks(ctx context.Context, client *deepl.Client, report *entity.DoctorReport) {
	// Test 3: Check source languages
	srcLangs, err := client.SourceLanguages(ctx)
	if err == nil {
		report.SourceLanguagesCount = len(srcLangs)
	}

	// Test 4: Check target languages
	tgtLangs, err := client.TargetLanguages(ctx)
	if err == nil {
		report.TargetLanguagesCount = len(tgtLangs)
	}

	// Test 5: Check glossaries
	glossaries, err := client.ListGlossaries(ctx)
	report.GlossariesWork = (err == nil)
	_ = glossaries // Suppress unused warning

	// Test 6: Check rephrase (Write API)
	_, err = client.Rephrase(ctx, []string{"Hello world"}, "EN")
	report.RephraseWorks = (err == nil && !errors.Is(err, deepl.ErrQuotaExceeded))
}
