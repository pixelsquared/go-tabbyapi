// Package tabby provides a Go client for TabbyAPI, an open-source self-hosted AI coding assistant.
//
// The tabby package allows Go applications to interact with TabbyAPI services, supporting:
// text completions, chat completions, embeddings generation, model management, LoRA adapter
// management, token operations, prompt templates, sampling parameters, health checks, and
// authentication.
//
// For complete documentation and examples, visit https://github.com/pixelsquared/go-tabbyapi
package tabby

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/pixelsquared/go-tabbyapi/internal/auth"
	"github.com/pixelsquared/go-tabbyapi/internal/rest"
)

// Client provides access to all TabbyAPI services and configuration options.
// Use NewClient to create a new client instance, then configure it with option
// methods like WithBaseURL or WithAPIKey.
type Client interface {
	// API Services

	// Completions returns the CompletionsService for generating text completions.
	Completions() CompletionsService

	// Chat returns the ChatService for handling chat-based interactions.
	Chat() ChatService

	// Embeddings returns the EmbeddingsService for generating vector embeddings.
	Embeddings() EmbeddingsService

	// Models returns the ModelsService for model management operations.
	Models() ModelsService

	// Lora returns the LoraService for managing LoRA adapters.
	Lora() LoraService

	// Templates returns the TemplatesService for managing prompt templates.
	Templates() TemplatesService

	// Tokens returns the TokensService for encoding and decoding tokens.
	Tokens() TokensService

	// Sampling returns the SamplingService for managing sampling parameters.
	Sampling() SamplingService

	// Health returns the HealthService for checking TabbyAPI health status.
	Health() HealthService

	// Auth returns the AuthService for managing authentication permissions.
	Auth() AuthService

	// Close releases resources used by the client.
	// Always call this method when you're done using the client.
	Close() error

	// Client configuration options

	// WithBaseURL sets the base URL for the TabbyAPI server.
	WithBaseURL(url string) Client

	// WithHTTPClient sets a custom HTTP client for making requests.
	WithHTTPClient(client *http.Client) Client

	// WithAPIKey sets the API key for authentication.
	WithAPIKey(key string) Client

	// WithAdminKey sets the admin key for administrative operations.
	WithAdminKey(key string) Client

	// WithBearerToken sets a bearer token for authentication.
	WithBearerToken(token string) Client

	// WithTimeout sets the timeout duration for API requests.
	WithTimeout(timeout time.Duration) Client

	// WithRetryPolicy sets the retry policy for failed requests.
	WithRetryPolicy(policy RetryPolicy) Client
}

// CompletionsService handles text completion requests for generating text
// based on provided prompts, with support for both synchronous and streaming responses.
type CompletionsService interface {
	// Create generates a completion for the provided request.
	// This method returns a complete, non-streaming response.
	//
	// To generate completions, pass a CompletionRequest with a prompt and desired parameters
	// such as max_tokens, temperature, top_p, etc.
	Create(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)

	// CreateStream generates a streaming completion where tokens are returned
	// incrementally as they're generated.
	//
	// The returned CompletionStream must be closed when no longer needed to release resources.
	// Use req.Stream = true when using this method.
	CreateStream(ctx context.Context, req *CompletionRequest) (CompletionStream, error)
}

// ChatService handles chat completion requests for multi-turn conversations
// with support for different message roles (system, user, assistant).
type ChatService interface {
	// Create generates a chat completion for the provided request.
	// This method returns a complete, non-streaming response.
	//
	// Pass a ChatCompletionRequest with an array of messages representing the conversation
	// history and desired parameters.
	Create(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error)

	// CreateStream generates a streaming chat completion where tokens are returned
	// incrementally as they're generated.
	//
	// The returned ChatCompletionStream must be closed when no longer needed to release resources.
	// Use req.Stream = true when using this method.
	CreateStream(ctx context.Context, req *ChatCompletionRequest) (ChatCompletionStream, error)
}

