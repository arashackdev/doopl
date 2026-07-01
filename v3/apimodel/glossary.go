package apimodel

// GlossaryResponse represents a glossary in the API.
type GlossaryResponse struct {
	GlossaryID   string `json:"glossary_id"`
	Name         string `json:"name"`
	Ready        bool   `json:"ready"`
	SourceLang   string `json:"source_lang"`
	TargetLang   string `json:"target_lang"`
	CreationTime string `json:"creation_time"`
	EntryCount   int64  `json:"entry_count"`
}

// GlossariesResponse lists glossaries.
type GlossariesResponse struct {
	Glossaries []GlossaryResponse `json:"glossaries"`
}

// GlossaryEntriesResponse contains entries in a glossary.
type GlossaryEntriesResponse struct {
	Entries map[string]string `json:"entries"`
}
