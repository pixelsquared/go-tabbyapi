# TabbyAPI Go Client - Completions Examples

This directory contains examples demonstrating how to use the Go client library for TabbyAPI text completions.

## Examples

### Basic Completion

Located in the `basic` directory, this example shows how to:
- Create a basic text completion request
- Process the completion response
- Extract generated text and metadata
- Handle token usage statistics

```bash
cd basic
go run main.go
```

### Streaming Completion

Located in the `streaming` directory, this example demonstrates:
- How to create a streaming text completion request
- How to process chunks of text as they are generated
- Proper handling of streaming responses and error conditions
- Collecting and combining stream chunks for further use

```bash
cd streaming
go run main.go
```

### JSON Schema Completion

Located in the `json_schema` directory, this example shows how to:
- Create a completion request with a JSON schema constraint
- Use the schema to enforce structured output format
- Parse and validate the resulting JSON response

```bash
cd json_schema
go run main.go
```

## Environment Variables

These examples use the following environment variables:

- `TABBY_API_ENDPOINT` - The URL of your TabbyAPI server (defaults to `http://localhost:8080`)
- `TABBY_API_KEY` - Your API key for authentication (if required by your server)

## Request Parameters

The examples demonstrate these common request parameters:

- `Prompt` - The text prompt to complete
- `MaxTokens` - Maximum number of tokens to generate
- `Temperature` - Controls randomness (0.0 = deterministic, higher values = more random)
- `TopP` - Top-p sampling (nucleus sampling)
- `TopK` - Top-k sampling
- `Stop` - Stop sequences that will end generation
- `Stream` - Whether to stream the response
- `JSONSchema` - JSON schema for structured output generation

## Further Customization

You can modify these examples to:
- Change the prompt text
- Adjust generation parameters for different creative outcomes
- Add retry logic for more robust error handling
- Save generated texts to files