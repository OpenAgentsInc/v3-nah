package nip90

import (
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
	// Check if the repo is a URL
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

	// If not a URL, expect the format "owner/repo"
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func analyzeRepository(owner, repo string) (string, error) {
	files := []string{
		"README.md",
		"go.mod",
		"main.go",
		"relay/internal/nip90/handler.go",
		"relay/internal/groq/tool_use.go",
	}

	var context strings.Builder
	context.WriteString(fmt.Sprintf("Repository: https://github.com/%s/%s\n\n", owner, repo))

	for _, file := range files {
		content, err := github.ViewFile(owner, repo, file, "")
		if err != nil {
			if err == github.ErrGitHubTokenNotSet {
				return "", err
			}
			log.Printf("Error viewing file %s: %v", file, err)
			continue
		}

		analysis, err := analyzeFileContent(file, content)
		if err != nil {
			log.Printf("Error analyzing file %s: %v", file, err)
			continue
		}

		context.WriteString(fmt.Sprintf("File: %s\n%s\n\n", file, analysis))
	}

	// Analyze the overall structure
	structure, err := github.ViewFolder(owner, repo, "", "")
	if err != nil {
		if err == github.ErrGitHubTokenNotSet {
			return "", err
		}
		log.Printf("Error viewing repository structure: %v", err)
	} else {
		structureAnalysis, err := analyzeRepoStructure(structure)
		if err != nil {
			log.Printf("Error analyzing repository structure: %v", err)
		} else {
			context.WriteString(fmt.Sprintf("Repository Structure:\n%s\n", structureAnalysis))
		}
	}

	return context.String(), nil
}

func analyzeFileContent(filename, content string) (string, error) {
	messages := []groq.ChatMessage{
		{Role: "system", Content: "You are a code analyst. Provide a brief summary of the given file content."},
		{Role: "user", Content: fmt.Sprintf("Analyze the following file (%s) content and provide a brief summary:\n\n%s", filename, content)},
	}

	response, err := groq.ChatCompletionWithTools(messages, nil, nil)
	if err != nil {
		return "", err
	}

	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no analysis generated")
}

func analyzeRepoStructure(structure string) (string, error) {
	messages := []groq.ChatMessage{
		{Role: "system", Content: "You are a repository structure analyst. Provide a brief summary of the given repository structure."},
		{Role: "user", Content: fmt.Sprintf("Analyze the following repository structure and provide a brief summary:\n\n%s", structure)},
	}

	response, err := groq.ChatCompletionWithTools(messages, nil, nil)
	if err != nil {
		return "", err
	}

	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no analysis generated")
}

func summarizeContext(context string) string {
	summary, err := SummarizeContext(context)
	if err != nil {
		log.Printf("Error summarizing context: %v", err)
		return "Error occurred while summarizing the context"
	}
	return summary
}

func SummarizeContext(context string) (string, error) {
	messages := []groq.ChatMessage{
		{Role: "system", Content: "You are a helpful assistant that summarizes repository contexts. Provide concise summaries in less than 200 words."},
		{Role: "user", Content: "Please summarize the following repository context in less than 200 words:\n\n" + context},
	}

	response, err := groq.ChatCompletionWithTools(messages, nil, nil)
	if err != nil {
		return "", err
	}

	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", nil
}