// ModelsService handles model management operations including listing, loading,
// unloading, and querying information about models.
type ModelsService interface {
	// List returns all available models.
	//
	// This method retrieves information about all models available to the TabbyAPI
	// server, including their IDs and parameters.
	List(ctx context.Context) (*ModelList, error)

	// Get returns the currently loaded model.
	//
	// This method provides information about the model that is currently loaded
	// and ready for use in the TabbyAPI server.
	Get(ctx context.Context) (*ModelCard, error)

	// Load loads a model with the specified parameters.
	//
	// This method initiates loading of a model into memory with the provided configuration
	// parameters. The operation is synchronous and returns when loading is complete.
	//
	// Parameters like max sequence length, RoPE scaling, cache size, and others can be
	// specified in the ModelLoadRequest.
	Load(ctx context.Context, req *ModelLoadRequest) (*ModelLoadResponse, error)

	// LoadStream loads a model and returns a stream of loading progress.
	//
	// This method is similar to Load but operates asynchronously, providing a stream
	// of updates during the loading process. This is useful for large models where
	// loading may take a significant amount of time.
	//
	// The returned ModelLoadStream must be closed when no longer needed.
	LoadStream(ctx context.Context, req *ModelLoadRequest) (ModelLoadStream, error)

	// Unload unloads the currently loaded model.
	//
	// This method releases memory and resources used by the currently loaded model.
	// After calling this method, no model will be available until another is loaded.
	Unload(ctx context.Context) error

	// GetProps returns model properties.
	//
	// This method provides additional properties of the currently loaded model,
	// such as context size, chat template, and default generation settings.
	GetProps(ctx context.Context) (*ModelPropsResponse, error)

	// Download downloads a model from HuggingFace.
	//
	// This method downloads a model from HuggingFace Model Hub to the local
	// TabbyAPI server. The model can then be loaded using the Load method.
	//
	// The DownloadRequest allows specifying the repository ID, revision,
	// authentication token, and file filters.
	Download(ctx context.Context, req *DownloadRequest) (*DownloadResponse, error)

	// ListDraft returns all available draft models.
	//
	// This method retrieves information about draft models available to the TabbyAPI
	// server. Draft models are models that are still in development or testing.
	ListDraft(ctx context.Context) (*ModelList, error)

	// ListEmbedding returns all available embedding models.
	//
	// This method retrieves information about embedding models available to the TabbyAPI
	// server. Embedding models are used to generate vector representations of text.
	ListEmbedding(ctx context.Context) (*ModelList, error)

	// GetEmbedding returns the currently loaded embedding model.
	//
	// This method provides information about the embedding model that is currently
	// loaded in the TabbyAPI server.
	GetEmbedding(ctx context.Context) (*ModelCard, error)

	// LoadEmbedding loads an embedding model.
	//
	// This method initiates loading of an embedding model into memory with the
	// provided configuration parameters. Embedding models are used to generate
	// vector representations of text through the EmbeddingsService.
	//
	// The EmbeddingModelLoadRequest allows specifying the model name and device.
	LoadEmbedding(ctx context.Context, req *EmbeddingModelLoadRequest) (*ModelLoadResponse, error)

	// UnloadEmbedding unloads the current embedding model.
	//
	// This method releases memory and resources used by the currently loaded
	// embedding model. After calling this method, no embedding model will be
	// available until another is loaded.
	UnloadEmbedding(ctx context.Context) error
}

// EmbeddingsService handles embedding generation for text inputs,
// converting text into numerical vector representations.
type EmbeddingsService interface {
	// Create generates embeddings for the provided input.
	//
	// This method converts text into vector embeddings that can be used for
	// semantic search, clustering, classification, and other machine learning tasks.
	//
	// The EmbeddingsRequest can contain either a single string or an array of
	// strings for batch processing. The response contains a corresponding
	// number of embedding vectors.
	//
	// The embedding model must be loaded via ModelsService.LoadEmbedding
	// before using this method, unless a default embedding model is configured.
	Create(ctx context.Context, req *EmbeddingsRequest) (*EmbeddingsResponse, error)
}

// LoraService handles Low-Rank Adaptation (LoRA) adapter management.
// LoRA adapters allow fine-tuning language models with significantly fewer parameters.
type LoraService interface {
	// List returns all available LoRA adapters.
	//
	// This method retrieves information about all LoRA adapters that are
	// available to the TabbyAPI server.
	List(ctx context.Context) (*LoraList, error)

	// GetActive returns the currently loaded LoRA adapters.
	//
	// This method provides information about the LoRA adapters that are
	// currently loaded and active in the TabbyAPI server.
	GetActive(ctx context.Context) (*LoraList, error)

	// Load loads specified LoRA adapters.
	//
	// This method loads one or more LoRA adapters into memory to modify the
	// behavior of the base model. The LoraLoadRequest allows specifying
	// multiple adapters with different scaling factors.
	//
	// LoRA adapters must be compatible with the currently loaded base model.
	Load(ctx context.Context, req *LoraLoadRequest) (*LoraLoadResponse, error)

	// Unload unloads all currently loaded LoRA adapters.
	//
	// This method removes all active LoRA adapters, returning the model
	// to its base behavior.
	Unload(ctx context.Context) error
}

