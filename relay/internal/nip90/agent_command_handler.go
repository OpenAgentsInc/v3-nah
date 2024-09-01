package nip90

import (
    "log"

    "github.com/gorilla/websocket"
    "github.com/openagentsinc/v3/relay/internal/nostr"
)

func HandleAgentCommandRequest(conn *websocket.Conn, event *nostr.Event) {
    // Log all of the fields of the event, one per line
    LogEventDetails(event)

    // Extract the repo parameter
    repo := extractRepoParam(event)
    if repo == "" {
        log.Println("Error: No repo parameter found in the event tags")
        SendAgentCommandResponse(conn, "Error: No repo parameter found")
        return
    }

    // Extract the user's prompt from the "i" tag
    prompt := extractPrompt(event)
    if prompt == "" {
        log.Println("Error: No prompt found in the event tags")
        SendAgentCommandResponse(conn, "Error: No prompt found")
        return
    }

    log.Printf("Received agent command request for repo: %s", repo)
    log.Printf("User prompt: %s", prompt)

    // Get repository context
    context := GetRepoContext(repo, conn, prompt)
    log.Printf("Repository context: %s", context)

    // Send the response back to the client
    SendAgentCommandResponse(conn, context)
}

func extractRepoParam(event *nostr.Event) string {
    for _, tag := range event.Tags {
        if len(tag) >= 3 && tag[0] == "param" && tag[1] == "repo" {
            return tag[2]
        }
    }
    return ""
}

func extractPrompt(event *nostr.Event) string {
    for _, tag := range event.Tags {
        if len(tag) >= 3 && tag[0] == "i" && tag[1] != "" {
            return tag[1]
        }
    }
    return ""
}