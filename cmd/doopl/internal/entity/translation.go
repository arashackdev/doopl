// Package entity holds presentation-layer types used only by cmd/doopl —
// shaped for terminal/JSON output, not for library consumption. Keeping
// these separate from model means CLI-specific concerns (column ordering,
// display formatting) never pressure the library's public API.
package entity

// TranslationRow is one row of `deepl translate` output.
type TranslationRow struct {
	Text               string `json:"text"`
	DetectedSourceLang string `json:"detected_source_lang"`
	BilledCharacters   int    `json:"billed_characters"`
}
