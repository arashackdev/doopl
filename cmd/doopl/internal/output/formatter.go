// Package output provides formatting for CLI output in text, TUI, and JSON modes.
package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/arashackdev/doopl/cmd/doopl/internal/entity"
	"github.com/charmbracelet/lipgloss"
)

// Formatter formats output data in different modes.
type Formatter interface {
	FormatTranslations(rows []entity.TranslationRow) string
	FormatLanguages(rows []entity.LanguageRow) string
	FormatUsage(usage entity.UsageRow) string
	FormatGlossaries(rows []entity.GlossaryRow) string
	FormatGlossaryEntries(entries map[string]string) string
	FormatDoctorReport(report entity.DoctorReport) string
}

// NewFormatter creates a formatter for the given mode.
func NewFormatter(mode string) Formatter {
	switch mode {
	case "json":
		return &JSONFormatter{}
	case "tui":
		return &TUIFormatter{}
	default: // text
		return &TextFormatter{}
	}
}

// TextFormatter outputs plain text.
type TextFormatter struct{}

// FormatTranslations formats translations as plain text.
func (f *TextFormatter) FormatTranslations(rows []entity.TranslationRow) string {
	var sb strings.Builder
	for i, row := range rows {
		if i > 0 {
			sb.WriteString("\n")
		}
		fmt.Fprintf(&sb, "Text: %s\n", row.Text)
		fmt.Fprintf(&sb, "Detected Lang: %s\n", row.DetectedSourceLang)
		fmt.Fprintf(&sb, "Billed Chars: %d\n", row.BilledCharacters)
	}
	return sb.String()
}

// FormatLanguages formats languages as plain text.
func (f *TextFormatter) FormatLanguages(rows []entity.LanguageRow) string {
	var sb strings.Builder
	for i, row := range rows {
		if i > 0 {
			sb.WriteString("\n")
		}
		fmt.Fprintf(&sb, "[%s] %s", row.Code, row.Name)
		if row.SupportsFormality != nil && *row.SupportsFormality {
			sb.WriteString(" (formality supported)")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// FormatUsage formats usage as plain text.
func (f *TextFormatter) FormatUsage(usage entity.UsageRow) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Characters: %d / %d\n", usage.CharacterCount, usage.CharacterLimit)
	fmt.Fprintf(&sb, "Documents: %d / %d\n", usage.DocumentCount, usage.DocumentLimit)
	if usage.TeamDocumentCount > 0 || usage.TeamDocumentLimit > 0 {
		fmt.Fprintf(&sb, "Team Documents: %d / %d\n", usage.TeamDocumentCount, usage.TeamDocumentLimit)
	}
	return sb.String()
}

// FormatGlossaries formats glossaries as plain text.
func (f *TextFormatter) FormatGlossaries(rows []entity.GlossaryRow) string {
	var sb strings.Builder
	for i, row := range rows {
		if i > 0 {
			sb.WriteString("\n")
		}
		fmt.Fprintf(&sb, "ID: %s\n", row.GlossaryID)
		fmt.Fprintf(&sb, "Name: %s\n", row.Name)
		fmt.Fprintf(&sb, "Language Pair: %s → %s\n", row.SourceLang, row.TargetLang)
		fmt.Fprintf(&sb, "Entries: %d\n", row.EntryCount)
		fmt.Fprintf(&sb, "Created: %s\n", row.CreationTime)
	}
	return sb.String()
}

// FormatGlossaryEntries formats glossary entries as plain text.
func (f *TextFormatter) FormatGlossaryEntries(entries map[string]string) string {
	var sb strings.Builder
	for source, target := range entries {
		fmt.Fprintf(&sb, "%s → %s\n", source, target)
	}
	return sb.String()
}

// FormatDoctorReport formats a doctor report as plain text.
func (f *TextFormatter) FormatDoctorReport(report entity.DoctorReport) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "API Health Report\n")
	fmt.Fprintf(&sb, "=================\n\n")

	fmt.Fprintf(&sb, "Connectivity:\n")
	status := "✗ FAIL"
	if report.Connected {
		status = "✓ PASS"
	}
	fmt.Fprintf(&sb, "  %s\n", status)
	if report.ConnectError != "" {
		fmt.Fprintf(&sb, "  Error: %s\n", report.ConnectError)
	}
	fmt.Fprintf(&sb, "  Latency: %dms\n\n", report.ConnectLatencyMs)

	fmt.Fprintf(&sb, "Translation Test:\n")
	status = "✗ FAIL"
	if report.TranslationWorks {
		status = "✓ PASS"
	}
	fmt.Fprintf(&sb, "  %s\n", status)
	if report.TranslationLatencyMs > 0 {
		fmt.Fprintf(&sb, "  Latency: %dms\n", report.TranslationLatencyMs)
	}

	fmt.Fprintf(&sb, "\nQuota:\n")
	fmt.Fprintf(&sb, "  Characters: %d / %d\n", report.CharacterCount, report.CharacterLimit)
	fmt.Fprintf(&sb, "  Documents: %d / %d\n", report.DocumentCount, report.DocumentLimit)

	if report.Verbose {
		fmt.Fprintf(&sb, "\nLanguages: %d sources, %d targets\n", report.SourceLanguagesCount, report.TargetLanguagesCount)
		if report.GlossariesWork {
			fmt.Fprintf(&sb, "Glossaries: ✓ Available\n")
		} else {
			fmt.Fprintf(&sb, "Glossaries: ✗ Not available\n")
		}
		if report.RephraseWorks {
			fmt.Fprintf(&sb, "Rephrase (Write API): ✓ Available\n")
		} else {
			fmt.Fprintf(&sb, "Rephrase (Write API): ✗ Not available\n")
		}
		fmt.Fprintf(&sb, "Detected Language: %s\n", report.DetectedLanguage)
	}

	return sb.String()
}

