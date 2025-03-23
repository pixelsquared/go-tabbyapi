// Package tabby provides a Go client for TabbyAPI.
package tabby

import (
	"net/http"
)

// Authenticator provides authentication for API requests.
// This interface is implemented by different authentication methods
// including API key, admin key, and bearer token authentication.
type Authenticator interface {
	// Apply adds authentication information to the provided HTTP request.
	// This is called automatically before each request is sent to the TabbyAPI server.
	Apply(req *http.Request)
}

// APIKeyAuthenticator uses the X-API-Key header for standard API access.
// This is the recommended authentication method for most applications.
//
// API keys typically provide read and write permissions to the API, allowing
// operations like generating completions and embeddings, but may restrict
// administrative actions like model management.
type APIKeyAuthenticator struct {
	// Key is the API key value to use for authentication
	Key string
}

// Apply implements the Authenticator interface by setting the X-API-Key header.
func (a *APIKeyAuthenticator) Apply(req *http.Request) {
	req.Header.Set("X-API-Key", a.Key)
}

// AdminKeyAuthenticator uses the X-Admin-Key header for administrative operations.
// This authentication method is required for management operations like loading
// or unloading models, managing LoRA adapters, and other privileged actions.
//
// Admin keys should be kept secure and used only by trusted applications or users
// who need administrative access to the TabbyAPI server.
type AdminKeyAuthenticator struct {
	// Key is the admin key value to use for authentication
	Key string
}

// Apply implements the Authenticator interface by setting the X-Admin-Key header.
func (a *AdminKeyAuthenticator) Apply(req *http.Request) {
	req.Header.Set("X-Admin-Key", a.Key)
}

// BearerTokenAuthenticator uses the Authorization header with Bearer authentication.
// This method is useful when TabbyAPI is configured to use OAuth or JWT authentication,
// or when integrated with other authentication systems.
//
// The permissions associated with a bearer token depend on the token's claims
// and how they are interpreted by the TabbyAPI server's authentication configuration.
type BearerTokenAuthenticator struct {
	// Token is the bearer token value to use for authentication
	Token string
}

// Apply implements the Authenticator interface by setting the Authorization header
// with the Bearer authentication scheme followed by the token.
func (a *BearerTokenAuthenticator) Apply(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+a.Token)
}

// Stream is a generic interface for Server-Sent Events (SSE) streams.
// It provides methods to receive items from the stream and to close the stream
// when it's no longer needed.
type Stream[T any] interface {
	// Recv returns the next item from the stream.
	//
	// This method waits for and returns the next item from the stream. If there
	// are no more items or the stream has been closed, it returns io.EOF.
	// If the context used to create the stream is canceled, it returns the context error.
	//
	// This method is designed to be called in a loop until an error is returned.
	// The returned error can be either io.EOF (end of stream), ErrStreamClosed
	// (stream was closed), or any other error that occurred while reading from the stream.
	Recv() (T, error)

	// Close releases resources associated with the stream.
	//
	// This method should always be called when done with the stream to prevent
	// resource leaks. It's safe to call Close multiple times.
	//
	// Best practice is to defer the Close call immediately after creating the stream:
	//   stream, err := client.Completions().CreateStream(ctx, req)
	//   if err != nil { return err }
	//   defer stream.Close()
	Close() error
}

// CompletionStream is a stream of incremental text completion responses.
type CompletionStream = Stream[*CompletionStreamResponse]

// ChatCompletionStream is a stream of incremental chat completion responses.
type ChatCompletionStream = Stream[*ChatCompletionStreamResponse]

// ModelLoadStream is a stream of model loading progress updates.
type ModelLoadStream = Stream[*ModelLoadResponse]

// ChatMessageRole defines the role of a message in a chat completion
type ChatMessageRole string

const (
	// ChatMessageRoleUser represents a user message
	ChatMessageRoleUser ChatMessageRole = "user"

	// ChatMessageRoleAssistant represents an assistant message
	ChatMessageRoleAssistant ChatMessageRole = "assistant"

	// ChatMessageRoleSystem represents a system message
	ChatMessageRoleSystem ChatMessageRole = "system"

	// ChatMessageRoleTool represents a tool message
	ChatMessageRoleTool ChatMessageRole = "tool"
)

