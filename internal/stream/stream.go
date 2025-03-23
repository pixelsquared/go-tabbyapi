// Package stream provides functionality for handling Server-Sent Events (SSE) streams.
package stream

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

	"github.com/pixelsquared/go-tabbyapi/tabby"
)

// Stream is a generic interface for SSE streams.
type Stream[T any] struct {
	ctx      context.Context
	cancel   context.CancelFunc
	response *http.Response
	reader   *bufio.Reader
	closed   bool
	mu       sync.Mutex
}

// New creates a new Stream from an HTTP response.
func New[T any](ctx context.Context, resp *http.Response) *Stream[T] {
	ctx, cancel := context.WithCancel(ctx)
	return &Stream[T]{
		ctx:      ctx,
		cancel:   cancel,
		response: resp,
		reader:   bufio.NewReader(resp.Body),
	}
}

// Recv reads the next item from the stream.
func (s *Stream[T]) Recv() (T, error) {
	var empty T
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return empty, tabby.ErrStreamClosed
	}

	// Check if the context has been cancelled
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
		return empty, &tabby.StreamError{
			Message: "failed to unmarshal event data",
			Err:     err,
		}
	}

	return item, nil
}

// sseEvent represents a Server-Sent Event.
type sseEvent struct {
	id    string
	event string
	data  string
}

// readEvent reads a single SSE event from the response body.
func (s *Stream[T]) readEvent() (*sseEvent, error) {
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
			return nil, &tabby.StreamError{
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

// parseEvent parses a string into an SSE event.
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

// Close closes the stream and releases associated resources.
func (s *Stream[T]) Close() error {
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

// Done returns a channel that's closed when the stream is closed.
func (s *Stream[T]) Done() <-chan struct{} {
	return s.ctx.Done()
}

// IsClosed returns true if the stream has been closed.
func (s *Stream[T]) IsClosed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closed
}

// ForType is a convenience function that creates a stream for a specific type.
// This can be used to create concrete stream implementations for specific types.
func ForType[T any](ctx context.Context, resp *http.Response) tabby.Stream[T] {
	return New[T](ctx, resp)
}

// CreateCompletionStream creates a stream for completion responses.
func CreateCompletionStream(ctx context.Context, resp *http.Response) tabby.CompletionStream {
	return ForType[*tabby.CompletionStreamResponse](ctx, resp)
}

// CreateChatCompletionStream creates a stream for chat completion responses.
func CreateChatCompletionStream(ctx context.Context, resp *http.Response) tabby.ChatCompletionStream {
	return ForType[*tabby.ChatCompletionStreamResponse](ctx, resp)
}

// CreateModelLoadStream creates a stream for model loading responses.
func CreateModelLoadStream(ctx context.Context, resp *http.Response) tabby.ModelLoadStream {
	return ForType[*tabby.ModelLoadResponse](ctx, resp)
}

// Function to convert a response body to a stream.
func ReadStream[T any](ctx context.Context, resp *http.Response) (tabby.Stream[T], error) {
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, &tabby.APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("unexpected status code: %d, body: %s", resp.StatusCode, string(body)),
		}
	}

	stream := New[T](ctx, resp)
	return stream, nil
}
