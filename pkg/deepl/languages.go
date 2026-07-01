package deepl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/arashackdev/doopl/pkg/model"
	"github.com/arashackdev/doopl/v3/apimodel"
)

// Languages returns the list of supported languages for the specified resource
// type. This is a low-level API; most users should prefer SourceLanguages,
// TargetLanguages, or DocumentLanguages convenience methods instead.
//
// Parameters:
//   - ctx: Context controlling request lifetime.
//   - resource: The resource type. Supported values: "translate" (default),
//     "document", "glossary", or "write".
//
// Returns a slice of Language objects sorted by language code. Each Language
// includes a code (e.g., "EN"), name (e.g., "English"), and supports_formality
// flag indicating whether the language supports formality levels.
//
// Example:
//
//	// List all languages supported for document translation
//	langs, err := client.Languages(ctx, "document")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, lang := range langs {
//		fmt.Printf("  %s: %s\n", lang.Code, lang.Name)
//	}
func (c *Client) Languages(ctx context.Context, resource string) ([]model.Language, error) {
	resp, err := c.transport.Do(ctx, transportRequest("GET", "/v3/languages?type="+resource, nil, ""))
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

	var result apimodel.LanguagesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("deepl: decoding response: %w", err)
	}

	return apiToModel.Languages(result.Languages), nil
}

// SourceLanguages returns the list of languages that can be used as input to
// TranslateText. This includes a special "auto" language for auto-detection.
// Convenience wrapper for Languages(ctx, "translate").
//
// Use this to validate user input or display a menu of supported source languages.
//
// Example:
//
//	langs, _ := client.SourceLanguages(ctx)
//	fmt.Println("Supported sources:", langs)
func (c *Client) SourceLanguages(ctx context.Context) ([]model.Language, error) {
	return c.Languages(ctx, "translate")
}

// TargetLanguages returns the list of languages that can be translated into.
// Convenience wrapper for Languages(ctx, "translate").
//
// Use this to validate or list available translation targets.
//
// Example:
//
//	langs, _ := client.TargetLanguages(ctx)
//	for _, lang := range langs {
//		fmt.Printf("%s (%s)\n", lang.Code, lang.Name)
//	}
func (c *Client) TargetLanguages(ctx context.Context) ([]model.Language, error) {
	return c.Languages(ctx, "translate")
}

// Usage returns the current API quota and consumption details for your DeepL
// account. This is useful for quota-aware workflows, monitoring, and billing
// integration. The returned Usage object includes character count, limits, and
// document counts.
//
// Returns: A Usage struct containing:
//   - CharacterCount: Number of characters already translated this billing period.
//   - CharacterLimit: Your account's monthly character quota.
//   - DocumentCount: Number of documents currently in progress.
//   - DocumentLimit: Maximum concurrent documents allowed.
//   - And team-level quotas if applicable.
//
// Example: Check if you're approaching your quota:
//
//	usage, err := client.Usage(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	percentUsed := float64(usage.CharacterCount) / float64(usage.CharacterLimit) * 100
//	fmt.Printf("Quota usage: %.1f%% (%d / %d chars)\n",
//		percentUsed, usage.CharacterCount, usage.CharacterLimit)
func (c *Client) Usage(ctx context.Context) (*model.Usage, error) {
	resp, err := c.transport.Do(ctx, transportRequest("GET", "/v3/usage", nil, ""))
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

	var result apimodel.UsageResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("deepl: decoding response: %w", err)
	}

	usage := apiToModel.Usage(result)
	return &usage, nil
}
