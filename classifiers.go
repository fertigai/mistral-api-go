package mistral

import (
	"context"
	"fmt"
	"net/http"
)

// ClassifiersService handles communication with the classifier-related methods of the Mistral AI API.
type ClassifiersService struct {
	client *Client
}

// Create submits a classification request.
func (s *ClassifiersService) Create(ctx context.Context, request *ClassifierRequest) (*ClassifierResponse, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "classifiers", request)
	if err != nil {
		return nil, err
	}

	var response ClassifierResponse
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Batch submits multiple texts for classification in a single request.
func (s *ClassifiersService) Batch(ctx context.Context, texts []string, labels []string, model string) (*ClassifierResponse, error) {
	request := &ClassifierRequest{
		Model:  model,
		Input:  texts,
		Labels: labels,
	}

	return s.Create(ctx, request)
}

// MultiLabel submits a multi-label classification request.
func (s *ClassifiersService) MultiLabel(ctx context.Context, request *ClassifierRequest) (*ClassifierResponse, error) {
	request.MultiLabel = true
	return s.Create(ctx, request)
}

// Confidence returns the confidence scores for a single text against all labels.
func (s *ClassifiersService) Confidence(ctx context.Context, text string, labels []string, model string) (map[string]float32, error) {
	request := &ClassifierRequest{
		Model:  model,
		Input:  []string{text},
		Labels: labels,
	}

	resp, err := s.Create(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(resp.Results) == 0 {
		return nil, fmt.Errorf("no classification results returned")
	}

	return resp.Results[0].Scores, nil
}
