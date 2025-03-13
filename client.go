package mistral

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://api.mistral.ai"
	userAgent      = "go-mistral"
	apiVersion     = "v1"
)

// Client manages communication with Mistral AI API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.
	baseURL *url.URL

	// User agent used when communicating with Mistral AI API.
	UserAgent string

	// API Key used for authentication
	apiKey string

	// Services used for communicating with different parts of the Mistral AI API.
	Chat        *ChatService
	Models      *ModelsService
	Embeddings  *EmbeddingsService
	Files       *FilesService
	FineTuning  *FineTuningService
	Moderations *ModerationsService
	OCR         *OCRService
	Agents      *AgentService
	FIM         *FIMService
	Classifiers *ClassifiersService
	Batch       *BatchService
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
		client:    http.DefaultClient,
		baseURL:   baseURL,
		UserAgent: userAgent,
		apiKey:    apiKey,
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	// Create services
	c.Chat = &ChatService{client: c}
	c.Models = &ModelsService{client: c}
	c.Embeddings = &EmbeddingsService{client: c}
	c.Files = &FilesService{client: c}
	c.FineTuning = &FineTuningService{client: c}
	c.Moderations = &ModerationsService{client: c}
	c.OCR = &OCRService{client: c}
	c.Agents = &AgentService{client: c}
	c.FIM = &FIMService{client: c}
	c.Classifiers = &ClassifiersService{client: c}
	c.Batch = &BatchService{client: c}

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

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) error {
		if httpClient == nil {
			return fmt.Errorf("HTTP client cannot be nil")
		}
		c.client = httpClient
		return nil
	}
}

// newRequest creates an API request.
func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(fmt.Sprintf("/%s/%s", apiVersion, path))
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	return req, nil
}

// do sends an API request and returns the API response.
func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = checkResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return resp, err
}

// Error represents an error response from the Mistral AI API.
type Error struct {
	Response *http.Response
	Message  string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		e.Response.Request.Method,
		e.Response.Request.URL,
		e.Response.StatusCode,
		e.Message)
}

// checkResponse checks the API response for errors.
func checkResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &Error{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.Message = string(data)
		}
	}

	return errorResponse
}
