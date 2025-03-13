package mistral

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// FilesService handles communication with the files related methods of the Mistral AI API.
type FilesService struct {
	client *Client
}

// List returns a list of files that belong to the user's organization.
func (s *FilesService) List(ctx context.Context) (*FileList, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "files", nil)
	if err != nil {
		return nil, err
	}

	var response FileList
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Get retrieves a specific file.
func (s *FilesService) Get(ctx context.Context, fileID string) (*File, error) {
	path := fmt.Sprintf("files/%s", fileID)
	req, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var response File
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Delete deletes a file.
func (s *FilesService) Delete(ctx context.Context, fileID string) error {
	path := fmt.Sprintf("files/%s", fileID)
	req, err := s.client.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(req, nil)
	return err
}

// Upload uploads a new file.
func (s *FilesService) Upload(ctx context.Context, filePath string, purpose string) (*File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	// Add the purpose
	err = writer.WriteField("purpose", purpose)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/v1/files", s.client.baseURL), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.client.apiKey))

	var response File
	_, err = s.client.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Download downloads a file's contents.
func (s *FilesService) Download(ctx context.Context, fileID string) (io.ReadCloser, error) {
	path := fmt.Sprintf("files/%s/content", fileID)
	req, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
