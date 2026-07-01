package entity

// GlossaryRow represents a glossary for display purposes.
type GlossaryRow struct {
	GlossaryID   string `json:"glossary_id"`
	Name         string `json:"name"`
	SourceLang   string `json:"source_lang"`
	TargetLang   string `json:"target_lang"`
	EntryCount   int    `json:"entry_count"`
	CreationTime string `json:"creation_time"`
}
