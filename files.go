package mistral

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/samber/lo"
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

type SampleType string

const (
	SampleTypePretrain     SampleType = "pretrain"
	SampleTypeInstruct     SampleType = "instruct"
	SampleTypeBatchRequest SampleType = "batch_request"
	SampleTypeBatchResult  SampleType = "batch_result"
	SampleTypeBatchError   SampleType = "batch_error"
	SampleTypeOCRInput     SampleType = "ocr_input"
)

type GetFilesRequest struct {
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	Source     string       `json:"source"`
	Search     string       `json:"search"`
	Purpose    string       `json:"purpose"`
	SampleType []SampleType `json:"sample_type"`
}

type GetFilesResponse struct {
	Data   []map[string]interface{} `json:"data"`
	Total  int                      `json:"total"`
	Object string                   `json:"object"`
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

func (s *FilesService) GetFiles(getFilesRequest *GetFilesRequest) (*GetFilesResponse, error) {
	var response GetFilesResponse
	resp, err := s.client.client.R().
		SetHeader("Authorization", s.client.AuthHeader()).
		SetResult(&response).
		SetQueryParams(map[string]string{
			"page":        strconv.Itoa(getFilesRequest.Page),
			"page_size":   strconv.Itoa(getFilesRequest.PageSize),
			"source":      getFilesRequest.Source,
			"search":      getFilesRequest.Search,
			"purpose":     getFilesRequest.Purpose,
			"sample_type": strings.Join(lo.Map(getFilesRequest.SampleType, func(s SampleType, _ int) string { return string(s) }), ","),
		}).
		Get(s.client.FormUrl("files"))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get file: %s", resp.String())
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
