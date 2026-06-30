// Package transport implements the low-level HTTP plumbing shared by every
// DeepL API call: auth header injection, retry with exponential backoff +
// jitter (honoring Retry-After), and User-Agent construction. It is
// intentionally unexported from the public deepl package — callers never see
// these types directly.
package transport

import (
	"bytes"
	"context"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	defaultMaxRetries = 5
	defaultBaseDelay  = 1 * time.Second
	defaultMaxDelay   = 60 * time.Second
)

// Transport wraps an *http.Client with DeepL-specific request behavior. It is
// safe for concurrent use, matching the underlying http.Client's contract.
type Transport struct {
	HTTPClient *http.Client
	BaseURL    string
	AuthKey    string
	UserAgent  string
	MaxRetries int
}

// Request describes a single API call before retries/auth are applied.
type Request struct {
	Method      string
	Path        string // joined with BaseURL, e.g. "/v2/translate"
	Query       url.Values
	Body        io.Reader
	ContentType string
}

// Do executes req against the DeepL API, retrying on 429/5xx with
// exponential backoff and jitter, and honoring the Retry-After header when
// present. It does not interpret the response body — callers decode the
// JSON themselves and use IsRetryable/status mapping as needed.
func (t *Transport) Do(ctx context.Context, req Request) (*http.Response, error) {
	maxRetries := t.MaxRetries
	if maxRetries <= 0 {
		maxRetries = defaultMaxRetries
	}

	fullURL := t.BaseURL + req.Path
	if len(req.Query) > 0 {
		fullURL += "?" + req.Query.Encode()
	}

	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
	}

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := backoffDelay(attempt)
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		var bodyReader io.Reader
		if bodyBytes != nil {
			bodyReader = bytes.NewReader(bodyBytes)
		}

		httpReq, err := http.NewRequestWithContext(ctx, req.Method, fullURL, bodyReader)
		if err != nil {
			return nil, err
		}
		httpReq.Header.Set("Authorization", "DeepL-Auth-Key "+t.AuthKey)
		httpReq.Header.Set("User-Agent", t.UserAgent)
		if req.ContentType != "" {
			httpReq.Header.Set("Content-Type", req.ContentType)
		}

		resp, err := t.HTTPClient.Do(httpReq)
		if err != nil {
			lastErr = err
			continue // network error: retry
		}

		if !isRetryableStatus(resp.StatusCode) || attempt == maxRetries {
			return resp, nil
		}

		// Retryable status: honor Retry-After if present, drain+close body,
		// then loop for another attempt.
		retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))
		resp.Body.Close()
		if retryAfter > 0 {
			select {
			case <-time.After(retryAfter):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}

	return nil, lastErr
}

func isRetryableStatus(code int) bool {
	return code == http.StatusTooManyRequests || (code >= 500 && code < 600)
}

func backoffDelay(attempt int) time.Duration {
	exp := float64(defaultBaseDelay) * math.Pow(2, float64(attempt-1))
	jitter := rand.Float64() * exp * 0.25 //nolint:gosec // jitter doesn't need crypto/rand
	d := time.Duration(exp + jitter)
	if d > defaultMaxDelay {
		return defaultMaxDelay
	}
	return d
}

func parseRetryAfter(v string) time.Duration {
	if v == "" {
		return 0
	}
	if secs, err := strconv.Atoi(v); err == nil {
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(v); err == nil {
		return time.Until(t)
	}
	return 0
}
