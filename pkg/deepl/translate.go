package deepl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/arashackdev/doopl/internal/convert"
	"github.com/arashackdev/doopl/pkg/model"
	"github.com/arashackdev/doopl/v3/apimodel"
)

// apiToModel is the package-level converter instance used by every method
// in this file. goverter generates a zero-field struct, so a package var is
// safe for concurrent use — same contract as the rest of the Client.
var apiToModel = &convert.APIToModelImpl{}

// Formality controls the formality level of the translation. Only honored for
// target languages that support it (e.g., most European languages, but not
// Chinese or Japanese). When not supported for a target language, this option
// is silently ignored.
type Formality string

// Formality level options. The "prefer_*" variants tell DeepL to use the
// specified formality level if available, otherwise use a neutral formality.
// The strict "more"/"less" variants will error if the target language does not
// support the requested formality.
const (
	FormalityDefault    Formality = "default"     // Use DeepL's default formality
	FormalityMore       Formality = "more"        // Formal (error if unsupported)
	FormalityLess       Formality = "less"        // Informal (error if unsupported)
	FormalityPreferMore Formality = "prefer_more" // Prefer formal, fallback to neutral
	FormalityPreferLess Formality = "prefer_less" // Prefer informal, fallback to neutral
)

// SplitSentences controls how the input is split into sentences before
// translation. Sentence segmentation can affect translation quality and
// consistency. Use this to handle special formatting or preserve structure.
type SplitSentences string

// SplitSentences options.
const (
	SplitSentencesOff        SplitSentences = "0"          // Do not split (translate whole input as one segment)
	SplitSentencesOn         SplitSentences = "1"          // Split on sentences and newlines (default)
	SplitSentencesNoNewlines SplitSentences = "nonewlines" // Split on sentences, preserve newlines
)

// ModelType expresses a preference between translation quality and latency.
// Note that DeepL chooses the specific underlying model dynamically; this
// only communicates your preference, not a specific model version. See DeepL's
// API documentation for current model selections.
type ModelType string

// ModelType options.
const (
	ModelTypeQualityOptimized ModelType = "quality_optimized"        // Prioritize quality (slower)
	ModelTypeLatencyOptimized ModelType = "latency_optimized"        // Prioritize speed (lower latency)
	ModelTypePreferQuality    ModelType = "prefer_quality_optimized" // Prefer quality, allow fallback to latency
)

// TagHandling specifies how XML/HTML tags in the input should be handled.
// When enabled, tags are preserved in the output and their content is
// translated appropriately.
type TagHandling string

