package mistral

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("test-key")
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "test-key", client.apiKey)
	assert.Equal(t, defaultBaseURL, client.baseURL.String())
}

func TestNewClientWithOptions(t *testing.T) {
	customURL := "https://custom.mistral.ai"
	customClient := &http.Client{}

	client, err := NewClient(
		"test-key",
		WithBaseURL(customURL),
		WithHTTPClient(customClient),
	)

	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, customURL, client.baseURL.String())
	assert.Equal(t, customClient, client.client)
}

func TestNewClientWithInvalidURL(t *testing.T) {
	_, err := NewClient(
		"test-key",
		WithBaseURL("://invalid-url"),
	)
	assert.Error(t, err)
}

func TestNewClientWithNilHTTPClient(t *testing.T) {
	_, err := NewClient(
		"test-key",
		WithHTTPClient(nil),
	)
	assert.Error(t, err)
}

func TestNewRequest(t *testing.T) {
	client, err := NewClient("test-key")
	assert.NoError(t, err)

	ctx := context.Background()
	req, err := client.newRequest(ctx, "GET", "test", nil)
	assert.NoError(t, err)
	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "/v1/test", req.URL.Path)
	assert.Equal(t, "Bearer test-key", req.Header.Get("Authorization"))
}

func TestDo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"test": "response"}`))
	}))
	defer server.Close()

	client, err := NewClient(
		"test-key",
		WithBaseURL(server.URL),
	)
	assert.NoError(t, err)

	req, err := client.newRequest(context.Background(), "GET", "test", nil)
	assert.NoError(t, err)

	var response map[string]string
	_, err = client.do(req, &response)
	assert.NoError(t, err)
	assert.Equal(t, "response", response["test"])
}

func TestDoError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "error message"}`))
	}))
	defer server.Close()

	client, err := NewClient(
		"test-key",
		WithBaseURL(server.URL),
	)
	assert.NoError(t, err)

	req, err := client.newRequest(context.Background(), "GET", "test", nil)
	assert.NoError(t, err)

	var response map[string]string
	_, err = client.do(req, &response)
	assert.Error(t, err)

	apiErr, ok := err.(*Error)
	assert.True(t, ok)
	assert.Equal(t, "error message", apiErr.Message)
}
