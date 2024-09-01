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

func summarizeContext(context, prompt string) string {
	messages := []groq.ChatMessage{
		{Role: "system", Content: "You are a helpful assistant that analyzes repository contexts. Provide specific and detailed answers focusing on the user's prompt. Always give a direct and comprehensive answer to the user's question, using information from the repository context."},
		{Role: "user", Content: fmt.Sprintf("Based on the following repository context, please provide a detailed and specific answer to the user's prompt: '%s'\n\nRepository context:\n%s", prompt, context)},
	}

	response, err := groq.ChatCompletionWithTools(messages, nil, nil)
	if err != nil {
		log.Printf("Error summarizing context: %v", err)
		return "Error occurred while analyzing the repository context"
	}

	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content
	}

	return "No specific information found related to the query"
}

// ... (keep all other functions unchanged)