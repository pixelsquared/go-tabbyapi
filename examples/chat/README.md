# TabbyAPI Go Client - Chat Examples

This directory contains examples demonstrating how to use the Go client library for TabbyAPI chat completions.

## Examples

### Basic Chat

Located in the `basic` directory, this example shows how to:
- Create and manage a multi-turn conversation with roles
- Use system, user, and assistant messages
- Process chat completion responses
- Continue a conversation by adding previous messages

```bash
cd basic
go run main.go
```

### Streaming Chat

Located in the `streaming` directory, this example demonstrates:
- How to create a streaming chat completion request
- Processing partial responses in real-time as they are generated
- Handling streaming chat deltas (content and role changes)
- Proper management of streaming resources

```bash
cd streaming
go run main.go
```

### JSON Schema Chat

Located in the `json_schema` directory, this example shows how to:
- Create a chat completion request with a JSON schema constraint
- Generate structured responses in valid JSON format
- Parse and validate the resulting JSON against the schema

```bash
cd json_schema
go run main.go
```

## Chat Message Roles

TabbyAPI chat completions support these message roles:

- `ChatMessageRoleSystem` - System instructions that guide the behavior of the assistant
- `ChatMessageRoleUser` - User messages in the conversation
- `ChatMessageRoleAssistant` - Assistant responses
- `ChatMessageRoleTool` - Tool usage messages (for function calling in supported models)

## Environment Variables

These examples use the following environment variables:

- `TABBY_API_ENDPOINT` - The URL of your TabbyAPI server (defaults to `http://localhost:8080`)
- `TABBY_API_KEY` - Your API key for authentication (if required by your server)

## Request Parameters

The examples demonstrate these common request parameters:

- `Messages` - Array of chat messages that form the conversation
- `MaxTokens` - Maximum number of tokens to generate
- `Temperature` - Controls randomness (0.0 = deterministic, higher values = more random)
- `TopP` - Top-p sampling (nucleus sampling)
- `TopK` - Top-k sampling
- `Stream` - Whether to stream the response
- `JSONSchema` - JSON schema for structured output generation

## Further Customization

You can extend these examples to:
- Implement a full conversational agent
- Add message history management
- Develop chat interfaces with specialized system prompts
- Integrate with other tools like retrieval systems