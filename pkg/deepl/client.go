// Package deepl is an idiomatic Go client for the DeepL API v3
// (https://developers.deepl.com). It provides full coverage of the DeepL
// translation platform, including text translation, document translation,
// glossaries, language information, usage tracking, and the Write API for
// text rephrasing. The client is built from the ground up for Go idioms:
// functional options, context support, proper error handling, and thread-safe
// concurrent operation.
//
// # Quick Start
//
// Create a client with your DeepL API key:
//
//	client, err := deepl.New(os.Getenv("DEEPL_AUTH_KEY"))
//	if err != nil {
//		log.Fatal(err)
//	}
//
// The free-vs-pro endpoint is detected automatically (free keys end in ":fx").
//
// # Text Translation
//
// Translate text into any supported target language:
//
//	results, err := client.TranslateText(ctx, []string{
//		"Hello, world!",
//		"How are you?",
//	}, "DE") // German
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(results[0].Text) // "Hallo, Welt!"
//
// Control translation with options:
//
//	results, err := client.TranslateText(ctx, texts, "FR",
//		deepl.WithSourceLang("EN"),        // Explicit source (usually auto-detected)
//		deepl.WithFormality(deepl.FormalityMore), // Formal French
//	)
//
// # Document Translation
//
// Upload a file for translation and download the result:
//
//	file, err := os.Open("report.pdf")
//	defer file.Close()
//
//	handle, err := client.DocumentUpload(ctx, file, "report.pdf", "ES")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer client.DocumentDelete(ctx, handle.DocumentID, handle.DocumentKey)
//
//	// Poll for completion
//	for {
//		status, err := client.DocumentStatus(ctx, handle.DocumentID, handle.DocumentKey)
//		if err != nil {
//			log.Fatal(err)
//		}
//		if status.Done {
//			break
//		}
//		time.Sleep(time.Second)
//	}
//
//	// Download the translated file
//	output, err := os.Create("report_es.pdf")
//	defer output.Close()
//	err = client.DocumentDownload(ctx, handle.DocumentID, handle.DocumentKey, output)
//
// # Glossaries
//
// Define custom term translations to enforce in all translations:
//
//	glossary, err := client.CreateGlossary(ctx, "tech-terms", "EN", "DE",
//		model.GlossaryEntries{
//			"API":      "Schnittstelle",
//			"backend":  "Datenbankseite",
//		},
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	results, err := client.TranslateText(ctx, texts, "DE",
//		deepl.WithSourceLang("EN"),
//		deepl.WithGlossaryID(glossary.GlossaryID),
//	)
//
// # Languages and Usage
//
// Check what languages are supported and track your API quota:
//
//	languages, err := client.SourceLanguages(ctx)
//	for _, lang := range languages {
//		fmt.Printf("%s: %s\n", lang.Code, lang.Name)
//	}
//
//	usage, err := client.Usage(ctx)
//	fmt.Printf("Used %d / %d characters\n", usage.CharacterCount, usage.CharacterLimit)
//
// # Write API (Rephrase)
//
// Rephrase text to match a specific tone or writing style:
//
//	results, err := client.Rephrase(ctx, []string{
//		"We are very pleased to inform you...",
//	}, "EN",
//		deepl.WithRephraseTone(model.WriteToneFormal),
//	)
//	fmt.Println(results[0].Text)
//
// # Configuration
//
// Functional options configure client behavior:
//
//	client, err := deepl.New(authKey,
//		deepl.WithAppInfo("myapp", "1.0.0"), // Sent in User-Agent
//		deepl.WithMaxRetries(3),              // Retry on 429/5xx
//		deepl.WithSendPlatformInfo(false),    // Privacy-sensitive deployments
//		deepl.WithServerURL("https://api.deepl.com"), // Override endpoint
//	)
//
// The Client is safe for concurrent use by multiple goroutines and composes
// cleanly with custom HTTP transports, proxies, and test doubles via
// WithHTTPClient.
//
// # Error Handling
//
// All methods return errors as *deepl.Error. Check for specific failures
// using errors.Is with sentinel errors:
//
//	results, err := client.TranslateText(ctx, texts, lang)
//	if err != nil {
//		var apiErr *deepl.Error
//		if errors.As(err, &apiErr) {
//			switch {
//			case errors.Is(err, deepl.ErrQuotaExceeded):
//				log.Println("Monthly character limit reached")
//			case errors.Is(err, deepl.ErrTooManyRequests):
//				log.Println("Rate limited; retry later")
//			case errors.Is(err, deepl.ErrAuthorizationFailed):
//				log.Println("Invalid or expired API key")
//			default:
//				log.Printf("DeepL error: %s", apiErr.Message)
//			}
//		}
//	}
//
// See the [documentation] for detailed API reference and additional examples.
//
// [documentation]: https://pkg.go.dev/github.com/arashackdev/doopl/pkg/deepl
package deepl

import (
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/arashackdev/doopl/internal/transport"
)

const (
	proBaseURL  = "https://api.deepl.com"
	freeBaseURL = "https://api-free.deepl.com"

	// Version is the SDK version, reported in the User-Agent header unless
	// platform info reporting is disabled. Bump this on release.
	Version = "0.1.0"
)

// Client is the entry point for all DeepL API v3 operations. It is built on
// top of a standard *http.Client and is safe for concurrent use by multiple
// goroutines. Create one with [New], configured via functional options.
type Client struct {
	transport *transport.Transport
}

// Option is a functional option that configures a Client at construction time.
// Common options include WithAppInfo, WithMaxRetries, WithServerURL, and
// WithHTTPClient. Pass multiple options to [New] to customize client behavior.
type Option func(*clientConfig)

type clientConfig struct {
	serverURL        string
	httpClient       *http.Client
	sendPlatformInfo bool
	appInfo          string // "appName/appVersion", appended to the User-Agent
	maxRetries       int
}

// WithServerURL overrides the DeepL API base URL. This disables the automatic
// free-vs-pro detection and is useful for:
//   - Testing against deepl-mock or other test doubles
//   - Using DeepL's regional endpoints (e.g., https://api-us.deepl.com)
//   - Custom or self-hosted deployments
//
// Example:
//
//	client, _ := deepl.New(authKey, deepl.WithServerURL("https://api-us.deepl.com"))
func WithServerURL(url string) Option {
	return func(c *clientConfig) { c.serverURL = url }
}

// WithHTTPClient sets the underlying *http.Client used for all requests.
// Use this to:
//   - Configure proxies or custom TLS settings
//   - Share connection pooling across multiple deepl.Client instances
//   - Inject custom transports for testing, observability, or authentication
//
// Example:
//
//	hc := &http.Client{
//		Transport: &http.Transport{MaxIdleConns: 100},
//		Timeout:   30 * time.Second,
//	}
//	client, _ := deepl.New(authKey, deepl.WithHTTPClient(hc))
func WithHTTPClient(hc *http.Client) Option {
	return func(c *clientConfig) { c.httpClient = hc }
}

// WithSendPlatformInfo controls whether Go runtime and OS information is
// included in the User-Agent header. Defaults to true. Set to false for
// privacy-sensitive deployments or when you want a minimal User-Agent string.
//
// Default User-Agent: "deepl-go/0.1.0 (go1.24; darwin/arm64)"
// With WithSendPlatformInfo(false): "deepl-go/0.1.0"
func WithSendPlatformInfo(send bool) Option {
	return func(c *clientConfig) { c.sendPlatformInfo = send }
}

// WithAppInfo identifies your application in the User-Agent header sent to
// DeepL. DeepL recommends this for all integrations so they can identify
// traffic patterns and provide support.
//
// Example User-Agent with WithAppInfo("myapp", "1.2.0"):
// "deepl-go/0.1.0 myapp/1.2.0 (go1.24; darwin/arm64)"
//
// Usage:
//
//	client, _ := deepl.New(authKey,
//		deepl.WithAppInfo("myapp", "1.2.0"),
//	)
func WithAppInfo(appName, appVersion string) Option {
	return func(c *clientConfig) { c.appInfo = appName + "/" + appVersion }
}

// WithMaxRetries sets how many times requests are retried on rate-limit (429)
// or server error (5xx) responses before returning an error. Defaults to 5.
// Each retry uses exponential backoff. Set to 0 to disable retries.
//
// Example: Disable retries for timeout-sensitive operations:
//
//	client, _ := deepl.New(authKey, deepl.WithMaxRetries(0))
func WithMaxRetries(n int) Option {
	return func(c *clientConfig) { c.maxRetries = n }
}

// New creates a Client authenticated with the provided API key. The free-vs-pro
// API endpoint is detected automatically from the key suffix (free keys end in
// ":fx"), matching the behavior of DeepL's other official SDKs.
//
// The returned Client is safe for concurrent use and ready for all API calls.
// A default 60-second request timeout is applied; customize with WithHTTPClient.
//
// Parameters:
//   - authKey: Your DeepL API key (required). Free keys end in ":fx".
//   - opts: Functional options like WithAppInfo, WithMaxRetries, etc.
//
// Returns an error if authKey is empty or any option fails.
//
// Example:
//
//	client, err := deepl.New(os.Getenv("DEEPL_AUTH_KEY"),
//		deepl.WithAppInfo("myapp", "1.0.0"),
//		deepl.WithMaxRetries(3),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
func New(authKey string, opts ...Option) (*Client, error) {
	if authKey == "" {
		return nil, &Error{StatusCode: 0, Message: "auth key must not be empty", sentinel: ErrAuthorizationFailed}
	}

	cfg := &clientConfig{
		sendPlatformInfo: true,
		httpClient:       &http.Client{Timeout: 60 * time.Second},
	}
	for _, opt := range opts {
		opt(cfg)
	}

	baseURL := cfg.serverURL
	if baseURL == "" {
		if strings.HasSuffix(authKey, ":fx") {
			baseURL = freeBaseURL
		} else {
			baseURL = proBaseURL
		}
	}

	return &Client{
		transport: &transport.Transport{
			HTTPClient: cfg.httpClient,
			BaseURL:    baseURL,
			AuthKey:    authKey,
			UserAgent:  buildUserAgent(cfg),
			MaxRetries: cfg.maxRetries,
		},
	}, nil
}

func buildUserAgent(cfg *clientConfig) string {
	ua := "deepl-go/" + Version
	if cfg.appInfo != "" {
		ua += " " + cfg.appInfo
	}
	if cfg.sendPlatformInfo {
		ua += " (" + runtime.Version() + "; " + runtime.GOOS + "/" + runtime.GOARCH + ")"
	}
	return ua
}
