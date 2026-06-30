package deepl

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arashackdev/doopl/pkg/model"
)

func TestDocumentUpload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/document" && r.Method == "POST" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"document_id": "doc123", "document_key": "key456"}`))
			return
		}
	}))
	defer server.Close()

	client, _ := New("test-key", WithServerURL(server.URL))
	handle, err := client.DocumentUpload(context.Background(), bytes.NewReader([]byte("test file content")), "test.txt", "DE")
	require.NoError(t, err)
	assert.Equal(t, "doc123", handle.DocumentID)
	assert.Equal(t, "key456", handle.DocumentKey)
}

func TestDocumentStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/document/status" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{
				"document_id": "doc123",
				"status": "done",
				"billed_characters": 150,
				"character_count": 150
			}`))
			return
		}
	}))
	defer server.Close()

	client, _ := New("test-key", WithServerURL(server.URL))
	handle := &model.DocumentHandle{DocumentID: "doc123", DocumentKey: "key456"}
	status, err := client.DocumentStatus(context.Background(), handle)
	require.NoError(t, err)
	assert.Equal(t, model.DocumentStatusDone, status.Status)
	assert.True(t, status.Done())
}

func TestDocumentStatusError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/document/status" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			msg := "File format not supported"
			w.Write([]byte(`{
				"document_id": "doc123",
				"status": "error",
				"error_message": "` + msg + `"
			}`))
			return
		}
	}))
	defer server.Close()

	client, _ := New("test-key", WithServerURL(server.URL))
	handle := &model.DocumentHandle{DocumentID: "doc123", DocumentKey: "key456"}
	status, err := client.DocumentStatus(context.Background(), handle)
	require.NoError(t, err)
	assert.Equal(t, model.DocumentStatusError, status.Status)
	assert.True(t, status.Done())
}