// JSONFormatter outputs JSON.
type JSONFormatter struct{}

// FormatTranslations formats translations as JSON.
func (f *JSONFormatter) FormatTranslations(rows []entity.TranslationRow) string {
	data, _ := json.MarshalIndent(rows, "", "  ")
	return string(data)
}

// FormatLanguages formats languages as JSON.
func (f *JSONFormatter) FormatLanguages(rows []entity.LanguageRow) string {
	data, _ := json.MarshalIndent(rows, "", "  ")
	return string(data)
}

// FormatUsage formats usage as JSON.
func (f *JSONFormatter) FormatUsage(usage entity.UsageRow) string {
	data, _ := json.MarshalIndent(usage, "", "  ")
	return string(data)
}

// FormatGlossaries formats glossaries as JSON.
func (f *JSONFormatter) FormatGlossaries(rows []entity.GlossaryRow) string {
	data, _ := json.MarshalIndent(rows, "", "  ")
	return string(data)
}

// FormatGlossaryEntries formats glossary entries as JSON.
func (f *JSONFormatter) FormatGlossaryEntries(entries map[string]string) string {
	data, _ := json.MarshalIndent(entries, "", "  ")
	return string(data)
}

// FormatDoctorReport formats a doctor report as JSON.
func (f *JSONFormatter) FormatDoctorReport(report entity.DoctorReport) string {
	data, _ := json.MarshalIndent(report, "", "  ")
	return string(data)
}

// TUIFormatter outputs rich terminal UI with lipgloss.
type TUIFormatter struct{}

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("4")).
			Padding(0, 1).
			MarginBottom(1)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	tableRowStyle = lipgloss.NewStyle().
			Padding(0, 1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("4")).
			Padding(1).
			MarginBottom(1)
)

// FormatTranslations formats translations as TUI.
func (f *TUIFormatter) FormatTranslations(rows []entity.TranslationRow) string {
	var sb strings.Builder
	sb.WriteString(headerStyle.Render("Translations") + "\n\n")

	for i, row := range rows {
		if i > 0 {
			sb.WriteString("\n" + strings.Repeat("-", 60) + "\n\n")
		}
		content := fmt.Sprintf(
			"Text: %s\nDetected Lang: %s\nBilled Chars: %d",
			row.Text,
			row.DetectedSourceLang,
			row.BilledCharacters,
		)
		sb.WriteString(boxStyle.Render(content))
	}
	return sb.String()
}

// FormatLanguages formats languages as TUI.
func (f *TUIFormatter) FormatLanguages(rows []entity.LanguageRow) string {
	var sb strings.Builder
	sb.WriteString(headerStyle.Render("Languages") + "\n\n")

	for _, row := range rows {
		line := fmt.Sprintf("[%s] %s", row.Code, row.Name)
		if row.SupportsFormality != nil && *row.SupportsFormality {
			line += " " + dimStyle.Render("(formality)")
		}
		sb.WriteString(tableRowStyle.Render(line) + "\n")
	}
	return sb.String()
}

