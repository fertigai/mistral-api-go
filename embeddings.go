package mistral

import (
	"context"
	"fmt"
	"math"
	"net/http"
)

// EmbeddingsService handles communication with the embeddings related methods of the Mistral AI API.
type EmbeddingsService struct {
	client *Client
}

// Create generates embeddings for the provided input text.
func (s *EmbeddingsService) Create(ctx context.Context, request *EmbeddingRequest) (*EmbeddingResponse, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "embeddings", request)
	if err != nil {
		return nil, err
	}

	var response EmbeddingResponse
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateEnhanced generates embeddings with additional options and metadata.
func (s *EmbeddingsService) CreateEnhanced(ctx context.Context, request *EnhancedEmbeddingRequest) (*EnhancedEmbeddingResponse, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "embeddings/enhanced", request)
	if err != nil {
		return nil, err
	}

	var response EnhancedEmbeddingResponse
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Similarity calculates the similarity between two texts using embeddings.
func (s *EmbeddingsService) Similarity(ctx context.Context, text1, text2 string, model string) (float32, error) {
	request := &EnhancedEmbeddingRequest{
		Model:     model,
		Input:     []string{text1, text2},
		Normalize: true,
	}

	resp, err := s.CreateEnhanced(ctx, request)
	if err != nil {
		return 0, err
	}

	if len(resp.Data) != 2 {
		return 0, fmt.Errorf("expected 2 embeddings, got %d", len(resp.Data))
	}

	// Calculate cosine similarity between the two embeddings
	vec1 := resp.Data[0].Embedding
	vec2 := resp.Data[1].Embedding

	var dotProduct float32
	var norm1 float32
	var norm2 float32

	for i := range vec1 {
		dotProduct += vec1[i] * vec2[i]
		norm1 += vec1[i] * vec1[i]
		norm2 += vec2[i] * vec2[i]
	}

	similarity := dotProduct / (float32(math.Sqrt(float64(norm1))) * float32(math.Sqrt(float64(norm2))))
	return similarity, nil
}

// Batch processes a large number of texts in batches to generate embeddings.
func (s *EmbeddingsService) Batch(ctx context.Context, texts []string, model string, batchSize int) ([][]float32, error) {
	if batchSize <= 0 {
		batchSize = 32 // Default batch size
	}

	var allEmbeddings [][]float32
	for i := 0; i < len(texts); i += batchSize {
		end := i + batchSize
		if end > len(texts) {
			end = len(texts)
		}

		batch := texts[i:end]
		request := &EmbeddingRequest{
			Model: model,
			Input: batch,
		}

		resp, err := s.Create(ctx, request)
		if err != nil {
			return nil, fmt.Errorf("error processing batch %d: %v", i/batchSize, err)
		}

		for _, data := range resp.Data {
			allEmbeddings = append(allEmbeddings, data.Embedding)
		}
	}

	return allEmbeddings, nil
}
