package apimodel

// LanguagesResponse is the wire format for the /v3/languages endpoint.
type LanguagesResponse struct {
	Languages []Language `json:"languages"`
}

// Language represents a language returned by the languages endpoint.
type Language struct {
	Code              string `json:"language"`
	Name              string `json:"name"`
	SupportsFormality *bool  `json:"supports_formality"`
}

// UsageResponse is the wire format for the /v3/usage endpoint.
type UsageResponse struct {
	CharacterCount    int64 `json:"character_count"`
	CharacterLimit    int64 `json:"character_limit"`
	DocumentCount     int64 `json:"document_count"`
	DocumentLimit     int64 `json:"document_limit"`
	TeamDocumentCount int64 `json:"team_document_count"`
	TeamDocumentLimit int64 `json:"team_document_limit"`
}
