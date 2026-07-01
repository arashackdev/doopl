// Package apimodel contains wire-format types that mirror DeepL's JSON API
// responses and request bodies exactly (snake_case json tags, raw string
// enums, etc.). These types are never exposed to library consumers — see
// internal/convert for the goverter-generated mapping into the public
// github.com/arashackdev/doopl/model types.
//
// Keeping this layer separate means a quirky or churny wire shape never
// leaks into the public API, and the OpenAPI-spec-drift check (see
// docs/build-plan.md) only ever has to touch this package.
package apimodel

// TranslateRequest is the request body for POST /v2/translate.
type TranslateRequest struct {
	Text                       []string `json:"text"`
	TargetLang                 string   `json:"target_lang"`
	SourceLang                 string   `json:"source_lang,omitempty"`
	Formality                  string   `json:"formality,omitempty"`
	GlossaryID                 string   `json:"glossary_id,omitempty"`
	Context                    string   `json:"context,omitempty"`
	SplitSentences             string   `json:"split_sentences,omitempty"`
	PreserveFormatting         *bool    `json:"preserve_formatting,omitempty"`
	ModelType                  string   `json:"model_type,omitempty"`
	TagHandling                string   `json:"tag_handling,omitempty"`
	TagHandlingVersion         string   `json:"tag_handling_version,omitempty"`
	CustomInstructions         []string `json:"custom_instructions,omitempty"`
	StyleID                    string   `json:"style_id,omitempty"`
	TranslationMemoryID        string   `json:"translation_memory_id,omitempty"`
	TranslationMemoryThreshold *float64 `json:"translation_memory_threshold,omitempty"`
}

// TranslateResponse is the response body for POST /v2/translate.
type TranslateResponse struct {
	Translations []Translation `json:"translations"`
}

// Translation is a single translated text as returned by the API, before
// conversion into the public model.TextResult.
type Translation struct {
	Text               string `json:"text"`
	DetectedSourceLang string `json:"detected_source_language"`
	BilledCharacters   int    `json:"billed_characters"`
	ModelTypeUsed      string `json:"model_type_used,omitempty"`
}

// ErrorResponse mirrors DeepL's JSON error body shape:
// {"message": "...", "detail": "..."}.
type ErrorResponse struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
}
