package mistral

import "time"

// UsageInfo represents token usage information for an API call.
type UsageInfo struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Message represents a chat message in the conversation.
type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

// ChatCompletionRequest represents a request to the chat completions API.
type ChatCompletionRequest struct {
	Model          string    `json:"model"`
	Messages       []Message `json:"messages"`
	Temperature    *float32  `json:"temperature,omitempty"`
	TopP           *float32  `json:"top_p,omitempty"`
	MaxTokens      *int      `json:"max_tokens,omitempty"`
	Stream         bool      `json:"stream"`
	Stop           []string  `json:"stop,omitempty"`
	RandomSeed     *int      `json:"random_seed,omitempty"`
	Tools          []Tool    `json:"tools,omitempty"`
	ToolChoice     *string   `json:"tool_choice,omitempty"`
	ResponseFormat *string   `json:"response_format,omitempty"`
}

// ChatCompletionResponse represents a response from the chat completions API.
type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   UsageInfo              `json:"usage"`
}

// ChatCompletionChoice represents a completion choice in a chat response.
type ChatCompletionChoice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Tool represents a function that can be called by the model.
type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// Function represents a callable function in a tool.
type Function struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// EmbeddingRequest represents a request to create embeddings.
type EmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

// EmbeddingResponse represents a response from the embeddings API.
type EmbeddingResponse struct {
	Object string          `json:"object"`
	Data   []EmbeddingData `json:"data"`
	Model  string          `json:"model"`
	Usage  UsageInfo       `json:"usage"`
}

// EmbeddingData represents a single embedding in the response.
type EmbeddingData struct {
	Object    string    `json:"object"`
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
}

// Model represents a model in the API.
type Model struct {
	ID          string   `json:"id"`
	Object      string   `json:"object"`
	Created     int64    `json:"created"`
	OwnedBy     string   `json:"owned_by"`
	Permissions []string `json:"permissions"`
	Root        string   `json:"root"`
	Parent      string   `json:"parent"`
}

// ModelList represents a list of models.
type ModelList struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

// File represents a file in the API.
type File struct {
	ID        string    `json:"id"`
	Object    string    `json:"object"`
	Bytes     int64     `json:"bytes"`
	CreatedAt time.Time `json:"created_at"`
	Filename  string    `json:"filename"`
	Purpose   string    `json:"purpose"`
}

// FileList represents a list of files.
type FileList struct {
	Object string `json:"object"`
	Data   []File `json:"data"`
}

// ModerationRequest represents a request to the moderation API.
type ModerationRequest struct {
	Input []string `json:"input"`
	Model string   `json:"model"`
}

// ModerationResponse represents a response from the moderation API.
type ModerationResponse struct {
	ID      string             `json:"id"`
	Model   string             `json:"model"`
	Results []ModerationResult `json:"results"`
}

// ModerationResult represents a single moderation result.
type ModerationResult struct {
	Categories     map[string]bool    `json:"categories"`
	CategoryScores map[string]float32 `json:"category_scores"`
}

// OCRRequest represents a request to perform OCR on an image.
type OCRRequest struct {
	Model     string   `json:"model"`
	Files     []string `json:"files"`     // File IDs of uploaded images
	Languages []string `json:"languages"` // Optional list of languages to detect
}

// OCRResponse represents a response from the OCR API.
type OCRResponse struct {
	ID      string      `json:"id"`
	Object  string      `json:"object"`
	Created int64       `json:"created"`
	Model   string      `json:"model"`
	Results []OCRResult `json:"results"`
	Usage   UsageInfo   `json:"usage"`
}

// OCRResult represents the OCR result for a single image.
type OCRResult struct {
	FileID   string     `json:"file_id"`
	Text     string     `json:"text"`     // Extracted text
	Language string     `json:"language"` // Detected language
	Blocks   []OCRBlock `json:"blocks"`   // Text blocks with position information
}

