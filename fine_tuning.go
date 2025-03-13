package mistral

import (
	"context"
	"fmt"
	"net/http"
)

// FineTuningService handles communication with the fine-tuning related methods of the Mistral AI API.
type FineTuningService struct {
	client *Client
}

// FineTuningJob represents a fine-tuning job.
type FineTuningJob struct {
	ID              string                 `json:"id"`
	Model           string                 `json:"model"`
	Status          string                 `json:"status"`
	TrainingFiles   []string               `json:"training_files"`
	ValidationFiles []string               `json:"validation_files,omitempty"`
	Hyperparameters map[string]interface{} `json:"hyperparameters"`
	ResultFiles     []string               `json:"result_files,omitempty"`
	CreatedAt       int64                  `json:"created_at"`
	FinishedAt      int64                  `json:"finished_at,omitempty"`
}

// FineTuningJobList represents a list of fine-tuning jobs.
type FineTuningJobList struct {
	Object string          `json:"object"`
	Data   []FineTuningJob `json:"data"`
}

// Create creates a new fine-tuning job.
func (s *FineTuningService) Create(ctx context.Context, model string, trainingFiles []string, hyperparameters map[string]interface{}) (*FineTuningJob, error) {
	request := map[string]interface{}{
		"model":           model,
		"training_files":  trainingFiles,
		"hyperparameters": hyperparameters,
	}

	req, err := s.client.newRequest(ctx, http.MethodPost, "fine_tuning/jobs", request)
	if err != nil {
		return nil, err
	}

	var response FineTuningJob
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// List returns a list of fine-tuning jobs.
func (s *FineTuningService) List(ctx context.Context) (*FineTuningJobList, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "fine_tuning/jobs", nil)
	if err != nil {
		return nil, err
	}

	var response FineTuningJobList
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Get retrieves a specific fine-tuning job.
func (s *FineTuningService) Get(ctx context.Context, jobID string) (*FineTuningJob, error) {
	path := fmt.Sprintf("fine_tuning/jobs/%s", jobID)
	req, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var response FineTuningJob
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Cancel cancels a fine-tuning job.
func (s *FineTuningService) Cancel(ctx context.Context, jobID string) (*FineTuningJob, error) {
	path := fmt.Sprintf("fine_tuning/jobs/%s/cancel", jobID)
	req, err := s.client.newRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	var response FineTuningJob
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
