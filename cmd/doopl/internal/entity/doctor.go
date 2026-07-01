package entity

// DoctorReport represents the result of a health check.
type DoctorReport struct {
	// Basic connectivity
	Connected        bool   `json:"connected"`
	ConnectError     string `json:"connect_error,omitempty"`
	ConnectLatencyMs int64  `json:"connect_latency_ms"`

	// Translation capability
	TranslationWorks     bool   `json:"translation_works"`
	TranslationLatencyMs int64  `json:"translation_latency_ms"`
	DetectedLanguage     string `json:"detected_language,omitempty"`

	// Quota
	CharacterCount int64 `json:"character_count"`
	CharacterLimit int64 `json:"character_limit"`
	DocumentCount  int64 `json:"document_count"`
	DocumentLimit  int64 `json:"document_limit"`

	// Verbose mode details
	Verbose              bool `json:"verbose"`
	SourceLanguagesCount int  `json:"source_languages_count,omitempty"`
	TargetLanguagesCount int  `json:"target_languages_count,omitempty"`
	GlossariesWork       bool `json:"glossaries_work,omitempty"`
	RephraseWorks        bool `json:"rephrase_works,omitempty"`
}