// TokensService handles tokenization operations including encoding text to token IDs
// and decoding token IDs back to text.
type TokensService interface {
	// Encode converts text into token IDs.
	//
	// This method takes input text and converts it into the numerical token IDs
	// used by the model. This is useful for understanding how the model processes
	// text and for precisely controlling token usage.
	//
	// The input can be either a string or an array of ChatMessage objects for
	// encoding conversations.
	Encode(ctx context.Context, req *TokenEncodeRequest) (*TokenEncodeResponse, error)

	// Decode converts token IDs back into text.
	//
	// This method takes an array of token IDs and converts them back into the
	// corresponding text. This is useful for debugging or for processing
	// model outputs at the token level.
	Decode(ctx context.Context, req *TokenDecodeRequest) (*TokenDecodeResponse, error)
}

// TemplatesService handles prompt template management for different model types.
// Templates define how to format prompts for specific model architectures.
type TemplatesService interface {
	// List returns all available prompt templates.
	//
	// This method retrieves information about all prompt templates that are
	// available to the TabbyAPI server. Templates define how prompts are
	// formatted for different model architectures.
	List(ctx context.Context) (*TemplateList, error)

	// Switch changes the active prompt template.
	//
	// This method sets the specified template as the active one to be used
	// for formatting prompts. Different model architectures may require
	// different templates for optimal performance.
	//
	// The template name must be one of those returned by the List method.
	Switch(ctx context.Context, req *TemplateSwitchRequest) error

	// Unload unloads the currently selected template.
	//
	// This method removes the currently active template, reverting to the
	// model's default templating behavior.
	Unload(ctx context.Context) error
}

// SamplingService handles sampling parameter overrides and presets
// to control text generation behavior across different requests.
type SamplingService interface {
	// ListOverrides returns all sampler overrides and presets.
	//
	// This method retrieves information about available sampling parameter presets
	// and any currently active overrides. Sampling parameters control the randomness
	// and creativity of text generation (temperature, top_p, top_k, etc.).
	ListOverrides(ctx context.Context) (*SamplerOverrideListResponse, error)

	// SwitchOverride changes the active sampler override.
	//
	// This method allows switching to a predefined preset or setting custom
	// sampling parameters that will apply globally to all requests that don't
	// explicitly override them.
	//
	// Use the preset field to select a named preset, or the overrides map to
	// set specific parameter values.
	SwitchOverride(ctx context.Context, req *SamplerOverrideSwitchRequest) error

	// UnloadOverride unloads the currently selected override preset.
	//
	// This method removes any active sampling parameter overrides, reverting
	// to the default sampling behavior for each model.
	UnloadOverride(ctx context.Context) error
}

// HealthService handles health checks to verify the TabbyAPI server status.
type HealthService interface {
	// Check returns the current health status of the TabbyAPI server.
	//
	// This method can be used to verify that the server is running properly
	// and to detect any issues that might affect its operation.
	//
	// The response includes a status field indicating the overall health
	// and an optional list of issues if problems were detected.
	Check(ctx context.Context) (*HealthCheckResponse, error)
}

// AuthService handles authentication permissions and access levels.
type AuthService interface {
	// GetPermission returns the access level for the current authentication method.
	//
	// This method allows checking what permissions the current authentication
	// credentials have, which determines what API operations can be performed.
	//
	// Common permission levels include "none", "read", "write", and "admin".
	GetPermission(ctx context.Context) (*AuthPermissionResponse, error)
}

