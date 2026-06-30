package deepl

import (
	"testing"
)

// TestTranslateText_MockHTTPClient demonstrates mock-based testing with mockery.
// This test is isolated from the real HTTP layer and verifies call behavior.
// Currently placeholders until Client is refactored to use the interfaces in internal/interfaces.
func TestTranslateText_MockHTTPClient(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping mock test in short mode")
	}
	// TODO: Once Client accepts HTTPClient interface, add mock-based test here.
	t.Skip("awaiting Client refactor to use interfaces")
}

// TestTranslateText_QuotaExceeded demonstrates error handling with mocks.
func TestTranslateText_QuotaExceeded(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping mock test in short mode")
	}
	// TODO: Once Client accepts HTTPClient interface, add mock-based error scenarios here.
	t.Skip("awaiting Client refactor to use interfaces")
}
