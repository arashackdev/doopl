package transport

// Dependencies used by this package:
//
// - net/http: Standard library HTTP client and request/response types.
//   Used by: Transport.Do() - executes HTTP requests with custom auth/retry logic
//   Key types: *http.Client, *http.Request, *http.Response
//
// Provides:
//
// - Transport: Wraps *http.Client with DeepL-specific features:
//   * Auth header injection (Authorization: Bearer <key>)
//   * Retry with exponential backoff on 429/5xx
//   * Retry-After header support
//   * User-Agent construction
//   Safe for concurrent use (same contract as *http.Client)
//
// See: transport.go
