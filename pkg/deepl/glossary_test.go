package deepl

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arashackdev/doopl/pkg/model"
)

func TestCreateGlossary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/glossaries" && r.Method == "POST" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(201)
			w.Write([]byte(`{
				"glossary_id": "glos-123",
				"name": "tech-terms",
				"ready": true,
				"source_lang": "EN",
				"target_lang": "DE",
				"creation_time": "2024-01-01T00:00:00Z",
				"entry_count": 10
			}`))
			return
		}
	}))
	defer server.Close()

	client, _ := New("test-key", WithServerURL(server.URL))
	entries := model.GlossaryEntries{
		"cloud":  "Wolke",
		"server": "Server",
	}
	glos, err := client.CreateGlossary(context.Background(), "tech-terms", "EN", "DE", entries)
	require.NoError(t, err)
	assert.Equal(t, "glos-123", glos.GlossaryID)
}

func TestListGlossaries(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/glossaries" && r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{
				"glossaries": [
					{"glossary_id": "glos-1", "name": "tech", "ready": true, "source_lang": "EN", "target_lang": "DE", "creation_time": "2024-01-01T00:00:00Z", "entry_count": 5},
					{"glossary_id": "glos-2", "name": "medical", "ready": true, "source_lang": "EN", "target_lang": "FR", "creation_time": "2024-01-02T00:00:00Z", "entry_count": 8}
				]
			}`))
			return
		}
	}))
	defer server.Close()

	client, _ := New("test-key", WithServerURL(server.URL))
	glos, err := client.ListGlossaries(context.Background())
	require.NoError(t, err)
	assert.Len(t, glos, 2)
}

func TestGlossaryEntries(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/glossaries/glos-123/entries" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{
				"entries": {
					"hello": "bonjour",
					"goodbye": "au revoir"
				}
			}`))
			return
		}
	}))
	defer server.Close()

	client, _ := New("test-key", WithServerURL(server.URL))
	entries, err := client.GlossaryEntries(context.Background(), "glos-123")
	require.NoError(t, err)
	assert.Equal(t, "bonjour", entries["hello"])
}
