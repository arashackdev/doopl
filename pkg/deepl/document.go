package deepl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/arashackdev/doopl/pkg/model"
	"github.com/arashackdev/doopl/v3/apimodel"
)

// DocumentUploadOption configures a document upload request.
type DocumentUploadOption func(*documentUploadParams)

type documentUploadParams struct {
	sourceLang                 string
	glossaryID                 string
	styleid                    string
	translationMemoryID        string
	translationMemoryThreshold *float64
}

// WithDocumentSourceLang sets the source language for document translation.
func WithDocumentSourceLang(lang string) DocumentUploadOption {
	return func(p *documentUploadParams) { p.sourceLang = lang }
}

// WithDocumentGlossaryID applies a glossary to the document translation.
func WithDocumentGlossaryID(id string) DocumentUploadOption {
	return func(p *documentUploadParams) { p.glossaryID = id }
}

// WithDocumentStyleID applies a style rule to the document translation.
func WithDocumentStyleID(id string) DocumentUploadOption {
	return func(p *documentUploadParams) { p.styleid = id }
}

// WithDocumentTranslationMemoryID applies a translation memory to the document.
func WithDocumentTranslationMemoryID(id string) DocumentUploadOption {
	return func(p *documentUploadParams) { p.translationMemoryID = id }
}

// WithDocumentTranslationMemoryThreshold sets the minimum translation memory match score.
func WithDocumentTranslationMemoryThreshold(threshold float64) DocumentUploadOption {
	return func(p *documentUploadParams) { p.translationMemoryThreshold = &threshold }
}

// DocumentUpload uploads a file to DeepL for translation. Returns a DocumentHandle
// that can be used to check status and download results. The filename argument is
// used to set the Content-Disposition; the caller is responsible for inferring
// language from extension (.docx, .pptx, .pdf, .xlsx, .txt, .html, .htm, .jpg, .jpeg, .png).
func (c *Client) DocumentUpload(ctx context.Context, data io.Reader, filename, targetLang string, opts ...DocumentUploadOption) (*model.DocumentHandle, error) {
	if filename == "" {
		return nil, fmt.Errorf("deepl: filename must not be empty")
	}
	if targetLang == "" {
		return nil, fmt.Errorf("deepl: targetLang must not be empty")
	}

	p := &documentUploadParams{}
	for _, opt := range opts {
		opt(p)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file field
	part, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return nil, fmt.Errorf("deepl: creating form file: %w", err)
	}
	if _, err := io.Copy(part, data); err != nil {
		return nil, fmt.Errorf("deepl: reading file: %w", err)
	}

	// Add required fields
	if err := writer.WriteField("target_lang", targetLang); err != nil {
		return nil, fmt.Errorf("deepl: writing target_lang: %w", err)
	}

	// Add optional fields
	if p.sourceLang != "" {
		if err := writer.WriteField("source_lang", p.sourceLang); err != nil {
			return nil, err
		}
	}
	if p.glossaryID != "" {
		if err := writer.WriteField("glossary_id", p.glossaryID); err != nil {
			return nil, err
		}
	}
	if p.styleid != "" {
		if err := writer.WriteField("style_id", p.styleid); err != nil {
			return nil, err
		}
	}
	if p.translationMemoryID != "" {
		if err := writer.WriteField("translation_memory_id", p.translationMemoryID); err != nil {
			return nil, err
		}
	}
	if p.translationMemoryThreshold != nil {
		if err := writer.WriteField("translation_memory_threshold", fmt.Sprintf("%f", *p.translationMemoryThreshold)); err != nil {
			return nil, err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("deepl: closing multipart writer: %w", err)
	}

	req := transportRequest("POST", "/v2/document", body.Bytes(), writer.FormDataContentType())
	resp, err := c.transport.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("deepl: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("deepl: reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newAPIError(resp.StatusCode, extractMessage(respBody))
	}

	var result apimodel.DocumentUploadResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("deepl: decoding response: %w", err)
	}

	return &model.DocumentHandle{DocumentID: result.DocumentID, DocumentKey: result.DocumentKey}, nil
}

// DocumentStatus checks the translation status of an uploaded document.
func (c *Client) DocumentStatus(ctx context.Context, handle *model.DocumentHandle) (*model.DocumentStatusInfo, error) {
	if handle.DocumentID == "" || handle.DocumentKey == "" {
		return nil, fmt.Errorf("deepl: invalid document handle")
	}

	payload := map[string]string{
		"document_id":  handle.DocumentID,
		"document_key": handle.DocumentKey,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req := transportRequest("POST", "/v2/document/status", body, "application/json")
	resp, err := c.transport.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("deepl: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("deepl: reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newAPIError(resp.StatusCode, extractMessage(respBody))
	}

	var result apimodel.DocumentStatusResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("deepl: decoding response: %w", err)
	}

	// Convert status string to enum
	statusInfo := apiToModel.DocumentStatus(result)
	statusInfo.Status = model.DocumentStatus(result.Status)
	return &statusInfo, nil
}

// DocumentDownload retrieves the translated document. Call this after DocumentStatus
// reports status=="done". The caller is responsible for closing the returned
// io.ReadCloser.
func (c *Client) DocumentDownload(ctx context.Context, handle *model.DocumentHandle, w io.Writer) error {
	if handle.DocumentID == "" || handle.DocumentKey == "" {
		return fmt.Errorf("deepl: invalid document handle")
	}

	payload := map[string]string{
		"document_id":  handle.DocumentID,
		"document_key": handle.DocumentKey,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req := transportRequest("POST", "/v2/document/result", body, "application/json")
	resp, err := c.transport.Do(ctx, req)
	if err != nil {
		return fmt.Errorf("deepl: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return newAPIError(resp.StatusCode, extractMessage(respBody))
	}

	_, err = io.Copy(w, resp.Body)
	return err
}

// TranslateDocument is a convenience function that uploads a document,
// polls until translation is complete, and downloads the result — all in one call.
// It respects ctx cancellation and backs off between status checks.
func (c *Client) TranslateDocument(ctx context.Context, data io.Reader, filename, targetLang string, w io.Writer, opts ...DocumentUploadOption) error {
	handle, err := c.DocumentUpload(ctx, data, filename, targetLang, opts...)
	if err != nil {
		return err
	}

	// Poll with exponential backoff until done or ctx cancelled
	backoff := time.Second
	maxBackoff := 30 * time.Second
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		status, err := c.DocumentStatus(ctx, handle)
		if err != nil {
			return err
		}

		if status.Done() {
			if status.Status == model.DocumentStatusError {
				if status.ErrorMessage != nil {
					return fmt.Errorf("deepl: document translation error: %s", *status.ErrorMessage)
				}
				return fmt.Errorf("deepl: document translation failed")
			}
			return c.DocumentDownload(ctx, handle, w)
		}

		// Back off before next poll
		select {
		case <-time.After(backoff):
			if backoff < maxBackoff {
				backoff *= 2
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
