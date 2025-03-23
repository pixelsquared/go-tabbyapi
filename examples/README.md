# TabbyAPI Go Client Examples

This directory contains examples demonstrating how to use the Go client library for TabbyAPI.

## Overview

These examples showcase various features of the TabbyAPI Go client:

- **Completions**: Text generation examples
- **Chat**: Conversational AI examples
- **Embeddings**: Vector embedding generation
- **Models**: Model management and monitoring

## Prerequisites

Before running these examples, make sure you have:

1. A running TabbyAPI server (locally or remote)
2. Appropriate environment variables set:
   - `TABBY_API_ENDPOINT` (defaults to `http://localhost:8080`)
   - `TABBY_API_KEY` (if your server requires API key authentication)
   - `TABBY_ADMIN_KEY` (for examples that require administrative access)
   - `TABBY_MODEL_NAME` (for model loading examples)

## Running the Examples

Navigate to the specific example directory and run the Go program. For example:

```bash
cd completions/basic
go run main.go
```

## Example Categories

### Completions

- **Basic**: Simple text completion with a single prompt
- **Streaming**: Real-time streaming of generated text

### Chat

- **Basic**: Multi-turn conversation with the model
- **Streaming**: Real-time streaming of chat responses

### Embeddings

- **Basic**: Generate vector embeddings from text inputs

### Models

- **List**: List available models and currently loaded model
- **Streaming**: Stream model loading progress in real-time

## Additional Resources

For more information about the TabbyAPI Go client, see the main repository documentation.