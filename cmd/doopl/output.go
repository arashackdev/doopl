// Package main provides the doopl CLI.
package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// outputJSON marshals v to JSON with indentation and prints to stdout.
func outputJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

// outputTable prints data as a formatted table.
// Each row is a line of formatted text.
func outputTable(rows []string) {
	for _, row := range rows {
		fmt.Println(row)
	}
}

// outputText prints plain text lines.
func outputText(lines []string) {
	for _, line := range lines {
		fmt.Println(line)
	}
}

// formatTableHeader formats a table header with columns separated by spacing.
// Example: formatTableHeader(5, "CODE", "NAME") -> "CODE  NAME"
func formatTableHeader(columns ...string) string {
	if len(columns) == 0 {
		return ""
	}
	widths := make([]int, len(columns))
	for i, col := range columns {
		widths[i] = len(col)
	}
	return formatTableRow(widths, columns...)
}

// formatTableRow formats a table row with columns aligned to specified widths.
func formatTableRow(widths []int, columns ...string) string {
	if len(columns) == 0 {
		return ""
	}
	parts := make([]string, len(columns))
	for i, col := range columns {
		if i < len(widths) {
			parts[i] = fmt.Sprintf("%-*s", widths[i], col)
		} else {
			parts[i] = col
		}
	}
	return strings.Join(parts, " ")
}

// formatTableDivider creates a horizontal divider line for tables.
func formatTableDivider(width int) string {
	return strings.Repeat("-", width)
}
