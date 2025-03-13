package mistral

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FIMService handles communication with the Fill-in-the-Middle related methods of the Mistral AI API.
type FIMService struct {
	client *Client
}

// Create submits a Fill-in-the-Middle request.
func (s *FIMService) Create(ctx context.Context, request *FIMRequest) (*FIMResponse, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "fim", request)
	if err != nil {
		return nil, err
	}

	var response FIMResponse
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateStream submits a streaming Fill-in-the-Middle request.
func (s *FIMService) CreateStream(ctx context.Context, request *FIMRequest) (*FIMStream, error) {
	if !request.Stream {
		return nil, fmt.Errorf("stream must be set to true for streaming requests")
	}

	req, err := s.client.newRequest(ctx, http.MethodPost, "fim", request)
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

	return &FIMStream{
		reader:   resp.Body,
		response: resp,
	}, nil
}

// FIMStream represents a streaming response from the Fill-in-the-Middle API.
type FIMStream struct {
	reader   io.ReadCloser
	response *http.Response
}

// Recv receives the next message from the stream.
func (s *FIMStream) Recv() (*FIMResponse, error) {
	if s.reader == nil {
		return nil, fmt.Errorf("stream is closed")
	}

	var response FIMResponse
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
func (s *FIMStream) Close() error {
	if s.reader != nil {
		err := s.reader.Close()
		s.reader = nil
		return err
	}
	return nil
}
