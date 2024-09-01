package nip90

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/openagentsinc/v3/relay/internal/github"
	"github.com/openagentsinc/v3/relay/internal/groq"
	"github.com/openagentsinc/v3/relay/internal/nostr"
	"github.com/openagentsinc/v3/relay/internal/common"
)

func GetRepoContext(repo string, conn *websocket.Conn, prompt string) string {
	log.Printf("GetRepoContext called for repo: %s", repo)
	log.Printf("User prompt: %s", prompt)

	owner, repoName := parseRepo(repo)
	if owner == "" || repoName == "" {
		return "Error: Invalid repository format. Expected 'owner/repo' or a valid GitHub URL."
	}

	// Check if the prompt is a simple structural question
	if isSimpleStructuralQuestion(prompt) {
		return handleSimpleStructuralQuestion(owner, repoName, prompt, conn)
	}

	context, err := analyzeRepository(owner, repoName, conn, prompt)
	if err != nil {
		if err == github.ErrGitHubTokenNotSet {
			return fmt.Sprintf("Error: %v", err)
		}
		log.Printf("Error analyzing repository: %v", err)
		return fmt.Sprintf("Error analyzing repository: %v", err)
	}

	return summarizeContext(context, prompt)
}

func isSimpleStructuralQuestion(prompt string) bool {
	lowercasePrompt := strings.ToLower(prompt)
	return strings.Contains(lowercasePrompt, "what folders") ||
		strings.Contains(lowercasePrompt, "list folders") ||
		strings.Contains(lowercasePrompt, "show folders") ||
		strings.Contains(lowercasePrompt, "what directories") ||
		strings.Contains(lowercasePrompt, "list directories") ||
		strings.Contains(lowercasePrompt, "show directories")
}

func handleSimpleStructuralQuestion(owner, repo, prompt string, conn *websocket.Conn) string {
	rootContent, err := github.ViewFolder(owner, repo, "", "")
	if err != nil {
		return fmt.Sprintf("Error viewing root folder: %v", err)
	}

	folders := extractFolders(rootContent)
	response := fmt.Sprintf("The repository contains the following folders:\n\n%s", strings.Join(folders, "\n"))

	return response
}

func extractFolders(content string) []string {
	lines := strings.Split(content, "\n")
	var folders []string
	for _, line := range lines {
		if strings.HasSuffix(line, "(dir)") {
			folders = append(folders, strings.TrimSuffix(line, " (dir)"))
		}
	}
	return folders
}

