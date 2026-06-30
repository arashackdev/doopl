// Package deepl is an idiomatic Go client for the DeepL API
// (https://developers.deepl.com). It covers text translation, document
// translation, glossaries, usage, languages, and text rephrasing (Write
// API), built against the current v3 API surface.
//
// # Basic Text Translation
//
//	client, err := deepl.New(os.Getenv("DEEPL_AUTH_KEY"))
//	if err != nil {
//		log.Fatal(err)
//	}
//	results, err := client.TranslateText(ctx, []string{"Hello, world!"}, "DE")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(results[0].Text) // "Hallo, Welt!"
//
// # Document Translation
//
//	file, err := os.Open("document.pdf")
//	defer file.Close()
//	output, err := os.Create("document_de.pdf")
//	defer output.Close()
//	err = client.TranslateDocument(ctx, file, "document.pdf", "DE", output)
//
// # Glossaries
//
//	entries := model.GlossaryEntries{"API": "Schnittstelle"}
//	glos, err := client.CreateGlossary(ctx, "tech", "EN", "DE", entries)
//	results, err := client.TranslateText(ctx, []string{"API call"}, "DE",
//		deepl.WithGlossaryID(glos.GlossaryID))
//
// # Languages & Usage
//
//	langs, err := client.SourceLanguages(ctx)
//	usage, err := client.Usage(ctx)
//
// # Write API (Rephrase)
//
//	results, err := client.Rephrase(ctx, []string{"Hello!"}, "EN",
//		deepl.WithTone(model.WriteToneFormal))
//
// # Configuration
//
// Free vs. Pro API endpoints are detected automatically from the API key
// (free keys end in ":fx"). Override with WithServerURL for custom endpoints
// or testing against deepl-mock.
//
// The Client is safe for concurrent use and embeds a standard *http.Client,
// so it composes cleanly with custom transports, proxies, and test doubles.
//
// # Error Handling
//
// Errors are returned as *deepl.Error, supporting errors.Is and errors.As:
//
//	var apiErr *deepl.Error
//	if errors.As(err, &apiErr) && errors.Is(err, deepl.ErrQuotaExceeded) {
//		// handle quota limit
//	}
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

// Client is the entry point for all DeepL API calls.
type Client struct {
	transport *transport.Transport
}

// Option configures a Client at construction time.
type Option func(*clientConfig)

type clientConfig struct {
	serverURL        string
	httpClient       *http.Client
	sendPlatformInfo bool
	appInfo          string // "appName/appVersion", appended to the User-Agent
	maxRetries       int
}

// WithServerURL overrides the DeepL API base URL. Use this for testing
// against deepl-mock, or for DeepL's regional endpoints (e.g. api-us.deepl.com).
// Overrides the automatic free/pro endpoint detection.
func WithServerURL(url string) Option {
	return func(c *clientConfig) { c.serverURL = url }
}

// WithHTTPClient sets the underlying *http.Client used for all requests.
// Use this to configure proxies, custom TLS config, or shared connection
// pooling across multiple SDK clients.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *clientConfig) { c.httpClient = hc }
}

// WithSendPlatformInfo controls whether Go/OS version info is included in
// the User-Agent sent to DeepL. Defaults to true. Disable for
// privacy-sensitive deployments.
func WithSendPlatformInfo(send bool) Option {
	return func(c *clientConfig) { c.sendPlatformInfo = send }
}

// WithAppInfo identifies your application in the User-Agent header, e.g.
// WithAppInfo("magus", "0.4.0"). DeepL recommends this for all integrations.
func WithAppInfo(appName, appVersion string) Option {
	return func(c *clientConfig) { c.appInfo = appName + "/" + appVersion }
}

// WithMaxRetries sets how many times a request is retried on 429/5xx
// responses before giving up. Defaults to 5.
func WithMaxRetries(n int) Option {
	return func(c *clientConfig) { c.maxRetries = n }
}

// New creates a Client authenticated with authKey. The free-vs-pro API
// endpoint is detected automatically from the key (free keys end in ":fx"),
// matching the behavior of DeepL's other official SDKs. Use WithServerURL
// to override this.
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
