package model

// Glossary represents a terminology dictionary that customizes translation behavior.
// Glossaries are language-pair specific.
type Glossary struct {
	// GlossaryID is the unique identifier for the glossary.
	GlossaryID string
	// Name is the human-readable name of the glossary.
	Name string
	// Ready indicates whether the glossary is ready for use.
	Ready bool
	// SourceLang is the language code for source terms.
	SourceLang string
	// TargetLang is the language code for target terms.
	TargetLang string
	// CreationTime is the ISO 8601 timestamp when the glossary was created.
	CreationTime string
	// EntryCount is the number of term pairs in the glossary.
	EntryCount int64
}

// GlossaryEntries represents key-value term pairs for a glossary.
type GlossaryEntries map[string]string