func parseRepo(repo string) (string, string) {
	if strings.HasPrefix(repo, "http://") || strings.HasPrefix(repo, "https://") {
		parsedURL, err := url.Parse(repo)
		if err != nil {
			return "", ""
		}
		parts := strings.Split(parsedURL.Path, "/")
		if len(parts) < 3 {
			return "", ""
		}
		return parts[1], parts[2]
	}

	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func analyzeRepository(owner, repo string, conn *websocket.Conn, prompt string) (string, error) {
	var context strings.Builder
	context.WriteString(fmt.Sprintf("Repository: https://github.com/%s/%s\n\n", owner, repo))

	rootContent, err := github.ViewFolder(owner, repo, "", "")
	if err != nil {
		return "", fmt.Errorf("error viewing root folder: %v", err)
	}

	tools := []groq.Tool{
		{
			Type: "function",
			Function: groq.ToolFunction{
				Name:        "view_file",
				Description: "View the contents of a file in the repository",
				Parameters: groq.Parameters{
					Type: "object",
					Properties: map[string]groq.Property{
						"path": {Type: "string", Description: "The path of the file to view"},
					},
					Required: []string{"path"},
				},
			},
		},
		{
			Type: "function",
			Function: groq.ToolFunction{
				Name:        "view_folder",
				Description: "View the contents of a folder in the repository",
				Parameters: groq.Parameters{
					Type: "object",
					Properties: map[string]groq.Property{
						"path": {Type: "string", Description: "The path of the folder to view"},
					},
					Required: []string{"path"},
				},
			},
		},
		{
			Type: "function",
			Function: groq.ToolFunction{
				Name:        "generate_summary",
				Description: "Generate a summary of the given content",
				Parameters: groq.Parameters{
					Type: "object",
					Properties: map[string]groq.Property{
						"content": {Type: "string", Description: "The content to summarize"},
					},
					Required: []string{"content"},
				},
			},
		},
	}

	messages := []groq.ChatMessage{
		{Role: "system", Content: "You are a repository analyzer. Analyze the repository structure and content using the provided tools. Focus on the user's prompt and find relevant information. Always provide a direct and detailed answer to the user's question."},
		{Role: "user", Content: fmt.Sprintf("Analyze the following repository structure and provide a detailed summary, focusing on answering the user's prompt: '%s'\n\nRepository structure:\n%s", prompt, rootContent)},
	}

	for i := 0; i < 5; i++ { // Limit to 5 iterations to prevent infinite loops
		response, err := groq.ChatCompletionWithTools(messages, tools, nil)
		if err != nil {
			return "", fmt.Errorf("error in ChatCompletionWithTools: %v", err)
		}

		if len(response.Choices) == 0 || len(response.Choices[0].Message.ToolCalls) == 0 {
			break
		}

		for _, toolCall := range response.Choices[0].Message.ToolCalls {
			result, err := executeToolCall(owner, repo, toolCall, conn)
			if err != nil {
				log.Printf("Error executing tool call: %v", err)
				continue
			}
			messages = append(messages, groq.ChatMessage{
				Role:    "function",
				Content: result,
			})
			context.WriteString(fmt.Sprintf("%s:\n%s\n\n", toolCall.Function.Name, result))
		}

		messages = append(messages, groq.ChatMessage{
			Role:    response.Choices[0].Message.Role,
			Content: response.Choices[0].Message.Content,
		})
	}

	return context.String(), nil
}

func executeToolCall(owner, repo string, toolCall groq.ToolCall, conn *websocket.Conn) (string, error) {
	var args map[string]string
	err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling tool call arguments: %v", err)
	}

	switch toolCall.Function.Name {
	case "view_file":
		content, err := github.ViewFile(owner, repo, args["path"], "")
		if err != nil {
			return "", err
		}
		sendViewedFileEvent(conn, args["path"])
		return content, nil
	case "view_folder":
		return github.ViewFolder(owner, repo, args["path"], "")
	case "generate_summary":
		return generateSummary(args["content"])
	default:
		return "", fmt.Errorf("unknown tool: %s", toolCall.Function.Name)
	}
}

func sendViewedFileEvent(conn *websocket.Conn, path string) {
	if conn == nil {
		log.Println("WebSocket connection is not set")
		return
	}

	viewedEvent := &nostr.Event{
		Kind:      6838,
		Content:   fmt.Sprintf("Viewed %s", path),
		CreatedAt: time.Now(),
		Tags:      [][]string{},
	}

	response := common.CreateEventMessage(viewedEvent)
	err := conn.WriteJSON(response)
	if err != nil {
		log.Printf("Error writing viewed file event to WebSocket: %v", err)
	}
}

func generateSummary(content string) (string, error) {
	messages := []groq.ChatMessage{
		{Role: "system", Content: "You are a helpful assistant that summarizes content. Provide concise summaries."},
		{Role: "user", Content: "Please summarize the following content:\n\n" + content},
	}

	response, err := groq.ChatCompletionWithTools(messages, nil, nil)
	if err != nil {
		return "", err
	}

	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no summary generated")
}

func summarizeContext(context, prompt string) string {
	messages := []groq.ChatMessage{
		{Role: "system", Content: "You are a helpful assistant that analyzes repository contexts. Provide specific and detailed answers focusing on the user's prompt. Always give a direct and comprehensive answer to the user's question, using information from the repository context. Limit your response to approximately 75 words."},
		{Role: "user", Content: fmt.Sprintf("Based on the following repository context, please provide a detailed and specific answer to the user's prompt in about 75 words: '%s'\n\nRepository context:\n%s", prompt, context)},
	}

	response, err := groq.ChatCompletionWithTools(messages, nil, nil)
	if err != nil {
		log.Printf("Error summarizing context: %v", err)
		return "Error occurred while analyzing the repository context"
	}

	if len(response.Choices) > 0 {
		return limitWords(response.Choices[0].Message.Content, 75)
	}

	return "No specific information found related to the query"
}

func limitWords(s string, maxWords int) string {
	words := strings.Fields(s)
	if len(words) <= maxWords {
		return s
	}
	return strings.Join(words[:maxWords], " ") + "..."
}