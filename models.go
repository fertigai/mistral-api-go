package mistral

import (
	"context"
	"fmt"
)

// ModelsService handles communication with the models related methods of the Mistral AI API.
type ModelsService struct {
	client *Client
}

// List returns a list of available models.
func (s *ModelsService) List(ctx context.Context) (*ModelList, error) {
	req, err := s.client.newRequest(ctx, "GET", "models", nil)
	if err != nil {
		return nil, err
	}

	var response ModelList
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Get retrieves a specific model by ID.
func (s *ModelsService) Get(ctx context.Context, modelID string) (*Model, error) {
	path := fmt.Sprintf("models/%s", modelID)
	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response Model
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetEnhanced retrieves detailed information about a model.
func (s *ModelsService) GetEnhanced(ctx context.Context, modelID string) (*EnhancedModel, error) {
	path := fmt.Sprintf("models/%s/enhanced", modelID)
	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response EnhancedModel
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// ListVersions returns all available versions of a model.
func (s *ModelsService) ListVersions(ctx context.Context, modelID string) ([]EnhancedModel, error) {
	path := fmt.Sprintf("models/%s/versions", modelID)
	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Versions []EnhancedModel `json:"versions"`
	}
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return response.Versions, nil
}

// GetCapabilities returns the capabilities of a model.
func (s *ModelsService) GetCapabilities(ctx context.Context, modelID string) ([]ModelCapability, error) {
	path := fmt.Sprintf("models/%s/capabilities", modelID)
	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Capabilities []ModelCapability `json:"capabilities"`
	}
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return response.Capabilities, nil
}

// GetPerformance returns the performance metrics of a model.
func (s *ModelsService) GetPerformance(ctx context.Context, modelID string) ([]ModelPerformance, error) {
	path := fmt.Sprintf("models/%s/performance", modelID)
	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Metrics []ModelPerformance `json:"metrics"`
	}
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return response.Metrics, nil
}

// EstimateTokens estimates the number of tokens in a text for a specific model.
func (s *ModelsService) EstimateTokens(ctx context.Context, modelID string, text string) (int, error) {
	path := fmt.Sprintf("models/%s/tokenize", modelID)
	req, err := s.client.newRequest(ctx, "POST", path, map[string]string{
		"text": text,
	})
	if err != nil {
		return 0, err
	}

	var response struct {
		TokenCount int `json:"token_count"`
	}
	_, err = s.client.do(req, &response)
	if err != nil {
		return 0, err
	}

	return response.TokenCount, nil
}