// ChatMessage represents a message in a chat completion request
type ChatMessage struct {
	Role    ChatMessageRole `json:"role"`
	Content interface{}     `json:"content"` // String or array of ChatMessageContent
	// Additional fields may be added later
}

// ChatMessageContent represents a part of a message content
type ChatMessageContent struct {
	Type     string        `json:"type"`
	Text     string        `json:"text,omitempty"`
	ImageURL *ChatImageURL `json:"image_url,omitempty"`
}

// ChatImageURL represents an image URL in a chat message
type ChatImageURL struct {
	URL string `json:"url"`
}

// CompletionRequest matches the TabbyAPI completion request schema
type CompletionRequest struct {
	Prompt      interface{} `json:"prompt"` // String or array of strings
	MaxTokens   int         `json:"max_tokens,omitempty"`
	Temperature float64     `json:"temperature,omitempty"`
	TopP        float64     `json:"top_p,omitempty"`
	TopK        int         `json:"top_k,omitempty"`
	Stream      bool        `json:"stream,omitempty"`
	Stop        interface{} `json:"stop,omitempty"` // String or array of strings
	Model       string      `json:"model,omitempty"`
	// Additional parameters will be added as needed
}

// CompletionResponse represents a response to a completion request
type CompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []CompletionRespChoice `json:"choices"`
	Usage   *UsageStats            `json:"usage,omitempty"`
}

// CompletionRespChoice represents a choice in a completion response
type CompletionRespChoice struct {
	Text         string              `json:"text"`
	Index        int                 `json:"index"`
	FinishReason string              `json:"finish_reason,omitempty"`
	LogProbs     *CompletionLogProbs `json:"logprobs,omitempty"`
}

// CompletionLogProbs represents log probabilities for a completion
type CompletionLogProbs struct {
	Tokens        []string             `json:"tokens"`
	TokenLogProbs []float64            `json:"token_logprobs"`
	TopLogProbs   []map[string]float64 `json:"top_logprobs"`
	TextOffset    []int                `json:"text_offset"`
}

// CompletionStreamResponse represents a streaming completion response
type CompletionStreamResponse struct {
	ID      string                   `json:"id"`
	Object  string                   `json:"object"`
	Created int64                    `json:"created"`
	Model   string                   `json:"model"`
	Choices []CompletionStreamChoice `json:"choices"`
}

// CompletionStreamChoice represents a streaming choice
type CompletionStreamChoice struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason,omitempty"`
}

// ChatCompletionRequest matches the TabbyAPI chat completion request schema
type ChatCompletionRequest struct {
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	TopP        float64       `json:"top_p,omitempty"`
	TopK        int           `json:"top_k,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
	Stop        interface{}   `json:"stop,omitempty"` // String or array of strings
	Model       string        `json:"model,omitempty"`
	// Additional parameters will be added as needed
}

// ChatCompletionResponse represents a response to a chat completion request
type ChatCompletionResponse struct {
	ID      string                     `json:"id"`
	Object  string                     `json:"object"`
	Created int64                      `json:"created"`
	Model   string                     `json:"model"`
	Choices []ChatCompletionRespChoice `json:"choices"`
	Usage   *UsageStats                `json:"usage,omitempty"`
}

// ChatCompletionRespChoice represents a choice in a chat completion response
type ChatCompletionRespChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason,omitempty"`
}

// ChatCompletionStreamResponse represents a streaming chat completion response
type ChatCompletionStreamResponse struct {
	ID      string                       `json:"id"`
	Object  string                       `json:"object"`
	Created int64                        `json:"created"`
	Model   string                       `json:"model"`
	Choices []ChatCompletionStreamChoice `json:"choices"`
}

// ChatCompletionStreamChoice represents a streaming chat choice
type ChatCompletionStreamChoice struct {
	Index        int    `json:"index"`
	Delta        *Delta `json:"delta"`
	FinishReason string `json:"finish_reason,omitempty"`
}

