package entity

// LanguageRow represents a language for display purposes (table/JSON output).
type LanguageRow struct {
	Code              string `json:"code"`
	Name              string `json:"name"`
	SupportsFormality *bool  `json:"supports_formality"`
}

// UsageRow represents API usage for display purposes.
type UsageRow struct {
	CharacterCount    int64 `json:"character_count"`
	CharacterLimit    int64 `json:"character_limit"`
	DocumentCount     int64 `json:"document_count"`
	DocumentLimit     int64 `json:"document_limit"`
	TeamDocumentCount int64 `json:"team_document_count"`
	TeamDocumentLimit int64 `json:"team_document_limit"`
}
