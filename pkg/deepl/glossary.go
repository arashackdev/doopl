package deepl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/arashackdev/doopl/internal/apimodel"
	"github.com/arashackdev/doopl/pkg/model"
)

// CreateGlossary creates a new glossary with the given name and entries.
// Returns the created glossary with its new ID.
func (c *Client) CreateGlossary(ctx context.Context, name, sourceLang, targetLang string, entries model.GlossaryEntries) (*model.Glossary, error) {
	if name == "" {
		return nil, fmt.Errorf("deepl: name must not be empty")
	}
	if sourceLang == "" || targetLang == "" {
		return nil, fmt.Errorf("deepl: sourceLang and targetLang must not be empty")
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("deepl: entries must not be empty")
	}

	// Format entries as tab-separated TSV
	var lines []string
	for src, tgt := range entries {
		lines = append(lines, src+"\t"+tgt)
	}
	entriesData := strings.Join(lines, "\n")

	// Build form data
	form := url.Values{}
	form.Add("name", name)
	form.Add("source_lang", sourceLang)
	form.Add("target_lang", targetLang)
	form.Add("entries", entriesData)
	form.Add("entries_format", "tsv")

	req := transportRequest("POST", "/v2/glossaries", []byte(form.Encode()), "application/x-www-form-urlencoded")
	resp, err := c.transport.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("deepl: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("deepl: reading response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, newAPIError(resp.StatusCode, extractMessage(body))
	}

	var result apimodel.GlossaryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("deepl: decoding response: %w", err)
	}

	glossary := apiToModel.Glossary(result)
	return &glossary, nil
}

// ListGlossaries returns all glossaries for the account.
func (c *Client) ListGlossaries(ctx context.Context) ([]*model.Glossary, error) {
	resp, err := c.transport.Do(ctx, transportRequest("GET", "/v2/glossaries", nil, ""))
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

	var result apimodel.GlossariesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("deepl: decoding response: %w", err)
	}

	glossaries := apiToModel.Glossaries(result.Glossaries)
	var result2 []*model.Glossary
	for i := range glossaries {
		result2 = append(result2, &glossaries[i])
	}
	return result2, nil
}

// GetGlossary retrieves a specific glossary by ID.
func (c *Client) GetGlossary(ctx context.Context, glossaryID string) (*model.Glossary, error) {
	if glossaryID == "" {
		return nil, fmt.Errorf("deepl: glossaryID must not be empty")
	}

	resp, err := c.transport.Do(ctx, transportRequest("GET", "/v2/glossaries/"+glossaryID, nil, ""))
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

	var result apimodel.GlossaryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("deepl: decoding response: %w", err)
	}

	glossary := apiToModel.Glossary(result)
	return &glossary, nil
}

// DeleteGlossary deletes a glossary by ID.
func (c *Client) DeleteGlossary(ctx context.Context, glossaryID string) error {
	if glossaryID == "" {
		return fmt.Errorf("deepl: glossaryID must not be empty")
	}

	resp, err := c.transport.Do(ctx, transportRequest("DELETE", "/v2/glossaries/"+glossaryID, nil, ""))
	if err != nil {
		return fmt.Errorf("deepl: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return newAPIError(resp.StatusCode, extractMessage(body))
	}
	return nil
}

// GlossaryEntries retrieves all entries from a glossary.
func (c *Client) GlossaryEntries(ctx context.Context, glossaryID string) (model.GlossaryEntries, error) {
	if glossaryID == "" {
		return nil, fmt.Errorf("deepl: glossaryID must not be empty")
	}

	resp, err := c.transport.Do(ctx, transportRequest("GET", "/v2/glossaries/"+glossaryID+"/entries", nil, ""))
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

	var result apimodel.GlossaryEntriesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("deepl: decoding response: %w", err)
	}

	return model.GlossaryEntries(result.Entries), nil
}