// Delta represents a delta in a streaming chat response
type Delta struct {
	Role    ChatMessageRole `json:"role,omitempty"`
	Content string          `json:"content,omitempty"`
}

// UsageStats represents token usage information
type UsageStats struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// EmbeddingsRequest represents a request to create embeddings
type EmbeddingsRequest struct {
	Input          interface{} `json:"input"` // String or array of strings
	Model          string      `json:"model,omitempty"`
	EncodingFormat string      `json:"encoding_format,omitempty"`
}

// EmbeddingsResponse represents a response to an embeddings request
type EmbeddingsResponse struct {
	Object string            `json:"object"`
	Data   []EmbeddingObject `json:"data"`
	Model  string            `json:"model"`
	Usage  UsageInfo         `json:"usage"`
}

// EmbeddingObject represents a single embedding result
type EmbeddingObject struct {
	Object    string      `json:"object"`
	Embedding interface{} `json:"embedding"` // Array of floats or base64 string
	Index     int         `json:"index"`
}

// UsageInfo represents token usage for embeddings
type UsageInfo struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// ModelCard represents information about a model
type ModelCard struct {
	ID         string               `json:"id"`
	Object     string               `json:"object"`
	Created    int64                `json:"created"`
	OwnedBy    string               `json:"owned_by"`
	Parameters *ModelCardParameters `json:"parameters,omitempty"`
}

// ModelCardParameters represents model parameters
type ModelCardParameters struct {
	MaxSeqLen      int     `json:"max_seq_len,omitempty"`
	RopeScale      float64 `json:"rope_scale,omitempty"`
	RopeAlpha      float64 `json:"rope_alpha,omitempty"`
	MaxBatchSize   int     `json:"max_batch_size,omitempty"`
	CacheSize      int     `json:"cache_size,omitempty"`
	CacheMode      string  `json:"cache_mode,omitempty"`
	ChunkSize      int     `json:"chunk_size,omitempty"`
	PromptTemplate string  `json:"prompt_template,omitempty"`
	UseVision      bool    `json:"use_vision,omitempty"`
}

// ModelList represents a list of models
type ModelList struct {
	Object string      `json:"object"`
	Data   []ModelCard `json:"data"`
}

// ModelLoadRequest represents a request to load a model
type ModelLoadRequest struct {
	ModelName      string      `json:"model_name"`
	MaxSeqLen      int         `json:"max_seq_len,omitempty"`
	RopeScale      float64     `json:"rope_scale,omitempty"`
	RopeAlpha      interface{} `json:"rope_alpha,omitempty"` // float64 or "auto"
	GPUSplit       []float64   `json:"gpu_split,omitempty"`
	CacheSize      int         `json:"cache_size,omitempty"`
	CacheMode      string      `json:"cache_mode,omitempty"`
	ChunkSize      int         `json:"chunk_size,omitempty"`
	PromptTemplate string      `json:"prompt_template,omitempty"`
	// Additional parameters may be added later
}

// ModelLoadResponse represents a response to a model load request
type ModelLoadResponse struct {
	ModelType string `json:"model_type"`
	Module    int    `json:"module"`
	Modules   int    `json:"modules"`
	Status    string `json:"status"`
}

// ModelPropsResponse represents a response to a model props request
type ModelPropsResponse struct {
	TotalSlots                int                             `json:"total_slots"`
	ChatTemplate              string                          `json:"chat_template"`
	DefaultGenerationSettings *ModelDefaultGenerationSettings `json:"default_generation_settings"`
}

// ModelDefaultGenerationSettings represents default generation settings
type ModelDefaultGenerationSettings struct {
	NCtx int `json:"n_ctx"`
}

// EmbeddingModelLoadRequest represents a request to load an embedding model
type EmbeddingModelLoadRequest struct {
	EmbeddingModelName string `json:"embedding_model_name"`
	EmbeddingsDevice   string `json:"embeddings_device,omitempty"`
}

