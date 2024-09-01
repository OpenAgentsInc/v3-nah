package nip90

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/openagentsinc/v3/relay/internal/common"
	"github.com/openagentsinc/v3/relay/internal/nostr"
)

func HandleAgentCommandRequest(conn *websocket.Conn, event *nostr.Event) {
	command := extractCommand(event)
	log.Printf("Received agent command request: %s", command)

	// Log all of the fields of the event, one per line
	LogEventDetails(event)

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

func extractCommand(event *nostr.Event) string {
	for _, tag := range event.Tags {
		if len(tag) >= 3 && tag[0] == "i" {
			return tag[1]
		}
	}
	return ""
}