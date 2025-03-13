package mistral

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// AgentService handles communication with the agent-related methods of the Mistral AI API.
type AgentService struct {
	client *Client
}

// Create starts a new agent chat session.
func (s *AgentService) Create(ctx context.Context, request *AgentRequest) (*AgentResponse, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "agents/chat", request)
	if err != nil {
		return nil, err
	}

	var response AgentResponse
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateStream starts a new streaming agent chat session.
func (s *AgentService) CreateStream(ctx context.Context, request *AgentRequest) (*AgentStream, error) {
	if !request.Stream {
		return nil, fmt.Errorf("stream must be set to true for streaming requests")
	}

	req, err := s.client.newRequest(ctx, http.MethodPost, "agents/chat", request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	resp, err := s.client.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return &AgentStream{
		reader:   resp.Body,
		response: resp,
	}, nil
}

// AgentStream represents a streaming response from the agent chat API.
type AgentStream struct {
	reader   io.ReadCloser
	response *http.Response
}

// Recv receives the next message from the stream.
func (s *AgentStream) Recv() (*AgentResponse, error) {
	if s.reader == nil {
		return nil, fmt.Errorf("stream is closed")
	}

	var response AgentResponse
	decoder := json.NewDecoder(s.reader)

	var chunk struct {
		Data string `json:"data"`
	}
	err := decoder.Decode(&chunk)
	if err == io.EOF {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	if chunk.Data == "[DONE]" {
		return nil, io.EOF
	}

	err = json.Unmarshal([]byte(chunk.Data), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Close closes the response body.
func (s *AgentStream) Close() error {
	if s.reader != nil {
		err := s.reader.Close()
		s.reader = nil
		return err
	}
	return nil
}
