package deepl

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arashackdev/doopl/internal/apimodel"
)

func TestTranslateText_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/translate", r.URL.Path)
		assert.Equal(t, "DeepL-Auth-Key test-key", r.Header.Get("Authorization"))

		var body apimodel.TranslateRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "DE", body.TargetLang)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(apimodel.TranslateResponse{
			Translations: []apimodel.Translation{
				{Text: "Hallo, Welt!", DetectedSourceLang: "EN", BilledCharacters: 13},
			},
		})
	}))
	defer srv.Close()

	client, err := New("test-key", WithServerURL(srv.URL), WithSendPlatformInfo(false))
	require.NoError(t, err)

	results, err := client.TranslateText(context.Background(), []string{"Hello, world!"}, "DE")
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "Hallo, Welt!", results[0].Text)
	assert.Equal(t, "EN", results[0].DetectedSourceLang)
}

func TestTranslateText_AuthError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "Wrong endpoint. Use api-free.deepl.com"})
	}))
	defer srv.Close()

	client, err := New("test-key:fx", WithServerURL(srv.URL))
	require.NoError(t, err)

	_, err = client.TranslateText(context.Background(), []string{"Hello"}, "DE")
	require.Error(t, err)
	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusForbidden, apiErr.StatusCode)
}

func TestEndpointAutoDetection(t *testing.T) {
	cases := []struct {
		key      string
		wantFree bool
	}{
		{"abc123:fx", true},
		{"abc123", false},
	}
	for _, tc := range cases {
		c, err := New(tc.key)
		require.NoError(t, err, "New(%q) error", tc.key)
		gotFree := c.transport.BaseURL == freeBaseURL
		assert.Equal(t, tc.wantFree, gotFree, "New(%q): baseURL = %s, wantFree = %v", tc.key, c.transport.BaseURL, tc.wantFree)
	}
}
