# TabbyAPI Go Client - Embeddings Examples

This directory contains examples demonstrating how to use the Go client library for generating vector embeddings with TabbyAPI.

## Examples

### Basic Embeddings

Located in the `basic` directory, this example shows how to:
- Generate vector embeddings from text input
- Process embedding responses
- Handle different embedding formats (float arrays or base64)
- Work with both single inputs and batched inputs

```bash
cd basic
go run main.go
```

## Embedding Use Cases

Vector embeddings are useful for many applications, including:

1. **Semantic Search**: Find documents or passages with similar meanings
2. **Clustering**: Group texts by semantic similarity
3. **Classification**: Train ML models on these fixed-dimension representations
4. **Recommendation Systems**: Find similar items based on embedding proximity
5. **Information Retrieval**: Enhance search results with semantic understanding

## Environment Variables

This example uses the following environment variables:

- `TABBY_API_ENDPOINT` - The URL of your TabbyAPI server (defaults to `http://localhost:8080`)
- `TABBY_API_KEY` - Your API key for authentication (if required by your server)

## Request Parameters

The example demonstrates these request parameters:

- `Input` - Text to embed (can be a single string or string array)
- `Model` - Optional embedding model name to use
- `EncodingFormat` - Format for the output embeddings (float or base64)

## Embedding Dimensions

The dimension of the embedding vectors depends on the model loaded on your TabbyAPI server. Common dimensions include:

- Small embedding models: 384 dimensions
- Medium embedding models: 768 dimensions
- Large embedding models: 1024 or 1536 dimensions

## Processing Embeddings

Once you have generated embeddings, you would typically:

1. Store them in a vector database like Pinecone, Weaviate, or Milvus
2. Perform similarity searches using cosine similarity, dot product, or Euclidean distance
3. Use them as input features for machine learning models

## Further Customization

You can modify this example to:
- Generate embeddings for entire document collections
- Implement semantic search functionality
- Create clustering algorithms based on embedding similarity
- Compare different similarity metrics