package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// BenchmarkClient_Get benchmarks the Get method of the Client.
func BenchmarkClient_Get(b *testing.B) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()

	// Create client
	client := New(server.URL)

	// Prepare result variable for reuse
	var result testResponse

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		err := client.Get(context.Background(), "/test", nil, &result)
		if err != nil {
			b.Fatalf("Get returned an error: %v", err)
		}
	}
}

// BenchmarkClient_Post benchmarks the Post method of the Client.
func BenchmarkClient_Post(b *testing.B) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"created"}`))
	}))
	defer server.Close()

	// Create client
	client := New(server.URL)

	// Create request body
	body := map[string]string{"key": "value"}

	// Prepare result variable for reuse
	var result testResponse

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		err := client.Post(context.Background(), "/test", body, &result)
		if err != nil {
			b.Fatalf("Post returned an error: %v", err)
		}
	}
}

// BenchmarkClient_DoRaw benchmarks the DoRaw method of the Client.
func BenchmarkClient_DoRaw(b *testing.B) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()

	// Create client
	client := New(server.URL)

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		resp, err := client.DoRaw(context.Background(), http.MethodGet, server.URL+"/test", nil)
		if err != nil {
			b.Fatalf("DoRaw returned an error: %v", err)
		}
		resp.Body.Close()
	}
}

// BenchmarkJSON_Marshal benchmarks JSON marshaling of a request body.
func BenchmarkJSON_Marshal(b *testing.B) {
	// Create a complex request body
	body := map[string]interface{}{
		"model":       "gpt-4",
		"temperature": 0.7,
		"max_tokens":  100,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a helpful assistant.",
			},
			{
				"role":    "user",
				"content": "Tell me about Go programming language.",
			},
		},
	}

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(body)
		if err != nil {
			b.Fatalf("json.Marshal returned an error: %v", err)
		}
	}
}

// BenchmarkJSON_Unmarshal benchmarks JSON unmarshaling of a response body.
func BenchmarkJSON_Unmarshal(b *testing.B) {
	// Create a complex response JSON
	responseJSON := []byte(`{
		"id": "chat-123",
		"object": "chat.completion",
		"created": 1677858242,
		"model": "gpt-4",
		"choices": [
			{
				"message": {
					"role": "assistant",
					"content": "Go is a statically typed, compiled programming language designed at Google."
				},
				"finish_reason": "stop",
				"index": 0
			}
		],
		"usage": {
			"prompt_tokens": 25,
			"completion_tokens": 20,
			"total_tokens": 45
		}
	}`)

	// Create a struct to unmarshal into
	type Choice struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	}

	type Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	}

	type Response struct {
		ID      string   `json:"id"`
		Object  string   `json:"object"`
		Created int64    `json:"created"`
		Model   string   `json:"model"`
		Choices []Choice `json:"choices"`
		Usage   Usage    `json:"usage"`
	}

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		var response Response
		err := json.Unmarshal(responseJSON, &response)
		if err != nil {
			b.Fatalf("json.Unmarshal returned an error: %v", err)
		}
	}
}
