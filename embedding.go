package mistral

import (
	"fmt"
	"net/http"
)

type EmbeddingService struct {
	client *Client
}

type EmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type EmbeddingUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type EmbeddingData struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
}

type EmbeddingResponse struct {
	Id     string          `json:"id"`
	Object string          `json:"object"`
	Model  string          `json:"model"`
	Data   []EmbeddingData `json:"data"`
	Usage  EmbeddingUsage  `json:"usage"`
}

func (s *EmbeddingService) Create(request *EmbeddingRequest) (*EmbeddingResponse, error) {
	var response EmbeddingResponse
	resp, err := s.client.client.R().
		SetHeader("Authorization", s.client.AuthHeader()).
		SetBody(request).
		SetResult(&response).
		Post(s.client.FormUrl("embeddings"))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to create embedding: %s", resp.String())
	}

	return &response, nil
}