// OCRBlock represents a block of text with its position in the image.
type OCRBlock struct {
	Text        string      `json:"text"`
	Confidence  float32     `json:"confidence"`
	BoundingBox BoundingBox `json:"bounding_box"`
}

// BoundingBox represents the position of text in an image.
type BoundingBox struct {
	X      float32 `json:"x"`      // X coordinate of top-left corner
	Y      float32 `json:"y"`      // Y coordinate of top-left corner
	Width  float32 `json:"width"`  // Width of the box
	Height float32 `json:"height"` // Height of the box
}

// AgentAction represents an action that an agent can take.
type AgentAction struct {
	Tool     string                 `json:"tool"`
	Input    map[string]interface{} `json:"input"`
	Thought  string                 `json:"thought,omitempty"`
	Response string                 `json:"response,omitempty"`
}

// AgentFunction represents a function that can be called by an agent.
type AgentFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// AgentTool represents a tool available to an agent.
type AgentTool struct {
	Type     string        `json:"type"`
	Function AgentFunction `json:"function"`
}

// AgentRequest represents a request to create an agent chat.
type AgentRequest struct {
	Model       string      `json:"model"`
	Tools       []AgentTool `json:"tools"`
	Messages    []Message   `json:"messages"`
	MaxActions  *int        `json:"max_actions,omitempty"`
	Stream      bool        `json:"stream"`
	Temperature *float32    `json:"temperature,omitempty"`
	TopP        *float32    `json:"top_p,omitempty"`
	MaxTokens   *int        `json:"max_tokens,omitempty"`
	Stop        []string    `json:"stop,omitempty"`
}

// AgentResponse represents a response from the agent chat API.
type AgentResponse struct {
	ID       string        `json:"id"`
	Object   string        `json:"object"`
	Created  int64         `json:"created"`
	Model    string        `json:"model"`
	Actions  []AgentAction `json:"actions"`
	Messages []Message     `json:"messages"`
	Usage    UsageInfo     `json:"usage"`
}

// EnhancedEmbeddingRequest represents a request to create embeddings with additional options.
type EnhancedEmbeddingRequest struct {
	Model          string   `json:"model"`
	Input          []string `json:"input"`
	EncodingFormat string   `json:"encoding_format,omitempty"` // e.g., "float", "base64"
	Normalize      bool     `json:"normalize,omitempty"`       // Whether to normalize the embeddings
	Truncate       bool     `json:"truncate,omitempty"`        // Whether to truncate inputs that are too long
}

// EnhancedEmbeddingResponse represents a response from the embeddings API with additional metadata.
type EnhancedEmbeddingResponse struct {
	Object   string            `json:"object"`
	Data     []EmbeddingData   `json:"data"`
	Model    string            `json:"model"`
	Usage    UsageInfo         `json:"usage"`
	Metadata EmbeddingMetadata `json:"metadata,omitempty"`
}

// EmbeddingMetadata represents additional metadata about the embeddings.
type EmbeddingMetadata struct {
	Dimensions int    `json:"dimensions"`
	Similarity string `json:"similarity_metric"` // e.g., "cosine", "euclidean"
	Truncated  []bool `json:"truncated,omitempty"`
}

// FIMRequest represents a request to the Fill-in-the-Middle API.
type FIMRequest struct {
	Model       string   `json:"model"`
	Prefix      string   `json:"prefix"`
	Suffix      string   `json:"suffix"`
	MaxTokens   *int     `json:"max_tokens,omitempty"`
	Temperature *float32 `json:"temperature,omitempty"`
	TopP        *float32 `json:"top_p,omitempty"`
	Stream      bool     `json:"stream"`
	Stop        []string `json:"stop,omitempty"`
}

// FIMResponse represents a response from the Fill-in-the-Middle API.
type FIMResponse struct {
	ID      string      `json:"id"`
	Object  string      `json:"object"`
	Created int64       `json:"created"`
	Model   string      `json:"model"`
	Choices []FIMChoice `json:"choices"`
	Usage   UsageInfo   `json:"usage"`
}

