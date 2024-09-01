package nip90

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/openagentsinc/v3/relay/internal/common"
	"github.com/openagentsinc/v3/relay/internal/nostr"
)

func HandleAgentCommandRequest(conn *websocket.Conn, event *nostr.Event) {
	command := extractCommand(event)
	log.Printf("Received agent command request: %s", command)

	// Log all of the fields of the event, one per line
	logEventDetails(event)

	// TODO: Implement agent command routing logic here
	responseEvent := &nostr.Event{
		Kind:      6838, // Updated event kind for agent command response
		Content:   "Acknowledged. Will respond shortly.",
		CreatedAt: time.Now(),
		Tags:      [][]string{},
	}

	// Send the response back to the client
	response := common.CreateEventMessage(responseEvent)
	err := conn.WriteJSON(response)
	if err != nil {
		log.Println("Error writing agent command response to WebSocket:", err)
	}
}

func logEventDetails(event *nostr.Event) {
	var sb strings.Builder

	sb.WriteString("Event Details:\n")
	sb.WriteString(fmt.Sprintf("  ID: %s\n", event.ID))
	sb.WriteString(fmt.Sprintf("  PubKey: %s\n", event.PubKey))
	sb.WriteString(fmt.Sprintf("  CreatedAt: %s\n", event.CreatedAt))
	sb.WriteString(fmt.Sprintf("  Kind: %d\n", event.Kind))
	sb.WriteString("  Tags:\n")
	for _, tag := range event.Tags {
		sb.WriteString(fmt.Sprintf("    - %v\n", tag))
	}
	sb.WriteString(fmt.Sprintf("  Content: %s\n", event.Content))
	sb.WriteString(fmt.Sprintf("  Sig: %s\n", event.Sig))

	log.Print(sb.String())
}

func extractCommand(event *nostr.Event) string {
	for _, tag := range event.Tags {
		if len(tag) >= 3 && tag[0] == "i" {
			return tag[1]
		}
	}
	return ""
}
