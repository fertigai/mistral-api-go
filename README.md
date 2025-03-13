# mistral-api-go

A Go client library for accessing the [Mistral AI API](https://docs.mistral.ai/)'s OCR capabilities.

Still under development and has limited functionality. Mostly focused on the OCR + File Upload endpoints.

## Installation

```bash
go get github.com/tforrest/mistral-api-go
```

## Usage

```go
import "github.com/tforrest/mistral-api-go"

// Create a new client
client, err := mistral.NewClient("your-api-key")
if err != nil {
    log.Fatal(err)
}

// Process a document with OCR
resp, err := client.OCR.Process(&OCRRequest{
    Document: OCRRequestModel{
        Type:         "document_url",
        DocumentURL:  "https://arxiv.org/pdf/2201.04234",
        DocumentName: "2201.04234.pdf",
    },
    Model: "mistral-ocr-latest",
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Processed %d pages\n", len(resp.Pages))
```
