package stream

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

// BenchmarkStream_Recv_SingleEvent benchmarks reading a single event from a stream.
func BenchmarkStream_Recv_SingleEvent(b *testing.B) {
	// Create a test event
	event := struct {
		Message string `json:"message"`
	}{
		Message: "test message",
	}

	// Marshal the event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		b.Fatalf("Failed to marshal event: %v", err)
	}

	// Create SSE data
	sseData := "data: " + string(eventJSON) + "\n\n"

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Create a fresh response for each iteration to reset the body reader
		resp := mockResponse(http.StatusOK, sseData)

		// Create a stream
		ctx := context.Background()
		stream := New[struct{ Message string }](ctx, resp)

		// Read from the stream
		_, err := stream.Recv()
		if err != nil {
			b.Fatalf("Recv returned an error: %v", err)
		}

		// Close the stream
		stream.Close()
	}
}

// BenchmarkStream_Recv_MultipleEvents benchmarks reading multiple events from a stream.
func BenchmarkStream_Recv_MultipleEvents(b *testing.B) {
	// Create test events
	events := []struct {
		Message string `json:"message"`
	}{
		{Message: "message 1"},
		{Message: "message 2"},
		{Message: "message 3"},
		{Message: "message 4"},
		{Message: "message 5"},
	}

	// Build SSE data
	var sseData bytes.Buffer
	for _, event := range events {
		eventJSON, err := json.Marshal(event)
		if err != nil {
			b.Fatalf("Failed to marshal event: %v", err)
		}
		sseData.WriteString("data: " + string(eventJSON) + "\n\n")
	}

	// Create a single SSE data string
	sseDataString := sseData.String()

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Create a fresh response for each iteration
		resp := mockResponse(http.StatusOK, sseDataString)

		// Create a stream
		ctx := context.Background()
		stream := New[struct{ Message string }](ctx, resp)

		// Read all events from the stream
		for j := 0; j < len(events); j++ {
			_, err := stream.Recv()
			if err != nil {
				b.Fatalf("Recv returned an error for event %d: %v", j, err)
			}
		}

		// Close the stream
		stream.Close()
	}
}

// BenchmarkStream_ParseEvent benchmarks the event parsing logic.
func BenchmarkStream_ParseEvent(b *testing.B) {
	// Create a complex SSE event
	eventStr := "id: 123\nevent: update\ndata: {\"message\":\"test\",\"count\":42}\n\n"

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		_, err := parseEvent(eventStr)
		if err != nil {
			b.Fatalf("parseEvent returned an error: %v", err)
		}
	}
}

// BenchmarkStream_ReadEvent benchmarks the event reading logic from a response body.
func BenchmarkStream_ReadEvent(b *testing.B) {
	// Create event data
	eventData := "data: {\"message\":\"test\"}\n\n" +
		"data: {\"message\":\"another\"}\n\n"

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Create a fresh response for each iteration
		resp := mockResponse(http.StatusOK, eventData)

		// Create a stream
		ctx := context.Background()
		stream := New[struct{ Message string }](ctx, resp)

		// Read one event
		_, err := stream.readEvent()
		if err != nil {
			b.Fatalf("readEvent returned an error: %v", err)
		}

		// Close the stream
		stream.Close()
	}
}

// BenchmarkStream_Unmarshal benchmarks unmarshaling event data to a struct.
func BenchmarkStream_Unmarshal(b *testing.B) {
	// Create a test event with multiple fields
	eventJSON := []byte(`{
		"id": "event-123",
		"type": "completion",
		"message": "This is a test message with more content to make it realistic",
		"timestamp": 1714142400,
		"tokens": 42,
		"model": "gpt-4",
		"finished": true,
		"metadata": {
			"source": "benchmark",
			"priority": "high"
		}
	}`)

	// Create a struct to unmarshal into
	type Metadata struct {
		Source   string `json:"source"`
		Priority string `json:"priority"`
	}

	type Event struct {
		ID        string   `json:"id"`
		Type      string   `json:"type"`
		Message   string   `json:"message"`
		Timestamp int64    `json:"timestamp"`
		Tokens    int      `json:"tokens"`
		Model     string   `json:"model"`
		Finished  bool     `json:"finished"`
		Metadata  Metadata `json:"metadata"`
	}

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		var event Event
		err := json.Unmarshal(eventJSON, &event)
		if err != nil {
			b.Fatalf("json.Unmarshal returned an error: %v", err)
		}
	}
}

// BenchmarkStream_LargePayload benchmarks processing a stream with large event payloads.
func BenchmarkStream_LargePayload(b *testing.B) {
	// Create a large payload (simulating a code completion response)
	largeData := struct {
		ID      string `json:"id"`
		Content string `json:"content"`
	}{
		ID:      "completion-123",
		Content: generateLargeString(5000), // 5KB of text
	}

	// Marshal to JSON
	largeJSON, err := json.Marshal(largeData)
	if err != nil {
		b.Fatalf("Failed to marshal large data: %v", err)
	}

	// Create SSE data
	sseData := "data: " + string(largeJSON) + "\n\n"

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Create a fresh response for each iteration
		resp := mockResponse(http.StatusOK, sseData)

		// Create a stream
		ctx := context.Background()
		stream := New[struct {
			ID      string `json:"id"`
			Content string `json:"content"`
		}](ctx, resp)

		// Read from the stream
		_, err := stream.Recv()
		if err != nil {
			b.Fatalf("Recv returned an error: %v", err)
		}

		// Close the stream
		stream.Close()
	}
}

// Helper function to generate a large string
func generateLargeString(size int) string {
	var sb bytes.Buffer
	sb.Grow(size)

	// Generate repeating pattern
	pattern := "This is a large completion response that simulates real-world code or text generation. "
	for sb.Len() < size {
		sb.WriteString(pattern)
	}

	return sb.String()[:size]
}
