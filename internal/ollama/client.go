// Package ollama talks to a local Ollama instance over its HTTP API.
// Default endpoint: http://localhost:11434
// Docs: https://github.com/ollama/ollama/blob/main/docs/api.md
package ollama

// TODO: implement Generate — send a prompt, return the model's response as a string.
// The endpoint is POST /api/generate with JSON body:
//   { "model": "<model-name>", "prompt": "<text>", "stream": false }
// Response field to extract: .response (string)
