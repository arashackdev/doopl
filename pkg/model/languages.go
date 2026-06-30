package model

// Language represents a supported language for DeepL translations.
// Some languages support formality control (see SupportsFormality).
type Language struct {
	// Code is the language code (e.g., "EN", "DE", "FR").
	Code string
	// Name is the human-readable language name (e.g., "English", "German").
	Name string
	// SupportsFormality indicates whether this language supports formal/informal speech (nil if unknown).
	SupportsFormality *bool
}

// Usage represents the current API quota and consumption.
// Character limits are per billing period; document limits are concurrent/active.
type Usage struct {
	// CharacterCount is the number of characters translated in the current billing period.
	CharacterCount int64
	// CharacterLimit is the maximum characters allowed in the current billing period.
	CharacterLimit int64
	// DocumentCount is the number of documents currently being processed or queued.
	DocumentCount int64
	// DocumentLimit is the maximum concurrent documents allowed.
	DocumentLimit int64
	// TeamDocumentCount is the number of documents being processed by the team (if applicable).
	TeamDocumentCount int64
	// TeamDocumentLimit is the maximum concurrent documents allowed for the team (if applicable).
	TeamDocumentLimit int64
}
