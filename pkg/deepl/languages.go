package deepl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/arashackdev/doopl/internal/apimodel"
	"github.com/arashackdev/doopl/pkg/model"
)

// Languages returns the list of supported languages for the given resource type.
// The resource parameter can be one of: "translate", "document", "glossary", or "write".
// Returns a slice of Language objects sorted by language code.
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

// SourceLanguages returns the list of supported source languages for translation.
// This is a convenience wrapper around Languages with resource="translate".
func (c *Client) SourceLanguages(ctx context.Context) ([]model.Language, error) {
	return c.Languages(ctx, "translate")
}

// TargetLanguages returns the list of supported target languages for translation.
// This is a convenience wrapper around Languages with resource="translate".
func (c *Client) TargetLanguages(ctx context.Context) ([]model.Language, error) {
	return c.Languages(ctx, "translate")
}

// Usage returns the current API quota and consumption for your account.
// Use this to check how many characters you have translated, how many
// documents are in progress, and what your limits are.
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
