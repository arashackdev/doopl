package deepl

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLanguages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v3/languages" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{
				"languages": [
					{"language": "EN", "name": "English", "supports_formality": true},
					{"language": "DE", "name": "German", "supports_formality": true},
					{"language": "FR", "name": "French", "supports_formality": true}
				]
			}`))
			return
		}
	}))
	defer server.Close()

	client, _ := New("test-key", WithServerURL(server.URL))
	langs, err := client.SourceLanguages(context.Background())
	require.NoError(t, err)
	require.Len(t, langs, 3)
	assert.Equal(t, "EN", langs[0].Code)
}

func TestUsage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v3/usage" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{
				"character_count": 42000,
				"character_limit": 1000000,
				"document_count": 2,
				"document_limit": 50,
				"team_document_count": 5,
				"team_document_limit": 100
			}`))
			return
		}
	}))
	defer server.Close()

	client, _ := New("test-key", WithServerURL(server.URL))
	usage, err := client.Usage(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(42000), usage.CharacterCount)
	assert.Equal(t, int64(2), usage.DocumentCount)
}
