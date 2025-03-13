package mistral

import (
	"context"
	"net/http"
)

// OCRService handles communication with the OCR related methods of the Mistral AI API.
type OCRService struct {
	client *Client
}

// Process submits an OCR request
func (s *OCRService) Process(request *OCRRequest) (*OCRResponse, error) {
	ctx := context.Background()
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