// FIMChoice represents a completion choice in a FIM response.
type FIMChoice struct {
	Index        int    `json:"index"`
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason"`
}

// ClassifierRequest represents a request to the classifier API.
type ClassifierRequest struct {
	Model       string   `json:"model"`
	Input       []string `json:"input"`
	Labels      []string `json:"labels"`
	MultiLabel  bool     `json:"multi_label,omitempty"`
	Temperature *float32 `json:"temperature,omitempty"`
}

// ClassifierResponse represents a response from the classifier API.
type ClassifierResponse struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Created int64              `json:"created"`
	Model   string             `json:"model"`
	Results []ClassifierResult `json:"results"`
	Usage   UsageInfo          `json:"usage"`
}

// ClassifierResult represents a single classification result.
type ClassifierResult struct {
	Input      string             `json:"input"`
	Labels     []string           `json:"labels"`
	Scores     map[string]float32 `json:"scores"`
	Confidence float32            `json:"confidence"`
}

// ModelCapability represents a capability of a model.
type ModelCapability struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Available   bool   `json:"available"`
}

// ModelPerformance represents performance metrics for a model.
type ModelPerformance struct {
	Metric string  `json:"metric"`
	Value  float32 `json:"value"`
	Unit   string  `json:"unit"`
}

// EnhancedModel extends the base Model type with additional information.
type EnhancedModel struct {
	Model
	Version      string             `json:"version"`
	Description  string             `json:"description"`
	Capabilities []ModelCapability  `json:"capabilities"`
	Performance  []ModelPerformance `json:"performance"`
	MaxTokens    int                `json:"max_tokens"`
	TokenCosts   struct {
		Input  float32 `json:"input"`
		Output float32 `json:"output"`
	} `json:"token_costs"`
}

// BatchRequest represents a request for batch processing.
type BatchRequest struct {
	Requests []interface{} `json:"requests"` // Can hold any request type
	Model    string        `json:"model"`
	Options  BatchOptions  `json:"options"`
}

// BatchOptions represents options for batch processing.
type BatchOptions struct {
	MaxConcurrency int    `json:"max_concurrency,omitempty"`
	RetryConfig    *Retry `json:"retry_config,omitempty"`
	Timeout        int    `json:"timeout,omitempty"`    // In seconds
	ChunkSize      int    `json:"chunk_size,omitempty"` // For large batches
}

// Retry represents retry configuration for batch processing.
type Retry struct {
	MaxAttempts  int `json:"max_attempts"`
	InitialDelay int `json:"initial_delay"` // In milliseconds
	MaxDelay     int `json:"max_delay"`     // In milliseconds
}

// BatchResponse represents a response from batch processing.
type BatchResponse struct {
	ID      string        `json:"id"`
	Object  string        `json:"object"`
	Created int64         `json:"created"`
	Model   string        `json:"model"`
	Results []BatchResult `json:"results"`
	Summary BatchSummary  `json:"summary"`
}

// BatchResult represents a single result in a batch response.
type BatchResult struct {
	Index    int         `json:"index"`
	Status   string      `json:"status"`
	Response interface{} `json:"response,omitempty"`
	Error    string      `json:"error,omitempty"`
	Attempts int         `json:"attempts"`
	Duration int         `json:"duration"` // In milliseconds
}

// BatchSummary represents summary statistics for a batch operation.
type BatchSummary struct {
	TotalRequests   int     `json:"total_requests"`
	Succeeded       int     `json:"succeeded"`
	Failed          int     `json:"failed"`
	TotalDuration   int     `json:"total_duration"`   // In milliseconds
	AverageDuration int     `json:"average_duration"` // In milliseconds
	TokensUsed      int     `json:"tokens_used"`
	TotalCost       float32 `json:"total_cost"`
}
