# mistral-api-go

A Go client library for accessing the [Mistral AI API](https://docs.mistral.ai/).

## Installation

```bash
go get github.com/tforrest/mistral-api-go
```

## Usage

```go
import "github.com/tforrest/mistral-api-go"

// Example usage
// Create a new client
client, err := mistral.NewClient("your-api-key")
if err != nil {
    log.Fatal(err)
}

// Create a chat completion
ctx := context.Background()
resp, err := client.Chat.Create(ctx, &mistral.ChatCompletionRequest{
    Model: "mistral-small-latest",
    Messages: []mistral.Message{
        {
            Role:    "user",
            Content: "What is the capital of France?",
        },
    },
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(resp.Choices[0].Message.Content)
```

## Services

The client is divided into several services:

- `Chat`: Chat completions
- `Models`: Model management with version tracking and performance metrics
- `Embeddings`: Text embeddings with similarity and batch processing
- `Files`: File operations
- `FineTuning`: Fine-tuning jobs
- `Moderations`: Content moderation
- `OCR`: Optical Character Recognition
- `Agents`: AI agents with custom tool support
- `FIM`: Fill-in-the-Middle text generation
- `Classifiers`: Text classification with multi-label support
- `Batch`: Efficient batch processing of requests

## Configuration

You can configure the client using options:

```go
client, err := mistral.NewClient(
    "your-api-key",
    mistral.WithBaseURL("https://custom-url.com"),
    mistral.WithHTTPClient(&http.Client{
        Timeout: time.Second * 30,
    }),
)
```

## Error Handling

The library returns detailed error messages from the API:

```go
_, err := client.Chat.Create(ctx, req)
if err != nil {
    if apiErr, ok := err.(*mistral.Error); ok {
        fmt.Printf("API error: %v\n", apiErr.Message)
        fmt.Printf("Status code: %d\n", apiErr.Response.StatusCode)
    }
    return err
}
```

## Usage Examples

```go
// Fill-in-the-Middle example
fimResp, err := client.FIM.Create(ctx, &mistral.FIMRequest{
    Model:  "mistral-large-latest",
    Prefix: "The quick brown",
    Suffix: "jumped over the lazy dog",
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(fimResp.Choices[0].Text)

// Text classification example
classResp, err := client.Classifiers.Create(ctx, &mistral.ClassifierRequest{
    Model:  "mistral-classify",
    Input:  []string{"This movie was fantastic!"},
    Labels: []string{"positive", "negative", "neutral"},
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Classification: %s (confidence: %.2f)\n",
    classResp.Results[0].Labels[0],
    classResp.Results[0].Confidence)

// Enhanced model information
modelInfo, err := client.Models.GetEnhanced(ctx, "mistral-large-latest")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Model version: %s\n", modelInfo.Version)
fmt.Printf("Max tokens: %d\n", modelInfo.MaxTokens)
for _, cap := range modelInfo.Capabilities {
    fmt.Printf("Capability: %s (%v)\n", cap.Name, cap.Available)
}

// Token estimation
tokenCount, err := client.Models.EstimateTokens(ctx, "mistral-large-latest",
    "How many tokens are in this text?")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Token count: %d\n", tokenCount)

// Batch processing example
texts := []string{"Text 1", "Text 2", "Text 3", /* ... */}
batchResp, err := client.Batch.CreateEmbeddings(ctx, texts, "mistral-embed", &mistral.BatchOptions{
    MaxConcurrency: 5,
    ChunkSize:      1000,
    RetryConfig: &mistral.Retry{
        MaxAttempts:  3,
        InitialDelay: 1000,
        MaxDelay:     5000,
    },
})
if err != nil {
    log.Fatal(err)
}

// Wait for batch completion
result, err := client.Batch.WaitForCompletion(ctx, batchResp.ID)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Processed %d/%d requests successfully\n",
    result.Summary.Succeeded,
    result.Summary.TotalRequests)
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This library is distributed under the MIT license.
