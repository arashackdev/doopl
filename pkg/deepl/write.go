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

// RephraseOption configures a rephrase request.
type RephraseOption func(*rephraseParams)

type rephraseParams struct {
	tone  string
	emoji string
}

// WithTone sets the desired tone for rephrasing (formal or informal).
func WithTone(tone model.WriteTone) RephraseOption {
	return func(p *rephraseParams) { p.tone = string(tone) }
}

// WithEmojiMode controls emoji usage in the rephrased text.
func WithEmojiMode(mode model.WriteEmojiMode) RephraseOption {
	return func(p *rephraseParams) { p.emoji = string(mode) }
}

// Rephrase rephrases one or more texts in the given language.
// The Write API is only available for certain languages and has additional
// constraints compared to translation. Returns rephrased results in order.
func (c *Client) Rephrase(ctx context.Context, texts []string, lang string, opts ...RephraseOption) ([]model.WriteResult, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("deepl: texts must not be empty")
	}
	if lang == "" {
		return nil, fmt.Errorf("deepl: lang must not be empty")
	}

	p := &rephraseParams{}
	for _, opt := range opts {
		opt(p)
	}

	reqBody := apimodel.RephraseRequest{
		Text:  texts,
		Lang:  lang,
		Tone:  p.tone,
		Emoji: p.emoji,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("deepl: encoding request: %w", err)
	}

	resp, err := c.transport.Do(ctx, transportRequest("POST", "/v2/write", payload, "application/json"))
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

	var result apimodel.RephraseResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("deepl: decoding response: %w", err)
	}

	var writeResults []model.WriteResult
	for _, r := range result.Results {
		writeResults = append(writeResults, model.WriteResult{Text: r.Text})
	}
	return writeResults, nil
}