// NewClient creates a new TabbyAPI client with the provided options.
//
// By default, the client connects to http://localhost:8080 with a 30-second timeout
// and no authentication. Use the With* option functions to customize the client
// configuration, such as setting the API endpoint, authentication credentials,
// timeout, HTTP client, or retry policy.
func NewClient(options ...Option) Client {
	c := &clientImpl{
		baseURL:    "http://localhost:8080",
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

// clientImpl is the concrete implementation of the Client interface.
// It maintains the configuration and internal state for API requests.
type clientImpl struct {
	baseURL     string
	httpClient  *http.Client
	auth        Authenticator
	retryPolicy RetryPolicy
	restClient  *rest.Client
}

// Close releases resources used by the client
func (c *clientImpl) Close() error {
	// No resources to release yet
	return nil
}

// Implement WithX methods for clientImpl
func (c *clientImpl) WithBaseURL(url string) Client {
	c.baseURL = url
	return c
}

func (c *clientImpl) WithHTTPClient(client *http.Client) Client {
	c.httpClient = client
	return c
}

func (c *clientImpl) WithAPIKey(key string) Client {
	c.auth = &APIKeyAuthenticator{Key: key}
	return c
}

func (c *clientImpl) WithAdminKey(key string) Client {
	c.auth = &AdminKeyAuthenticator{Key: key}
	return c
}

func (c *clientImpl) WithBearerToken(token string) Client {
	c.auth = &BearerTokenAuthenticator{Token: token}
	return c
}

func (c *clientImpl) WithTimeout(timeout time.Duration) Client {
	c.httpClient.Timeout = timeout
	return c
}

func (c *clientImpl) WithRetryPolicy(policy RetryPolicy) Client {
	c.retryPolicy = policy
	return c
}

// getRestClient returns a REST client, initializing it if needed
func (c *clientImpl) getRestClient() *rest.Client {
	if c.restClient == nil {
		authProvider := auth.NoAuth
		if c.auth != nil {
			authProvider = c.auth
		}

		c.restClient = rest.New(
			c.baseURL,
			rest.WithHTTPClient(c.httpClient),
			rest.WithAuth(authProvider),
		)
	}
	return c.restClient
}

// buildURL constructs the URL for the API request
func (c *clientImpl) buildURL(endpoint string) string {
	endpoint = strings.TrimLeft(endpoint, "/")
	return fmt.Sprintf("%s/%s", c.baseURL, endpoint)
}

// Service getters
func (c *clientImpl) Completions() CompletionsService {
	return &completionsService{client: c.getRestClient(), baseURL: c.baseURL}
}

func (c *clientImpl) Chat() ChatService {
	return &chatService{client: c.getRestClient(), baseURL: c.baseURL}
}

func (c *clientImpl) Models() ModelsService {
	return &modelsService{client: c.getRestClient(), baseURL: c.baseURL}
}

func (c *clientImpl) Embeddings() EmbeddingsService {
	return &embeddingsService{client: c.getRestClient()}
}

func (c *clientImpl) Lora() LoraService {
	return &loraService{client: c.getRestClient()}
}

func (c *clientImpl) Templates() TemplatesService {
	return &templatesService{client: c.getRestClient()}
}

func (c *clientImpl) Tokens() TokensService {
	return &tokensService{client: c.getRestClient()}
}

func (c *clientImpl) Sampling() SamplingService {
	return &samplingService{client: c.getRestClient()}
}

func (c *clientImpl) Health() HealthService {
	return &healthService{client: c.getRestClient()}
}

func (c *clientImpl) Auth() AuthService {
	return &authService{client: c.getRestClient()}
}

// Internal stream implementation to avoid circular imports
// GenericStream implements a generic SSE stream
type GenericStream[T any] struct {
	ctx      context.Context
	cancel   context.CancelFunc
	response *http.Response
	reader   *bufio.Reader
	closed   bool
	mu       sync.Mutex
}

// newGenericStream creates a new stream for handling SSE responses
func newGenericStream[T any](ctx context.Context, resp *http.Response) *GenericStream[T] {
	ctx, cancel := context.WithCancel(ctx)
	return &GenericStream[T]{
		ctx:      ctx,
		cancel:   cancel,
		response: resp,
		reader:   bufio.NewReader(resp.Body),
	}
}

// Recv reads the next item from the stream
func (s *GenericStream[T]) Recv() (T, error) {
	var empty T
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return empty, ErrStreamClosed
	}

	// Check if context has been canceled
	select {
	case <-s.ctx.Done():
		return empty, s.ctx.Err()
	default:
	}

	// Read and parse the next event
	event, err := s.readEvent()
	if err != nil {
		return empty, err
	}

	// Parse the data
	var item T
	if err := json.Unmarshal([]byte(event.data), &item); err != nil {
		return empty, &StreamError{
			Message: "failed to unmarshal event data",
			Err:     err,
		}
	}

	return item, nil
}

// Close closes the stream and releases resources
func (s *GenericStream[T]) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true
	s.cancel()

	if s.response != nil && s.response.Body != nil {
		return s.response.Body.Close()
	}
	return nil
}

