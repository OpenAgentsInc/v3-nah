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

// ... (keep all other functions unchanged)

func analyzeRepository(owner, repo string, conn *websocket.Conn, prompt string) (string, error) {
	var context strings.Builder
	context.WriteString(fmt.Sprintf("Repository: https://github.com/%s/%s\n\n", owner, repo))

	rootContent, err := github.ViewFolder(owner, repo, "", "")
	if err != nil {
		return "", fmt.Errorf("error viewing root folder: %v", err)
	}

	// Look for authentication-related files
	authFiles := []string{
		"auth.go", "authentication.go", "login.go", "signin.go",
		"auth.js", "authentication.js", "login.js", "signin.js",
		"auth.ts", "authentication.ts", "login.ts", "signin.ts",
	}

	for _, file := range authFiles {
		content, err := github.ViewFile(owner, repo, file, "")
		if err == nil {
			context.WriteString(fmt.Sprintf("Found authentication-related file: %s\n", file))
			context.WriteString(fmt.Sprintf("Content:\n%s\n\n", content))
		}
	}

	// If no auth files found, check README for auth info
	if !strings.Contains(context.String(), "Found authentication-related file") {
		readmeContent, err := github.ViewFile(owner, repo, "README.md", "")
		if err == nil {
			context.WriteString("README.md content:\n")
			context.WriteString(readmeContent)
			context.WriteString("\n\n")
		}
	}

	return context.String(), nil
}

// ... (keep all other functions unchanged)