// DownloadRequest represents a request to download a model
type DownloadRequest struct {
	RepoID     string   `json:"repo_id"`
	RepoType   string   `json:"repo_type,omitempty"`
	FolderName string   `json:"folder_name,omitempty"`
	Revision   string   `json:"revision,omitempty"`
	Token      string   `json:"token,omitempty"`
	Include    []string `json:"include,omitempty"`
	Exclude    []string `json:"exclude,omitempty"`
}

// DownloadResponse represents a response to a download request
type DownloadResponse struct {
	DownloadPath string `json:"download_path"`
}

// LoraCard represents information about a LoRA adapter
type LoraCard struct {
	ID      string  `json:"id"`
	Object  string  `json:"object"`
	Created int64   `json:"created"`
	OwnedBy string  `json:"owned_by"`
	Scaling float64 `json:"scaling,omitempty"`
}

// LoraList represents a list of LoRA adapters
type LoraList struct {
	Object string     `json:"object"`
	Data   []LoraCard `json:"data"`
}

// LoraLoadInfo represents information for loading a LoRA adapter
type LoraLoadInfo struct {
	Name    string  `json:"name"`
	Scaling float64 `json:"scaling,omitempty"`
}

// LoraLoadRequest represents a request to load LoRA adapters
type LoraLoadRequest struct {
	Loras     []LoraLoadInfo `json:"loras"`
	SkipQueue bool           `json:"skip_queue,omitempty"`
}

// LoraLoadResponse represents a response to a LoRA load request
type LoraLoadResponse struct {
	Success []string `json:"success"`
	Failure []string `json:"failure"`
}

// TokenEncodeRequest represents a request to encode tokens
type TokenEncodeRequest struct {
	Text                interface{} `json:"text"` // String or []ChatMessage
	AddBOSToken         bool        `json:"add_bos_token,omitempty"`
	EncodeSpecialTokens bool        `json:"encode_special_tokens,omitempty"`
	DecodeSpecialTokens bool        `json:"decode_special_tokens,omitempty"`
}

// TokenEncodeResponse represents a response to a token encode request
type TokenEncodeResponse struct {
	Tokens []int `json:"tokens"`
	Length int   `json:"length"`
}

// TokenDecodeRequest represents a request to decode tokens
type TokenDecodeRequest struct {
	Tokens              []int `json:"tokens"`
	AddBOSToken         bool  `json:"add_bos_token,omitempty"`
	EncodeSpecialTokens bool  `json:"encode_special_tokens,omitempty"`
	DecodeSpecialTokens bool  `json:"decode_special_tokens,omitempty"`
}

// TokenDecodeResponse represents a response to a token decode request
type TokenDecodeResponse struct {
	Text string `json:"text"`
}

// TemplateList represents a list of templates
type TemplateList struct {
	Object string   `json:"object"`
	Data   []string `json:"data"`
}

// TemplateSwitchRequest represents a request to switch templates
type TemplateSwitchRequest struct {
	PromptTemplateName string `json:"prompt_template_name"`
}

// SamplerOverrideListResponse represents a response to a sampler override list request
type SamplerOverrideListResponse struct {
	SelectedPreset string                 `json:"selected_preset,omitempty"`
	Overrides      map[string]interface{} `json:"overrides"`
	Presets        []string               `json:"presets"`
}

// SamplerOverrideSwitchRequest represents a request to switch sampler overrides
type SamplerOverrideSwitchRequest struct {
	Preset    string                 `json:"preset,omitempty"`
	Overrides map[string]interface{} `json:"overrides,omitempty"`
}

// HealthCheckResponse represents a response to a health check
type HealthCheckResponse struct {
	Status string           `json:"status"`
	Issues []UnhealthyEvent `json:"issues,omitempty"`
}

// UnhealthyEvent represents an issue in a health check
type UnhealthyEvent struct {
	Time        string `json:"time"`
	Description string `json:"description"`
}

// AuthPermissionResponse represents a response to an auth permission check
type AuthPermissionResponse struct {
	Permission string `json:"permission"`
}
