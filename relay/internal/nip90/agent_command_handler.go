package nip90

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/openagentsinc/v3/relay/internal/nostr"
)

func HandleAgentCommandRequest(conn *websocket.Conn, event *nostr.Event) {
	command := extractCommand(event)
	log.Printf("Received agent command request: %s", command)

	// Log all of the fields of the event, one per line
	LogEventDetails(event)

	// TODO: Implement agent command routing logic here

	// Send the response back to the client
	SendAgentCommandResponse(conn)
}

func extractCommand(event *nostr.Event) string {
	for _, tag := range event.Tags {
		if len(tag) >= 3 && tag[0] == "i" {
			return tag[1]
		}
	}
	return ""
}