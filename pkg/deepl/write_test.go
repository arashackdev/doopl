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

func TestRephrase(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/write" && r.Method == "POST" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{
				"results": [
					{"text": "Greetings, esteemed colleague."}
				]
			}`))
			return
		}
	}))
	defer server.Close()

	client, _ := New("test-key", WithServerURL(server.URL))
	results, err := client.Rephrase(context.Background(), []string{"Hello, friend."}, "EN", WithTone(model.WriteToneFormal))
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "Greetings, esteemed colleague.", results[0].Text)
}
