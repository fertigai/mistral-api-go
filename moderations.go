package mistral

import (
	"context"
	"net/http"
)

// ModerationsService handles communication with the moderations related methods of the Mistral AI API.
type ModerationsService struct {
	client *Client
}

// Create submits a moderation request.
func (s *ModerationsService) Create(ctx context.Context, request *ModerationRequest) (*ModerationResponse, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "moderations", request)
	if err != nil {
		return nil, err
	}

	var response ModerationResponse
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateChat submits a chat moderation request.
func (s *ModerationsService) CreateChat(ctx context.Context, messages []Message, model string) (*ModerationResponse, error) {
	request := struct {
		Input []Message `json:"input"`
		Model string    `json:"model"`
	}{
		Input: messages,
		Model: model,
	}

	req, err := s.client.newRequest(ctx, http.MethodPost, "chat/moderations", request)
	if err != nil {
		return nil, err
	}

	var response ModerationResponse
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
