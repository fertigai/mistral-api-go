package mistral

import (
	"context"
	"net/http"
)

// OCRService handles communication with the OCR related methods of the Mistral AI API.
type OCRService struct {
	client *Client
}

// Create submits an OCR request for one or more images.
func (s *OCRService) Create(ctx context.Context, request *OCRRequest) (*OCRResponse, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "ocr", request)
	if err != nil {
		return nil, err
	}

	var response OCRResponse
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateAsync submits an asynchronous OCR request and returns a job ID.
func (s *OCRService) CreateAsync(ctx context.Context, request *OCRRequest) (string, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "ocr/async", request)
	if err != nil {
		return "", err
	}

	var response struct {
		JobID string `json:"job_id"`
	}
	_, err = s.client.do(req, &response)
	if err != nil {
		return "", err
	}

	return response.JobID, nil
}

// GetAsyncResult retrieves the result of an asynchronous OCR request.
func (s *OCRService) GetAsyncResult(ctx context.Context, jobID string) (*OCRResponse, error) {
	path := "ocr/async/" + jobID
	req, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var response OCRResponse
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
