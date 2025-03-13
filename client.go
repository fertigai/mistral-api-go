package mistral

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
)

const (
	defaultBaseURL = "https://api.mistral.ai"
	userAgent      = "go-mistral"
	apiVersion     = "v1"
)

// Client manages communication with Mistral AI API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *resty.Client

	// Base URL for API requests.
	baseURL *url.URL

	// User agent used when communicating with Mistral AI API.
	UserAgent string

	// API Key used for authentication
	apiKey string

	// Services used for communicating with different parts of the Mistral AI API.
	OCR *OCRService

	Files *FilesService
}

func (c *Client) AuthHeader() string {
	return fmt.Sprintf("Bearer %s", c.apiKey)
}

func (c *Client) FormUrl(path string) string {
	return fmt.Sprintf("%s/%s/%s", c.baseURL.String(), apiVersion, path)
}

// ClientOption is a function that modifies the client.
type ClientOption func(*Client) error

// NewClient creates a new Mistral AI API client.
func NewClient(apiKey string, opts ...ClientOption) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		client:    resty.New(),
		baseURL:   baseURL,
		UserAgent: userAgent,
		apiKey:    apiKey,
	}

	// Apply any custom options
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	// Initialize services
	c.OCR = &OCRService{client: c}
	c.Files = &FilesService{client: c}

	return c, nil
}

// WithBaseURL sets a custom base URL for the client.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		u, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.baseURL = u
		return nil
	}
}

// Error represents an error response from the API.
type Error struct {
	Response *http.Response
	Message  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		e.Response.Request.Method,
		e.Response.Request.URL,
		e.Response.StatusCode,
		e.Message,
	)
}
