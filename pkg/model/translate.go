// Package model holds the public, Go-idiomatic domain types returned by the
// deepl.Client. These types are the primary API surface that library consumers
// interact with. They are kept in a separate package from deepl.Client to
// isolate internal wire-format details (internal/apimodel) from the public API.
package model

// TextResult is a single translated text and its metadata, returned by
// Client.TranslateText. The Text field contains the translated string; the
// other fields provide diagnostic and billing information.
//
// Example:
//
//	result := results[0]
//	fmt.Printf("Translated: %s\n", result.Text)
//	fmt.Printf("From: %s\n", result.DetectedSourceLang)
//	fmt.Printf("Billed: %d characters\n", result.BilledCharacters)
type TextResult struct {
	// Text is the translated text segment.
	Text string

	// DetectedSourceLang is the source language code that was auto-detected
	// (e.g., "EN") or the one you explicitly set via WithSourceLang. Always
	// populated after a successful translation.
	DetectedSourceLang string

	// BilledCharacters is the number of characters billed for this translation.
	// May differ from len(Text) due to language-specific character counting rules
	// (e.g., CJK languages, combining diacritics). See DeepL's billing documentation.
	BilledCharacters int

	// ModelTypeUsed reports which model DeepL selected for this request when
	// explicitly expressed via WithModelType. May be empty if DeepL didn't
	// include it in the response.
	ModelTypeUsed string
}
