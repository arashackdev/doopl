package deepl

import "fmt"

// Sentinel errors. Use errors.Is(err, deepl.ErrQuotaExceeded) etc. to check
// for specific failure conditions without parsing strings.
var (
	ErrAuthorizationFailed = fmt.Errorf("deepl: authorization failed (check your auth key)")
	ErrQuotaExceeded       = fmt.Errorf("deepl: quota exceeded")
	ErrTooManyRequests     = fmt.Errorf("deepl: too many requests (rate limited)")
	ErrResourceNotFound    = fmt.Errorf("deepl: resource not found")
	ErrBadRequest          = fmt.Errorf("deepl: bad request")
	ErrLimitExceeded       = fmt.Errorf("deepl: request size limit exceeded")
	ErrServiceUnavailable  = fmt.Errorf("deepl: service temporarily unavailable")
)

// Error is returned for any non-2xx response from the DeepL API. It wraps a
// sentinel error (see the Err* vars above) so callers can use errors.Is/As,
// while still retaining the raw status code and message body for logging.
type Error struct {
	StatusCode int
	Message    string
	sentinel   error
}

func (e *Error) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("deepl: %s (status %d)", e.Message, e.StatusCode)
	}
	return fmt.Sprintf("deepl: request failed with status %d", e.StatusCode)
}

// Unwrap allows errors.Is(err, deepl.ErrQuotaExceeded) to work.
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
