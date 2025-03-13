package mistral

import (
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
