package github

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	githubAPIBaseURL = "https://api.github.com"
)

type GitHubFile struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

func ViewFile(owner, repo, path, branch string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", githubAPIBaseURL, owner, repo, path)
	if branch != "" {
		url += fmt.Sprintf("?ref=%s", branch)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return "", fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var file GitHubFile
	err = json.Unmarshal(body, &file)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %v", err)
	}

	if file.Encoding != "base64" {
		return "", fmt.Errorf("unexpected file encoding: %s", file.Encoding)
	}

	decodedContent, err := base64.StdEncoding.DecodeString(file.Content)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 content: %v", err)
	}

	return string(decodedContent), nil
}