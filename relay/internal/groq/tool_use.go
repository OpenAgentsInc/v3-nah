package groq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const GroqChatCompletionURL = "https://api.groq.com/openai/v1/chat/completions"

type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Tools       []Tool        `json:"tools,omitempty"`
	ToolChoice  interface{}   `json:"tool_choice,omitempty"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Tool struct {
	Type     string        `json:"type"`
	Function ToolFunction  `json:"function"`
}

type ToolFunction struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Parameters  Parameters `json:"parameters"`
}

type Parameters struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type ChatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Role       string     `json:"role"`
			Content    string     `json:"content"`
			ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
		} `json:"message"`
	} `json:"choices"`
}

type ToolCall struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"`
	Function ToolCallFunction `json:"function"`
}

type ToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

func ChatCompletionWithTools(messages []ChatMessage, tools []Tool, toolChoice interface{}) (*ChatCompletionResponse, error) {
	request := ChatCompletionRequest{
		Model:       "llama3-groq-70b-8192-tool-use-preview", // Using the recommended model for tool use
		Messages:    messages,
		Tools:       tools,
		ToolChoice:  toolChoice,
		Temperature: 0.7,
		MaxTokens:   4096,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", GroqChatCompletionURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("GROQ_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var result ChatCompletionResponse
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &result, nil
}

// TODO: Implement functions to handle tool calls and their results
// func HandleToolCalls(toolCalls []ToolCall) ([]ChatMessage, error) {
//     // Implement logic to execute tool calls and format their results
// }

// TODO: Implement a function to summarize context using Groq API
// func SummarizeContext(context string) (string, error) {
//     // Implement logic to summarize context using Groq API
// }

// TODO: Implement a function for the indexing process
// func IndexRepository(repo string) (string, error) {
//     // Implement logic to index a repository and generate context
// }