// FormatUsage formats usage as TUI.
func (f *TUIFormatter) FormatUsage(usage entity.UsageRow) string {
	var sb strings.Builder
	sb.WriteString(headerStyle.Render("API Usage") + "\n\n")

	charsPct := float64(usage.CharacterCount) / float64(usage.CharacterLimit) * 100
	docsPct := float64(usage.DocumentCount) / float64(usage.DocumentLimit) * 100

	content := fmt.Sprintf(
		"Characters: %d / %d (%.1f%%)\nDocuments: %d / %d (%.1f%%)",
		usage.CharacterCount,
		usage.CharacterLimit,
		charsPct,
		usage.DocumentCount,
		usage.DocumentLimit,
		docsPct,
	)
	if usage.TeamDocumentCount > 0 || usage.TeamDocumentLimit > 0 {
		teamPct := float64(usage.TeamDocumentCount) / float64(usage.TeamDocumentLimit) * 100
		content += fmt.Sprintf("\nTeam Documents: %d / %d (%.1f%%)",
			usage.TeamDocumentCount,
			usage.TeamDocumentLimit,
			teamPct,
		)
	}

	sb.WriteString(boxStyle.Render(content))
	return sb.String()
}

// FormatGlossaries formats glossaries as TUI.
func (f *TUIFormatter) FormatGlossaries(rows []entity.GlossaryRow) string {
	var sb strings.Builder
	sb.WriteString(headerStyle.Render("Glossaries") + "\n\n")

	for i, row := range rows {
		if i > 0 {
			sb.WriteString("\n")
		}
		line := fmt.Sprintf(
			"%s → %s\n%s (%s) %d entries\n%s",
			row.SourceLang,
			row.TargetLang,
			row.Name,
			row.GlossaryID,
			row.EntryCount,
			dimStyle.Render(row.CreationTime),
		)
		sb.WriteString(boxStyle.Render(line))
	}
	return sb.String()
}

// FormatGlossaryEntries formats glossary entries as TUI.
func (f *TUIFormatter) FormatGlossaryEntries(entries map[string]string) string {
	var sb strings.Builder
	sb.WriteString(headerStyle.Render("Entries") + "\n\n")

	for source, target := range entries {
		line := fmt.Sprintf("%s %s %s", source, dimStyle.Render("→"), target)
		sb.WriteString(tableRowStyle.Render(line) + "\n")
	}
	return sb.String()
}

// FormatDoctorReport formats a doctor report as TUI.
func (f *TUIFormatter) FormatDoctorReport(report entity.DoctorReport) string {
	var sb strings.Builder
	sb.WriteString(headerStyle.Render("🏥 API Health Report") + "\n\n")

	connStatus := errorStyle.Render("✗ FAIL")
	if report.Connected {
		connStatus = successStyle.Render("✓ PASS")
	}
	connBox := fmt.Sprintf("Connectivity: %s\nLatency: %dms", connStatus, report.ConnectLatencyMs)
	if report.ConnectError != "" {
		connBox += fmt.Sprintf("\nError: %s", report.ConnectError)
	}
	sb.WriteString(boxStyle.Render(connBox) + "\n")

	transStatus := errorStyle.Render("✗ FAIL")
	if report.TranslationWorks {
		transStatus = successStyle.Render("✓ PASS")
	}
	transBox := fmt.Sprintf("Translation Test: %s\nLatency: %dms", transStatus, report.TranslationLatencyMs)
	sb.WriteString(boxStyle.Render(transBox) + "\n")

	charsPct := float64(report.CharacterCount) / float64(report.CharacterLimit) * 100
	docsPct := float64(report.DocumentCount) / float64(report.DocumentLimit) * 100
	quotaBox := fmt.Sprintf(
		"Characters: %d / %d (%.1f%%)\nDocuments: %d / %d (%.1f%%)",
		report.CharacterCount,
		report.CharacterLimit,
		charsPct,
		report.DocumentCount,
		report.DocumentLimit,
		docsPct,
	)
	sb.WriteString(boxStyle.Render(quotaBox))

	if report.Verbose {
		sb.WriteString("\n")
		verboseBox := fmt.Sprintf(
			"Languages: %d sources, %d targets\nGlossaries: %v\nRephrase: %v\nDetected: %s",
			report.SourceLanguagesCount,
			report.TargetLanguagesCount,
			report.GlossariesWork,
			report.RephraseWorks,
			report.DetectedLanguage,
		)
		sb.WriteString(boxStyle.Render(verboseBox))
	}

	return sb.String()
}
