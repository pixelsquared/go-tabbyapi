package stream

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/pixelsquared/go-tabbyapi/tabby"
)

// mockResponse creates a mock HTTP response with SSE data.
func mockResponse(statusCode int, data string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(data)),
		Header:     make(http.Header),
	}
}

// TestStream_Recv_SingleEvent tests reading a single event from a stream.
func TestStream_Recv_SingleEvent(t *testing.T) {
	// Create a test event
	event := struct {
		Message string `json:"message"`
	}{
		Message: "test message",
	}

	// Marshal the event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal event: %v", err)
	}

	// Create a mock SSE response
	sseData := "data: " + string(eventJSON) + "\n\n"
	resp := mockResponse(http.StatusOK, sseData)

	// Create a stream
	ctx := context.Background()
	stream := New[struct{ Message string }](ctx, resp)
	defer stream.Close()

	// Read from the stream
	result, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv returned an error: %v", err)
	}

	// Check the result
	if result.Message != event.Message {
		t.Errorf("Expected message %q, got %q", event.Message, result.Message)
	}

	// Reading again should return EOF since there are no more events
	_, err = stream.Recv()
	if err != io.EOF {
		t.Errorf("Expected EOF, got %v", err)
	}
}

// TestStream_Recv_MultipleEvents tests reading multiple events from a stream.
func TestStream_Recv_MultipleEvents(t *testing.T) {
	// Create test events
	events := []struct {
		Message string `json:"message"`
	}{
		{Message: "message 1"},
		{Message: "message 2"},
		{Message: "message 3"},
	}

	// Build SSE data
	var sseData bytes.Buffer
	for _, event := range events {
		eventJSON, err := json.Marshal(event)
		if err != nil {
			t.Fatalf("Failed to marshal event: %v", err)
		}
		sseData.WriteString("data: " + string(eventJSON) + "\n\n")
	}

	// Create a mock SSE response
	resp := mockResponse(http.StatusOK, sseData.String())

	// Create a stream
	ctx := context.Background()
	stream := New[struct{ Message string }](ctx, resp)
	defer stream.Close()

	// Read all events from the stream
	for i, expected := range events {
		result, err := stream.Recv()
		if err != nil {
			t.Fatalf("Recv returned an error for event %d: %v", i, err)
		}

		if result.Message != expected.Message {
			t.Errorf("Event %d: expected message %q, got %q", i, expected.Message, result.Message)
		}
	}

	// Reading again should return EOF since there are no more events
	_, err := stream.Recv()
	if err != io.EOF {
		t.Errorf("Expected EOF, got %v", err)
	}
}

// TestStream_Recv_ComplexEvents tests reading events with multiple fields.
func TestStream_Recv_ComplexEvents(t *testing.T) {
	// Create a test event with id and event type
	sseData := "id: 123\nevent: update\ndata: {\"message\":\"complex event\"}\n\n"
	resp := mockResponse(http.StatusOK, sseData)

	// Create a stream
	ctx := context.Background()
	stream := New[struct{ Message string }](ctx, resp)
	defer stream.Close()

	// Read from the stream
	result, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv returned an error: %v", err)
	}

	// Check the result
	if result.Message != "complex event" {
		t.Errorf("Expected message %q, got %q", "complex event", result.Message)
	}
}

// TestStream_Recv_MultilineData tests reading events with multiline data.
func TestStream_Recv_MultilineData(t *testing.T) {
	// Create a test event with multiline data
	sseData := "data: {\"line1\":\"first line\",\n" +
		"data: \"line2\":\"second line\"}\n\n"
	resp := mockResponse(http.StatusOK, sseData)

	// Create a stream
	ctx := context.Background()
	stream := New[struct {
		Line1 string `json:"line1"`
		Line2 string `json:"line2"`
	}](ctx, resp)
	defer stream.Close()

	// Read from the stream
	result, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv returned an error: %v", err)
	}

	// Check the result
	if result.Line1 != "first line" || result.Line2 != "second line" {
		t.Errorf("Expected lines %q and %q, got %q and %q", "first line", "second line", result.Line1, result.Line2)
	}
}

// TestStream_Recv_InvalidJSON tests handling of invalid JSON.
func TestStream_Recv_InvalidJSON(t *testing.T) {
	// Create a test event with invalid JSON
	sseData := "data: {invalid json}\n\n"
	resp := mockResponse(http.StatusOK, sseData)

	// Create a stream
	ctx := context.Background()
	stream := New[struct{ Message string }](ctx, resp)
	defer stream.Close()

	// Read from the stream
	_, err := stream.Recv()
	if err == nil {
		t.Fatal("Expected an error for invalid JSON, got nil")
	}

	// Check that it's a StreamError
	var streamErr *tabby.StreamError
	if !errors.As(err, &streamErr) {
		t.Errorf("Expected a StreamError, got %T: %v", err, err)
	}
}

