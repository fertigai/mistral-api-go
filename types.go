package mistral

type OCRRequestModel struct {
	Type         string `json:"type" enum:"document_url"`
	DocumentURL  string `json:"document_url,omitempty" binding:"required_if=Type document_url"`
	ImageURL     string `json:"image_url,omitempty" binding:"required_if=Type image_url"`
	DocumentName string `json:"document_name,omitempty"`
}

// OCRRequest represents a request to perform OCR on an image.
type OCRRequest struct {
	Model              string          `json:"model"`
	ID                 string          `json:"id"`
	Document           OCRRequestModel `json:"document"`
	Pages              []int           `json:"pages"`
	IncludeImageBase64 bool            `json:"include_image_base64"`
	ImageLimit         int             `json:"image_limit"`
	ImageMinSize       int             `json:"image_min_size"`
}

type OCRImage struct {
	ID           int    `json:"id"`
	TopLeftX     int    `json:"top_left_x"`
	TopLeftY     int    `json:"top_left_y"`
	BottomRightX int    `json:"bottom_right_x"`
	BottomRightY int    `json:"bottom_right_y"`
	ImageBase64  string `json:"image_base64"`
}

type Dimensions struct {
	DPI    int `json:"dpi"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type OCRPage struct {
	Index      int        `json:"index"`
	Markdown   string     `json:"markdown"`
	Images     []OCRImage `json:"images"`
	Dimensions Dimensions `json:"dimensions"`
}

type OCRUsage struct {
	PagesProcessed int `json:"pages_processed"`
	DocSizeBytes   int `json:"doc_size_bytes"`
}

type OCRResponse struct {
	Pages []OCRPage `json:"pages" binding:"required"`
	Model string    `json:"model" binding:"required"`
	Usage OCRUsage  `json:"usage" binding:"required"`
}
