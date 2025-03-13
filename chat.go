package mistral

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ChatService handles communication with the chat related methods of the Mistral AI API.
type ChatService struct {
	client *Client
}

// Create submits a chat completion request.
func (s *ChatService) Create(ctx context.Context, request *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "chat/completions", request)
	if err != nil {
		return nil, err
	}

	var response ChatCompletionResponse
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateStream submits a chat completion request and returns a stream of responses.
func (s *ChatService) CreateStream(ctx context.Context, request *ChatCompletionRequest) (*ChatCompletionStream, error) {
	if !request.Stream {
		return nil, fmt.Errorf("stream must be set to true for streaming requests")
	}

	req, err := s.client.newRequest(ctx, http.MethodPost, "chat/completions", request)
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

	return &ChatCompletionStream{
		reader:   resp.Body,
		response: resp,
	}, nil
}

// ChatCompletionStream represents a streaming response from the chat completions API.
type ChatCompletionStream struct {
	reader   io.ReadCloser
	response *http.Response
}

// Recv receives the next message from the stream.
func (s *ChatCompletionStream) Recv() (*ChatCompletionResponse, error) {
	if s.reader == nil {
		return nil, fmt.Errorf("stream is closed")
	}

	var response ChatCompletionResponse
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
func (s *ChatCompletionStream) Close() error {
	if s.reader != nil {
		err := s.reader.Close()
		s.reader = nil
		return err
	}
	return nil
}
