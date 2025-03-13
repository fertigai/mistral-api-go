package mistral

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// BatchService handles communication with the batch processing methods of the Mistral AI API.
type BatchService struct {
	client *Client
}

// Create submits a batch processing request.
func (s *BatchService) Create(ctx context.Context, request *BatchRequest) (*BatchResponse, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "batch", request)
	if err != nil {
		return nil, err
	}

	var response BatchResponse
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Get retrieves the status and results of a batch operation.
func (s *BatchService) Get(ctx context.Context, batchID string) (*BatchResponse, error) {
	path := fmt.Sprintf("batch/%s", batchID)
	req, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var response BatchResponse
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Cancel cancels a running batch operation.
func (s *BatchService) Cancel(ctx context.Context, batchID string) error {
	path := fmt.Sprintf("batch/%s/cancel", batchID)
	req, err := s.client.newRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(req, nil)
	return err
}

// CreateChat submits a batch of chat completion requests.
func (s *BatchService) CreateChat(ctx context.Context, requests []*ChatCompletionRequest, options *BatchOptions) (*BatchResponse, error) {
	if options == nil {
		options = &BatchOptions{
			MaxConcurrency: 5,
			ChunkSize:      100,
		}
	}

	batchReq := &BatchRequest{
		Requests: make([]interface{}, len(requests)),
		Model:    requests[0].Model, // Use the model from the first request
		Options:  *options,
	}

	for i, req := range requests {
		batchReq.Requests[i] = req
	}

	return s.Create(ctx, batchReq)
}

// CreateEmbeddings submits a batch of embedding requests.
func (s *BatchService) CreateEmbeddings(ctx context.Context, texts []string, model string, options *BatchOptions) (*BatchResponse, error) {
	if options == nil {
		options = &BatchOptions{
			MaxConcurrency: 5,
			ChunkSize:      1000,
		}
	}

	requests := make([]*EmbeddingRequest, len(texts))
	for i, text := range texts {
		requests[i] = &EmbeddingRequest{
			Model: model,
			Input: []string{text},
		}
	}

	batchReq := &BatchRequest{
		Requests: make([]interface{}, len(requests)),
		Model:    model,
		Options:  *options,
	}

	for i, req := range requests {
		batchReq.Requests[i] = req
	}

	return s.Create(ctx, batchReq)
}

// WaitForCompletion waits for a batch operation to complete.
func (s *BatchService) WaitForCompletion(ctx context.Context, batchID string) (*BatchResponse, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			resp, err := s.Get(ctx, batchID)
			if err != nil {
				return nil, err
			}

			if resp.Summary.Failed+resp.Summary.Succeeded == resp.Summary.TotalRequests {
				return resp, nil
			}

			// Add a small delay before the next check
			time.Sleep(time.Second * 2)
		}
	}
}
