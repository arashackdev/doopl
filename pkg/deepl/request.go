package deepl

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/arashackdev/doopl/internal/transport"
	"github.com/arashackdev/doopl/v3/apimodel"
)

// transportRequest builds a transport.Request from a JSON (or pre-encoded)
// body. Pass a nil body for GET requests.
func transportRequest(method, path string, body []byte, contentType string) transport.Request {
	req := transport.Request{
		Method:      method,
		Path:        path,
		ContentType: contentType,
	}
	if body != nil {
		req.Body = bytes.NewReader(body)
	}
	return req
}

// extractMessage best-effort parses a DeepL error response body into a
// human-readable message. Falls back to the raw body if it isn't JSON.
func extractMessage(body []byte) string {
	var eb apimodel.ErrorResponse
	if err := json.Unmarshal(body, &eb); err == nil {
		if eb.Detail != "" {
			return eb.Detail
		}
		if eb.Message != "" {
			return eb.Message
		}
	}
	return string(body)
}

// drainAndClose is a small helper to fully consume and close a response body
// without leaking connections, used in code paths that don't need the body.
func drainAndClose(body io.ReadCloser) {
	if body == nil {
		return
	}
	_, _ = io.Copy(io.Discard, body)
	_ = body.Close()
}
