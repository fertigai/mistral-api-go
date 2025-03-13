package mistral

import (
	"fmt"
	"net/http"
)

// OCRService handles communication with the OCR related methods of the Mistral AI API.
type OCRService struct {
	client *Client
}

// Process submits an OCR request
func (s *OCRService) Process(request *OCRRequest) (*OCRResponse, error) {
	var response OCRResponse
	resp, err := s.client.client.R().
		SetBody(request).
		SetResult(&response).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", s.client.AuthHeader()).
		Post(s.client.FormUrl("ocr"))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to process OCR: %s", resp.String())
	}

	return &response, nil
}