// sseEvent represents a Server-Sent Event
type sseEvent struct {
	id    string
	event string
	data  string
}

// readEvent reads a single SSE event from the response body
func (s *GenericStream[T]) readEvent() (*sseEvent, error) {
	var buffer bytes.Buffer

	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// If we've reached EOF but have some data, return what we have
				if buffer.Len() > 0 {
					return parseEvent(buffer.String())
				}
				return nil, io.EOF
			}
			return nil, &StreamError{
				Message: "error reading from stream",
				Err:     err,
			}
		}

		// Detect the end of an event (empty line)
		if line == "\n" || line == "\r\n" {
			if buffer.Len() > 0 {
				// We have a complete event
				return parseEvent(buffer.String())
			}
			// Empty event, continue reading
			continue
		}

		buffer.WriteString(line)
	}
}

// parseEvent parses a string into an SSE event
func parseEvent(eventStr string) (*sseEvent, error) {
	event := &sseEvent{}
	for _, line := range strings.Split(eventStr, "\n") {
		line = strings.TrimRight(line, "\r")
		if line == "" {
			continue
		}

		// Check for field prefixes
		if strings.HasPrefix(line, "id:") {
			event.id = strings.TrimSpace(line[3:])
		} else if strings.HasPrefix(line, "event:") {
			event.event = strings.TrimSpace(line[6:])
		} else if strings.HasPrefix(line, "data:") {
			if event.data != "" {
				event.data += "\n"
			}
			event.data += strings.TrimSpace(line[5:])
		} else if strings.HasPrefix(line, ":") {
			// Comment, ignore
		} else if strings.Contains(line, ":") {
			// Handle non-standard fields
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				field := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				if field == "id" {
					event.id = value
				} else if field == "event" {
					event.event = value
				} else if field == "data" {
					if event.data != "" {
						event.data += "\n"
					}
					event.data += value
				}
			}
		}
	}

	// If no event type is specified, it's a "message" event
	if event.event == "" {
		event.event = "message"
	}

	return event, nil
}

// Helper functions to create typed streams
func createCompletionStream(ctx context.Context, resp *http.Response) CompletionStream {
	return newGenericStream[*CompletionStreamResponse](ctx, resp)
}

func createChatCompletionStream(ctx context.Context, resp *http.Response) ChatCompletionStream {
	return newGenericStream[*ChatCompletionStreamResponse](ctx, resp)
}

func createModelLoadStream(ctx context.Context, resp *http.Response) ModelLoadStream {
	return newGenericStream[*ModelLoadResponse](ctx, resp)
}

// Service implementations

// completionsService implements the CompletionsService interface
type completionsService struct {
	client  *rest.Client
	baseURL string
}

func (s *completionsService) Create(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	// Force stream to false to ensure we get a regular response
	reqCopy := *req
	reqCopy.Stream = false

	// Create a response object
	var response CompletionResponse

	// Send the request to the completions endpoint
	err := s.client.Post(ctx, "v1/completions", &reqCopy, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create completion: %w", err)
	}

	return &response, nil
}

func (s *completionsService) CreateStream(ctx context.Context, req *CompletionRequest) (CompletionStream, error) {
	// Force stream to true to ensure we get a streaming response
	reqCopy := *req
	reqCopy.Stream = true

	// Construct the URL manually
	endpoint := "v1/completions"
	endpoint = strings.TrimLeft(endpoint, "/")
	url := fmt.Sprintf("%s/%s", s.baseURL, endpoint)

	// Create the request
	reqBody, err := json.Marshal(&reqCopy)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// Send the request
	resp, err := s.client.DoRaw(ctx, http.MethodPost, url, &reqCopy)
	if err != nil {
		return nil, fmt.Errorf("failed to create completion stream: %w", err)
	}

	// Create a stream from the response
	return createCompletionStream(ctx, resp), nil
}

// chatService implements the ChatService interface
type chatService struct {
	client  *rest.Client
	baseURL string
}

func (s *chatService) Create(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// Force stream to false to ensure we get a regular response
	reqCopy := *req
	reqCopy.Stream = false

	// Create a response object
	var response ChatCompletionResponse

	// Send the request to the chat completions endpoint
	err := s.client.Post(ctx, "v1/chat/completions", &reqCopy, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat completion: %w", err)
	}

	return &response, nil
}

