// Package interfaces defines mockable interfaces for the doopl client.
package interfaces

import (
	"context"
	"io"
	"net/http"

	"github.com/arashackdev/doopl/internal/apimodel"
)

// HTTPClient is a mockable interface for HTTP operations.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Translator handles text translation requests.
type Translator interface {
	Translate(ctx context.Context, req *apimodel.TranslateRequest) (*apimodel.TranslateResponse, error)
}

// DocumentHandler manages document translation operations.
type DocumentHandler interface {
	Upload(ctx context.Context, file io.Reader, filename string, lang string) (*apimodel.DocumentUploadResponse, error)
	Status(ctx context.Context, docID string, apiKey string) (*apimodel.DocumentStatusResponse, error)
	Download(ctx context.Context, docID string, apiKey string) (io.ReadCloser, error)
}

// GlossaryManager handles glossary CRUD operations.
type GlossaryManager interface {
	Create(ctx context.Context, name, sourceLang, targetLang string, entries map[string]string) (*apimodel.GlossaryResponse, error)
	List(ctx context.Context) (*apimodel.GlossariesResponse, error)
	Entries(ctx context.Context, glossaryID string) (*apimodel.GlossaryEntriesResponse, error)
	Delete(ctx context.Context, glossaryID string) error
}

// LanguagesProvider handles language and usage queries.
type LanguagesProvider interface {
	SourceLanguages(ctx context.Context) ([]apimodel.Language, error)
	TargetLanguages(ctx context.Context) ([]apimodel.Language, error)
	Usage(ctx context.Context) (*apimodel.UsageResponse, error)
}

// Writer handles text rephrasing (Write API).
type Writer interface {
	Rephrase(ctx context.Context, req *apimodel.RephraseRequest) (*apimodel.RephraseResponse, error)
}