// TagHandling options.
const (
	TagHandlingXML  TagHandling = "xml"  // Preserve XML tags in translation
	TagHandlingHTML TagHandling = "html" // Preserve HTML tags in translation
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
// auto-detects the source language, and the result's DetectedSourceLang field
// reports what was found. Specifying the source can improve performance and
// accuracy if you know the source language.
//
// Example: Force translation from English even if auto-detection might fail
//
//	deepl.WithSourceLang("EN")
func WithSourceLang(lang string) TranslateTextOption {
	return func(p *translateTextParams) { p.sourceLang = lang }
}

// WithFormality sets the formality level for languages that support it
// (e.g., German, French, Spanish). For languages that don't support formality,
// the option is silently ignored. Use FormalityPreferMore/PreferLess if you
// want a graceful fallback rather than an error.
func WithFormality(f Formality) TranslateTextOption {
	return func(p *translateTextParams) { p.formality = f }
}

// WithGlossaryID applies a previously-created glossary to the translation.
// Important: Using a glossary requires WithSourceLang to also be set, matching
// the DeepL API's requirement that glossary use pins the source language.
// Returns an error if the source language is not specified.
//
// Example:
//
//	deepl.WithSourceLang("EN"),
//	deepl.WithGlossaryID("my-glossary-id")
func WithGlossaryID(id string) TranslateTextOption {
	return func(p *translateTextParams) { p.glossaryID = id }
}

// WithTranslationContext supplies additional background text to improve
// translation quality of the main text. The context itself is not translated
// or included in results — it is used only to influence the translation of
// the main input. Useful for domain-specific terminology or style guidance.
//
// Example:
//
//	deepl.WithTranslationContext("We translate between programming languages")
func WithTranslationContext(context string) TranslateTextOption {
	return func(p *translateTextParams) { p.context = context }
}

// WithSplitSentences controls how the input is segmented into sentences before
// translation. By default, sentences are split and newlines are treated as
// sentence boundaries. Use SplitSentencesOff to translate the entire input as
// a single unit, or SplitSentencesNoNewlines to preserve input formatting.
func WithSplitSentences(s SplitSentences) TranslateTextOption {
	return func(p *translateTextParams) { p.splitSentences = s }
}

// WithPreserveFormatting disables DeepL's automatic formatting corrections when
// set to true. By default (false), DeepL normalizes spacing and punctuation.
// Set to true to keep the exact formatting of the input text.
func WithPreserveFormatting(preserve bool) TranslateTextOption {
	return func(p *translateTextParams) { p.preserveFormatting = &preserve }
}

// WithModelType expresses a preference between translation quality and latency.
// DeepL selects the actual underlying model dynamically; this only communicates
// your preference. Performance and model selection may vary over time.
func WithModelType(m ModelType) TranslateTextOption {
	return func(p *translateTextParams) { p.modelType = m }
}

// WithTagHandling enables XML or HTML tag-aware translation. When set, tags are
// preserved in the output and not translated. Choose TagHandlingXML or
// TagHandlingHTML based on your input format.
func WithTagHandling(t TagHandling) TranslateTextOption {
	return func(p *translateTextParams) { p.tagHandling = t }
}

// WithTagHandlingVersion selects the tag-handling algorithm version. Most users
// should omit this; it defaults to v2. Use "v1" only if your XML/HTML structure
// requires the v1 parser behavior.
func WithTagHandlingVersion(version string) TranslateTextOption {
	return func(p *translateTextParams) { p.tagHandlingVersion = version }
}

// WithCustomInstructions provides up to 10 natural-language instructions
// (each max 300 characters) to customize translation behavior. For example:
//
//	[]string{"Use a friendly tone", "Prefer 'software' over 'app'"}
//
// Custom instructions are only supported for a subset of target languages and
// automatically enforce ModelTypeQualityOptimized. See DeepL's documentation for
// supported languages and instruction best practices.
func WithCustomInstructions(instructions []string) TranslateTextOption {
	return func(p *translateTextParams) { p.customInstructions = instructions }
}

// WithStyleID applies a style rule list to the translation. Requires a style ID
// created via the DeepL platform. Only supported for a subset of language pairs.
func WithStyleID(id string) TranslateTextOption {
	return func(p *translateTextParams) { p.styleID = id }
}

// WithTranslationMemoryID applies a translation memory to the request. When a
// translation memory is set, exact and fuzzy matches are checked and used if
// their scores meet or exceed WithTranslationMemoryThreshold.
func WithTranslationMemoryID(id string) TranslateTextOption {
	return func(p *translateTextParams) { p.translationMemoryID = id }
}

// WithTranslationMemoryThreshold sets the minimum match score (0.0–1.0) for
// translation memory matches to be used. Defaults to 0.0 (all matches used).
// Increase this to accept only high-confidence matches from the translation memory.
// Only meaningful when WithTranslationMemoryID is also set.
func WithTranslationMemoryThreshold(threshold float64) TranslateTextOption {
	return func(p *translateTextParams) { p.translationMemoryThreshold = &threshold }
}

// TranslateText translates one or more texts into the target language. The
// source language is auto-detected unless explicitly specified via
// WithSourceLang. The order and count of results matches the input texts.
//
// Parameters:
//   - ctx: Context controlling the request lifetime. If canceled, the in-flight
//     request is abandoned and its context error is returned.
//   - texts: One or more texts to translate (must not be empty).
//   - targetLang: Target language code (e.g., "DE", "FR", "ES"). See
//     Client.TargetLanguages for supported values.
//   - opts: Optional parameters like WithFormality, WithGlossaryID, etc.
//
// Returns one TextResult per input text, preserving order. Each result includes
// the translated text, detected source language (if not explicitly set), and
// billing information.
//
// Example:
//
//	results, err := client.TranslateText(ctx,
//		[]string{
//			"Hello, world!",
//			"Good morning!",
//		},
//		"DE",
//		deepl.WithFormality(deepl.FormalityMore),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for i, r := range results {
//		fmt.Printf("[%d] %s (from %s, %d chars billed)\n",
//			i, r.Text, r.DetectedSourceLang, r.BilledCharacters)
//	}
//
// Errors: Returns [ErrQuotaExceeded] if you've exhausted your monthly character
// quota, [ErrTooManyRequests] if rate-limited, [ErrBadRequest] for invalid
// parameters, or other API errors. Use errors.Is to check for specific failures.
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