func (s *chatService) CreateStream(ctx context.Context, req *ChatCompletionRequest) (ChatCompletionStream, error) {
	// Force stream to true to ensure we get a streaming response
	reqCopy := *req
	reqCopy.Stream = true

	// Construct the URL manually
	endpoint := "v1/chat/completions"
	endpoint = strings.TrimLeft(endpoint, "/")
	url := fmt.Sprintf("%s/%s", s.baseURL, endpoint)

	// Send the request
	resp, err := s.client.DoRaw(ctx, http.MethodPost, url, &reqCopy)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat completion stream: %w", err)
	}

	// Create a stream from the response
	return createChatCompletionStream(ctx, resp), nil
}

// embeddingsService implements the EmbeddingsService interface
type embeddingsService struct {
	client *rest.Client
}

func (s *embeddingsService) Create(ctx context.Context, req *EmbeddingsRequest) (*EmbeddingsResponse, error) {
	var response EmbeddingsResponse
	err := s.client.Post(ctx, "v1/embeddings", req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create embeddings: %w", err)
	}
	return &response, nil
}

// modelsService implements the ModelsService interface
type modelsService struct {
	client  *rest.Client
	baseURL string
}

func (s *modelsService) List(ctx context.Context) (*ModelList, error) {
	var response ModelList
	err := s.client.Get(ctx, "v1/models", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	return &response, nil
}

func (s *modelsService) Get(ctx context.Context) (*ModelCard, error) {
	var response ModelCard
	err := s.client.Get(ctx, "v1/models/current", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get current model: %w", err)
	}
	return &response, nil
}

func (s *modelsService) Load(ctx context.Context, req *ModelLoadRequest) (*ModelLoadResponse, error) {
	reqCopy := *req
	var response ModelLoadResponse
	err := s.client.Post(ctx, "v1/models/load", &reqCopy, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to load model: %w", err)
	}
	return &response, nil
}

func (s *modelsService) LoadStream(ctx context.Context, req *ModelLoadRequest) (ModelLoadStream, error) {
	reqCopy := *req

	// Construct the URL manually
	endpoint := "v1/models/load"
	endpoint = strings.TrimLeft(endpoint, "/")
	url := fmt.Sprintf("%s/%s", s.baseURL, endpoint)

	// Send the request
	resp, err := s.client.DoRaw(ctx, http.MethodPost, url, &reqCopy)
	if err != nil {
		return nil, fmt.Errorf("failed to load model stream: %w", err)
	}
	return createModelLoadStream(ctx, resp), nil
}

func (s *modelsService) Unload(ctx context.Context) error {
	err := s.client.Delete(ctx, "v1/models/current", nil)
	if err != nil {
		return fmt.Errorf("failed to unload model: %w", err)
	}
	return nil
}

func (s *modelsService) GetProps(ctx context.Context) (*ModelPropsResponse, error) {
	var response ModelPropsResponse
	err := s.client.Get(ctx, "v1/models/props", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get model properties: %w", err)
	}
	return &response, nil
}

func (s *modelsService) Download(ctx context.Context, req *DownloadRequest) (*DownloadResponse, error) {
	var response DownloadResponse
	err := s.client.Post(ctx, "v1/models/download", req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to download model: %w", err)
	}
	return &response, nil
}

func (s *modelsService) ListDraft(ctx context.Context) (*ModelList, error) {
	var response ModelList
	err := s.client.Get(ctx, "v1/models/draft", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list draft models: %w", err)
	}
	return &response, nil
}

func (s *modelsService) ListEmbedding(ctx context.Context) (*ModelList, error) {
	var response ModelList
	err := s.client.Get(ctx, "v1/models/embedding", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list embedding models: %w", err)
	}
	return &response, nil
}

func (s *modelsService) GetEmbedding(ctx context.Context) (*ModelCard, error) {
	var response ModelCard
	err := s.client.Get(ctx, "v1/models/embedding/current", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get current embedding model: %w", err)
	}
	return &response, nil
}

func (s *modelsService) LoadEmbedding(ctx context.Context, req *EmbeddingModelLoadRequest) (*ModelLoadResponse, error) {
	var response ModelLoadResponse
	err := s.client.Post(ctx, "v1/models/embedding/load", req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to load embedding model: %w", err)
	}
	return &response, nil
}

func (s *modelsService) UnloadEmbedding(ctx context.Context) error {
	err := s.client.Delete(ctx, "v1/models/embedding/current", nil)
	if err != nil {
		return fmt.Errorf("failed to unload embedding model: %w", err)
	}
	return nil
}

// loraService implements the LoraService interface
type loraService struct {
	client *rest.Client
}

func (s *loraService) List(ctx context.Context) (*LoraList, error) {
	var response LoraList
	err := s.client.Get(ctx, "v1/loras", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list LoRAs: %w", err)
	}
	return &response, nil
}

func (s *loraService) GetActive(ctx context.Context) (*LoraList, error) {
	var response LoraList
	err := s.client.Get(ctx, "v1/loras/active", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get active LoRAs: %w", err)
	}
	return &response, nil
}

func (s *loraService) Load(ctx context.Context, req *LoraLoadRequest) (*LoraLoadResponse, error) {
	var response LoraLoadResponse
	err := s.client.Post(ctx, "v1/loras/load", req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to load LoRAs: %w", err)
	}
	return &response, nil
}

func (s *loraService) Unload(ctx context.Context) error {
	err := s.client.Delete(ctx, "v1/loras/active", nil)
	if err != nil {
		return fmt.Errorf("failed to unload LoRAs: %w", err)
	}
	return nil
}

// tokensService implements the TokensService interface
type tokensService struct {
	client *rest.Client
}

func (s *tokensService) Encode(ctx context.Context, req *TokenEncodeRequest) (*TokenEncodeResponse, error) {
	var response TokenEncodeResponse
	err := s.client.Post(ctx, "v1/tokens/encode", req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to encode tokens: %w", err)
	}
	return &response, nil
}

func (s *tokensService) Decode(ctx context.Context, req *TokenDecodeRequest) (*TokenDecodeResponse, error) {
	var response TokenDecodeResponse
	err := s.client.Post(ctx, "v1/tokens/decode", req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode tokens: %w", err)
	}
	return &response, nil
}

// templatesService implements the TemplatesService interface
type templatesService struct {
	client *rest.Client
}

func (s *templatesService) List(ctx context.Context) (*TemplateList, error) {
	var response TemplateList
	err := s.client.Get(ctx, "v1/templates", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}
	return &response, nil
}

func (s *templatesService) Switch(ctx context.Context, req *TemplateSwitchRequest) error {
	err := s.client.Post(ctx, "v1/templates/switch", req, nil)
	if err != nil {
		return fmt.Errorf("failed to switch template: %w", err)
	}
	return nil
}

func (s *templatesService) Unload(ctx context.Context) error {
	err := s.client.Delete(ctx, "v1/templates/active", nil)
	if err != nil {
		return fmt.Errorf("failed to unload template: %w", err)
	}
	return nil
}

// samplingService implements the SamplingService interface
type samplingService struct {
	client *rest.Client
}

func (s *samplingService) ListOverrides(ctx context.Context) (*SamplerOverrideListResponse, error) {
	var response SamplerOverrideListResponse
	err := s.client.Get(ctx, "v1/sampler/overrides", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list sampler overrides: %w", err)
	}
	return &response, nil
}

func (s *samplingService) SwitchOverride(ctx context.Context, req *SamplerOverrideSwitchRequest) error {
	err := s.client.Post(ctx, "v1/sampler/overrides/switch", req, nil)
	if err != nil {
		return fmt.Errorf("failed to switch sampler override: %w", err)
	}
	return nil
}

func (s *samplingService) UnloadOverride(ctx context.Context) error {
	err := s.client.Delete(ctx, "v1/sampler/overrides/active", nil)
	if err != nil {
		return fmt.Errorf("failed to unload sampler override: %w", err)
	}
	return nil
}

// healthService implements the HealthService interface
type healthService struct {
	client *rest.Client
}

func (s *healthService) Check(ctx context.Context) (*HealthCheckResponse, error) {
	var response HealthCheckResponse
	err := s.client.Get(ctx, "v1/health", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to check health: %w", err)
	}
	return &response, nil
}

// authService implements the AuthService interface
type authService struct {
	client *rest.Client
}

func (s *authService) GetPermission(ctx context.Context) (*AuthPermissionResponse, error) {
	var response AuthPermissionResponse
	err := s.client.Get(ctx, "v1/auth/permission", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get authentication permissions: %w", err)
	}
	return &response, nil
}
