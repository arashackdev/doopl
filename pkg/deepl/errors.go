package deepl

import "fmt"

// Sentinel errors for use with errors.Is. Each maps to a specific HTTP status
// code or failure condition. Use errors.Is(err, deepl.ErrQuotaExceeded) etc. to
// write robust error handling without string parsing. All API errors will
// unwrap to one of these sentinels.
//
// Example error handling:
//
//	results, err := client.TranslateText(ctx, texts, lang)
//	if err != nil {
//		if errors.Is(err, deepl.ErrQuotaExceeded) {
//			// Handle quota exhaustion
//		} else if errors.Is(err, deepl.ErrTooManyRequests) {
//			// Handle rate limiting (retry with backoff)
//		} else if errors.Is(err, deepl.ErrAuthorizationFailed) {
//			// Handle auth failure (check API key)
//		}
//	}
var (
	ErrAuthorizationFailed = fmt.Errorf("deepl: authorization failed (check your auth key)")
	ErrQuotaExceeded       = fmt.Errorf("deepl: quota exceeded (monthly character limit reached)")
	ErrTooManyRequests     = fmt.Errorf("deepl: too many requests (rate limited; retry with backoff)")
	ErrResourceNotFound    = fmt.Errorf("deepl: resource not found (invalid ID or document finished)")
	ErrBadRequest          = fmt.Errorf("deepl: bad request (invalid parameters)")
	ErrLimitExceeded       = fmt.Errorf("deepl: request size limit exceeded (text or file too large)")
	ErrServiceUnavailable  = fmt.Errorf("deepl: service temporarily unavailable (retry later)")
)

// Error represents an API error response from DeepL. It includes the HTTP status
// code, the error message from DeepL's response body, and wraps a sentinel error
// (see Err* vars) to enable robust error handling via errors.Is and errors.As.
//
// Access fields for detailed logging or alternative error messages:
//
//	var apiErr *deepl.Error
//	if errors.As(err, &apiErr) {
//		log.Printf("API Error %d: %s", apiErr.StatusCode, apiErr.Message)
//	}
type Error struct {
	StatusCode int    // HTTP status code (e.g., 429, 456, 503)
	Message    string // Error message from DeepL's response body (if any)
	sentinel   error  // Wrapped sentinel error (e.g., ErrQuotaExceeded)
}

func (e *Error) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("deepl: %s (status %d)", e.Message, e.StatusCode)
	}
	return fmt.Sprintf("deepl: request failed with status %d", e.StatusCode)
}

// Unwrap returns the wrapped sentinel error, enabling errors.Is(err, deepl.ErrQuotaExceeded).
func (e *Error) Unwrap() error { return e.sentinel }

// newAPIError maps an HTTP status code to the appropriate sentinel error and
// wraps it alongside the response body's message, if any.
func newAPIError(statusCode int, message string) error {
	var sentinel error
	switch statusCode {
	case 401, 403:
		sentinel = ErrAuthorizationFailed
	case 404:
		sentinel = ErrResourceNotFound
	case 413:
		sentinel = ErrLimitExceeded
	case 429:
		sentinel = ErrTooManyRequests
	case 456:
		sentinel = ErrQuotaExceeded
	case 503:
		sentinel = ErrServiceUnavailable
	case 400:
		sentinel = ErrBadRequest
	default:
		sentinel = fmt.Errorf("deepl: unexpected status %d", statusCode)
	}
	return &Error{StatusCode: statusCode, Message: message, sentinel: sentinel}
}