// TestStream_Close tests closing a stream.
func TestStream_Close(t *testing.T) {
	// Create a long-running stream
	sseData := "data: {\"message\":\"test\"}\n\n" +
		"data: {\"message\":\"should not receive\"}\n\n"
	resp := mockResponse(http.StatusOK, sseData)

	// Create a stream
	ctx := context.Background()
	stream := New[struct{ Message string }](ctx, resp)

	// Read the first event
	result, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv returned an error: %v", err)
	}
	if result.Message != "test" {
		t.Errorf("Expected message %q, got %q", "test", result.Message)
	}

	// Close the stream
	err = stream.Close()
	if err != nil {
		t.Fatalf("Close returned an error: %v", err)
	}

	// Attempt to read after closing should fail
	_, err = stream.Recv()
	if err != tabby.ErrStreamClosed {
		t.Errorf("Expected ErrStreamClosed, got %v", err)
	}

	// Verify the stream is marked as closed
	if !stream.IsClosed() {
		t.Error("Expected IsClosed() to return true after closing")
	}

	// Closing again should be a no-op
	err = stream.Close()
	if err != nil {
		t.Errorf("Second Close returned an error: %v", err)
	}
}

// TestStream_ContextCancellation tests cancelling a stream via context.
func TestStream_ContextCancellation(t *testing.T) {
	// Create a long-running stream
	var longData bytes.Buffer
	for i := 0; i < 100; i++ {
		longData.WriteString(fmt.Sprintf("data: {\"index\":%d}\n\n", i))
	}
	resp := mockResponse(http.StatusOK, longData.String())

	// Create a stream with a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	stream := New[struct{ Index int }](ctx, resp)
	defer stream.Close()

	// Read a few events
	for i := 0; i < 5; i++ {
		result, err := stream.Recv()
		if err != nil {
			t.Fatalf("Recv returned an error for event %d: %v", i, err)
		}
		if result.Index != i {
			t.Errorf("Expected index %d, got %d", i, result.Index)
		}
	}

	// Cancel the context
	cancel()

	// Allow a little time for cancellation to take effect
	time.Sleep(10 * time.Millisecond)

	// Attempt to read after cancellation should fail
	_, err := stream.Recv()
	if err == nil || !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
}

// TestCreateTypeSpecificStreams tests the type-specific stream creator functions.
func TestCreateTypeSpecificStreams(t *testing.T) {
	// Create a simple SSE response
	sseData := "data: {\"id\":\"test\"}\n\n"
	resp := mockResponse(http.StatusOK, sseData)

	// Test CompletionStream
	ctx := context.Background()
	compStream := CreateCompletionStream(ctx, resp)
	defer compStream.Close()

	// Test ChatCompletionStream
	resp = mockResponse(http.StatusOK, sseData) // Need a fresh response since the previous one is consumed
	chatStream := CreateChatCompletionStream(ctx, resp)
	defer chatStream.Close()

	// Test ModelLoadStream
	resp = mockResponse(http.StatusOK, sseData)
	modelStream := CreateModelLoadStream(ctx, resp)
	defer modelStream.Close()
}

// TestReadStream tests the ReadStream function.
func TestReadStream(t *testing.T) {
	// Test successful case
	sseData := "data: {\"message\":\"test\"}\n\n"
	resp := mockResponse(http.StatusOK, sseData)

	ctx := context.Background()
	stream, err := ReadStream[struct{ Message string }](ctx, resp)
	if err != nil {
		t.Fatalf("ReadStream returned an error: %v", err)
	}
	defer stream.Close()

	// Read from the stream
	result, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv returned an error: %v", err)
	}
	if result.Message != "test" {
		t.Errorf("Expected message %q, got %q", "test", result.Message)
	}

	// Test error case - non-200 response
	errResp := mockResponse(http.StatusBadRequest, "error")
	_, err = ReadStream[struct{ Message string }](ctx, errResp)
	if err == nil {
		t.Fatal("Expected an error for non-200 response, got nil")
	}

	// Check that it's an APIError
	var apiErr *tabby.APIError
	if !errors.As(err, &apiErr) {
		t.Errorf("Expected an APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, apiErr.StatusCode)
	}
}
