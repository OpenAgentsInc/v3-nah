package nip90

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/openagentsinc/v3/relay/internal/common"
	"github.com/openagentsinc/v3/relay/internal/nostr"
)

func SendAgentCommandResponse(conn *websocket.Conn) {
	responseEvent := &nostr.Event{
		Kind:      6838, // Event kind for agent command response
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