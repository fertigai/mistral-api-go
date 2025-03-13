package mistral

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type FilesService struct {
	client *Client
}

type FileUploadResponse struct {
	ID         string `json:"id"`
	Object     string `json:"object"`
	Bytes      int    `json:"bytes"`
	CreatedAt  int64  `json:"created_at"`
	Filename   string `json:"filename"`
	Purpose    string `json:"purpose"`
	SampleType string `json:"sample_type"`
	NumLines   *int   `json:"num_lines,omitempty"`
	Source     string `json:"source"`
}

type FileType string

const (
	FileTypePDF   FileType = "pdf"
	FileTypeImage FileType = "image"
)

type FileUploadRequest struct {
	File    *os.File
	Purpose string
	Type    FileType
}

type DocumentUrlResponse struct {
	Url string `json:"url"`
}

func (s *FilesService) GetDocumentUrl(id string, expiry int) (*DocumentUrlResponse, error) {
	var response DocumentUrlResponse
	// todo validate expiry
	resp, err := s.client.client.R().
		SetResult(&response).
		SetHeader("Authorization", s.client.AuthHeader()).
		SetQueryParam("expiry", strconv.Itoa(expiry)).
		Get(s.client.FormUrl("files/" + id + "/url"))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get document url: %s", resp.String())
	}

	return &response, nil
}

func (s *FilesService) Upload(fileUploadRequest *FileUploadRequest) (*FileUploadResponse, error) {
	var response FileUploadResponse
	request := s.client.client.R().
		SetFileReader("file", fileUploadRequest.File.Name(), fileUploadRequest.File).
		SetFormData(map[string]string{
			"purpose": "ocr",
		}).
		SetHeader("Content-Type", "multipart/form-data").
		SetResult(&response).
		SetHeader("Authorization", s.client.AuthHeader())

	if fileUploadRequest.Type == FileTypePDF {
		request.SetHeader("Content-Type", "application/pdf")
	} else if fileUploadRequest.Type == FileTypeImage {
		request.SetHeader("Content-Type", "image/jpeg")
	} else {
		return nil, fmt.Errorf("invalid file type: %s", fileUploadRequest.Type)
	}

	resp, err := request.Post(s.client.FormUrl("files"))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to upload file: %s", resp.String())
	}

	return &response, nil
}
