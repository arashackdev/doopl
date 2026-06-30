package deepl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/arashackdev/doopl/internal/apimodel"
	"github.com/arashackdev/doopl/internal/convert"
	"github.com/arashackdev/doopl/pkg/model"
)

// apiToModel is the package-level converter instance used by every method
// in this file. goverter generates a zero-field struct, so a package var is
// safe for concurrent use — same contract as the rest of the Client.
var apiToModel = &convert.APIToModelImpl{}

// Formality controls the formality level of the translation. Only honored
// for target languages that support it.
type Formality string

// Formality values for the formality option.
const (
	FormalityDefault    Formality = "default"
	FormalityMore       Formality = "more"
	FormalityLess       Formality = "less"
	FormalityPreferMore Formality = "prefer_more"
	FormalityPreferLess Formality = "prefer_less"
)

// SplitSentences controls how the input is split into sentences before translation.
type SplitSentences string

// SplitSentences values.
const (
	SplitSentencesOff        SplitSentences = "0"
	SplitSentencesOn         SplitSentences = "1"
	SplitSentencesNoNewlines SplitSentences = "nonewlines"
)

// ModelType selects between translation quality and latency. DeepL chooses
// the underlying model; this only expresses your preference, not a specific
// model version — see the API docs for why that's intentional.
type ModelType string

// ModelType values.
const (
	ModelTypeQualityOptimized ModelType = "quality_optimized"
	ModelTypeLatencyOptimized ModelType = "latency_optimized"
	ModelTypePreferQuality    ModelType = "prefer_quality_optimized"
)

// TagHandling specifies how XML/HTML tags in the input should be handled.
type TagHandling string

// TagHandling values.
const (
	TagHandlingXML  TagHandling = "xml"
	TagHandlingHTML TagHandling = "html"
)

type translateTextParams struct {
	sourceLang                 string
	formality                  Formality
	glossaryID                 string
	context                    string
	splitSentences             SplitSentences
	preserveFormatting         *bool
	modelType                  ModelType
	tagHandling                TagHandling
	tagHandlingVersion         string
	customInstructions         []string
	styleID                    string
	translationMemoryID        string
	translationMemoryThreshold *float64
}

// TranslateTextOption configures an individual TranslateText call.
type TranslateTextOption func(*translateTextParams)

// WithSourceLang sets the source language explicitly. If omitted, DeepL
// auto-detects it (and the result's DetectedSourceLang reports what it found).
func WithSourceLang(lang string) TranslateTextOption {
	return func(p *translateTextParams) { p.sourceLang = lang }
}

// WithFormality sets the formality level for languages that support it.
func WithFormality(f Formality) TranslateTextOption {
	return func(p *translateTextParams) { p.formality = f }
}

// WithGlossaryID applies a previously-created glossary to the translation.
// Requires WithSourceLang to also be set, matching the API's requirement
// that glossary use pins the source language.
func WithGlossaryID(id string) TranslateTextOption {
	return func(p *translateTextParams) { p.glossaryID = id }
}

// WithTranslationContext supplies additional text used only to improve
// translation quality of the main text — the context itself is not
// translated or returned. See DeepL's context parameter guide.
func WithTranslationContext(context string) TranslateTextOption {
	return func(p *translateTextParams) { p.context = context }
}

// WithSplitSentences controls sentence-splitting behavior on the input.
func WithSplitSentences(s SplitSentences) TranslateTextOption {
	return func(p *translateTextParams) { p.splitSentences = s }
}

// WithPreserveFormatting disables DeepL's automatic formatting corrections
// (e.g. punctuation spacing) when set to true.
func WithPreserveFormatting(preserve bool) TranslateTextOption {
	return func(p *translateTextParams) { p.preserveFormatting = &preserve }
}

// WithModelType expresses a preference for translation quality vs latency.
// DeepL still selects the specific underlying model.
func WithModelType(m ModelType) TranslateTextOption {
	return func(p *translateTextParams) { p.modelType = m }
}

// WithTagHandling enables XML/HTML tag-aware translation.
func WithTagHandling(t TagHandling) TranslateTextOption {
	return func(p *translateTextParams) { p.tagHandling = t }
}

// WithTagHandlingVersion selects the tag-handling algorithm version ("v1" or "v2").
func WithTagHandlingVersion(version string) TranslateTextOption {
	return func(p *translateTextParams) { p.tagHandlingVersion = version }
}

// WithCustomInstructions provides up to 10 natural-language instructions
// (max 300 chars each) to customize translation behavior, e.g.
// []string{"Use a friendly, diplomatic tone"}. Only supported for a subset
// of target languages, and forces ModelTypeQualityOptimized — see DeepL docs.
func WithCustomInstructions(instructions []string) TranslateTextOption {
	return func(p *translateTextParams) { p.customInstructions = instructions }
}

// WithStyleID applies a configured style rule list to the translation.
func WithStyleID(id string) TranslateTextOption {
	return func(p *translateTextParams) { p.styleID = id }
}

// WithTranslationMemoryID and WithTranslationMemoryThreshold control use of
// a translation memory for this request.
func WithTranslationMemoryID(id string) TranslateTextOption {
	return func(p *translateTextParams) { p.translationMemoryID = id }
}

// WithTranslationMemoryThreshold sets the minimum match score (0.0-1.0)
// required for a translation memory match to be used.
func WithTranslationMemoryThreshold(threshold float64) TranslateTextOption {
	return func(p *translateTextParams) { p.translationMemoryThreshold = &threshold }
}

// TranslateText translates one or more texts into targetLang. Source
// language is auto-detected unless WithSourceLang is supplied.
//
//	results, err := client.TranslateText(ctx, []string{"Hello, world!"}, "DE")
//	results[0].Text // "Hallo, Welt!"
func (c *Client) TranslateText(ctx context.Context, texts []string, targetLang string, opts ...TranslateTextOption) ([]model.TextResult, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("deepl: texts must not be empty")
	}
	if targetLang == "" {
		return nil, fmt.Errorf("deepl: targetLang must not be empty")
	}

	p := &translateTextParams{}
	for _, opt := range opts {
		opt(p)
	}

	reqBody := apimodel.TranslateRequest{
		Text:                       texts,
		TargetLang:                 targetLang,
		SourceLang:                 p.sourceLang,
		Formality:                  string(p.formality),
		GlossaryID:                 p.glossaryID,
		Context:                    p.context,
		SplitSentences:             string(p.splitSentences),
		PreserveFormatting:         p.preserveFormatting,
		ModelType:                  string(p.modelType),
		TagHandling:                string(p.tagHandling),
		TagHandlingVersion:         p.tagHandlingVersion,
		CustomInstructions:         p.customInstructions,
		StyleID:                    p.styleID,
		TranslationMemoryID:        p.translationMemoryID,
		TranslationMemoryThreshold: p.translationMemoryThreshold,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("deepl: encoding request: %w", err)
	}

	resp, err := c.transport.Do(ctx, transportRequest("POST", "/v2/translate", payload, "application/json"))
	if err != nil {
		return nil, fmt.Errorf("deepl: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("deepl: reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newAPIError(resp.StatusCode, extractMessage(body))
	}

	var result apimodel.TranslateResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("deepl: decoding response: %w", err)
	}

	return apiToModel.Translations(result.Translations), nil
}
