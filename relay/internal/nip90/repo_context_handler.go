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
        {Role: "system", Content: "You are a helpful assistant that summarizes repository contexts. Provide concise summaries focusing on the user's prompt. Always give a direct answer to the user's question."},
        {Role: "user", Content: fmt.Sprintf("Please summarize the following repository context, focusing on answering the user's prompt: '%s'\n\n%s", prompt, context)},
    }

    response, err := groq.ChatCompletionWithTools(messages, nil, nil)
    if err != nil {
        log.Printf("Error summarizing context: %v", err)
        return "Error occurred while summarizing the context"
    }

    if len(response.Choices) > 0 {
        return fmt.Sprintf("User prompt: %s\n\nAnswer: %s", prompt, response.Choices[0].Message.Content)
    }

    return "No summary generated"
}

// ... (keep all other functions unchanged)