// Package model holds the public, Go-idiomatic domain types returned by the
// deepl.Client — as opposed to internal/apimodel, which mirrors the raw
// DeepL JSON wire format. Library consumers interact with these types, not
// apimodel's.
package model

// TextResult is a single translated text and its metadata, returned by
// Client.TranslateText.
type TextResult struct {
	// Text is the translated text.
	Text string
	// DetectedSourceLang is the source language DeepL detected (or the one
	// you explicitly supplied via WithSourceLang).
	DetectedSourceLang string
	// BilledCharacters is the number of characters billed for this
	// translation — may differ from len(Text) for some languages/scripts.
	BilledCharacters int
	// ModelTypeUsed reports which underlying model served the request, when
	// DeepL includes it in the response.
	ModelTypeUsed string
}
