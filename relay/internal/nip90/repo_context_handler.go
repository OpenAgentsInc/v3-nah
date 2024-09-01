package nip90

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"net/url"

	"github.com/openagentsinc/v3/relay/internal/github"
	"github.com/openagentsinc/v3/relay/internal/groq"
)

func GetRepoContext(repo string) string {
	log.Printf("GetRepoContext called for repo: %s", repo)

	owner, repoName := parseRepo(repo)
	if owner == "" || repoName == "" {
		return "Error: Invalid repository format. Expected 'owner/repo' or a valid GitHub URL."
	}

	context, err := analyzeRepository(owner, repoName)
	if err != nil {
		if err == github.ErrGitHubTokenNotSet {
			return fmt.Sprintf("Error: %v", err)
		}
		log.Printf("Error analyzing repository: %v", err)
		return fmt.Sprintf("Error analyzing repository: %v", err)
	}

	return summarizeContext(context)
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

func analyzeRepository(owner, repo string) (string, error) {
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
		{Role: "system", Content: "You are a repository analyzer. Analyze the repository structure and content using the provided tools."},
		{Role: "user", Content: fmt.Sprintf("Analyze the following repository structure and provide a summary:\n\n%s", rootContent)},
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
			result, err := executeToolCall(owner, repo, toolCall)
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

		messages = append(messages, response.Choices[0].Message)
	}

	return context.String(), nil
}

func executeToolCall(owner, repo string, toolCall groq.ToolCall) (string, error) {
	var args map[string]string
	err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling tool call arguments: %v", err)
	}

	switch toolCall.Function.Name {
	case "view_file":
		return github.ViewFile(owner, repo, args["path"], "")
	case "view_folder":
		return github.ViewFolder(owner, repo, args["path"], "")
	case "generate_summary":
		return generateSummary(args["content"])
	default:
		return "", fmt.Errorf("unknown tool: %s", toolCall.Function.Name)
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

func summarizeContext(context string) string {
	summary, err := generateSummary(context)
	if err != nil {
		log.Printf("Error summarizing context: %v", err)
		return "Error occurred while summarizing the context"
	}
	return summary
}