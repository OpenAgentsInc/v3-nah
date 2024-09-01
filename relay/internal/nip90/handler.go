package nip90

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/openagentsinc/v3/relay/internal/nostr"
	"github.com/openagentsinc/v3/relay/internal/groq"
	"github.com/openagentsinc/v3/relay/internal/common"
)

type AudioData struct {
	Data   string
	Format string
}

func HandleAudioMessage(conn *websocket.Conn, event *nostr.Event) {
	audioData := extractAudioData(event)
	log.Printf("Received audio message. Format: %s, Length: %d\n", audioData.Format, len(audioData.Data))

	// Transcribe the audio using Groq API
	transcription, err := groq.TranscribeAudio(audioData.Data, audioData.Format)
	if err != nil {
		log.Printf("Error transcribing audio: %v", err)
		transcription = "Error transcribing audio"
	}

	// Create a response event
	responseEvent := &nostr.Event{
		Kind:      6252, // Updated event kind for transcription response
		Content:   transcription,
		CreatedAt: time.Now(),
		Tags:      [][]string{},
	}

	// Send the response back to the client
	response := common.CreateEventMessage(responseEvent)
	err = conn.WriteJSON(response)
	if err != nil {
		log.Println("Error writing audio response to WebSocket:", err)
	}
}

func HandleNIP90Event(conn *websocket.Conn, event *nostr.Event) {
	switch event.Kind {
	case 5252:
		HandleAudioMessage(conn, event)
	case 5838:
		HandleAgentCommandRequest(conn, event)
	default:
		log.Printf("Unhandled NIP-90 event kind: %d", event.Kind)
	}
}

func extractAudioData(event *nostr.Event) *AudioData {
	var audioData AudioData
	for _, tag := range event.Tags {
		if len(tag) >= 2 {
			switch tag[0] {
			case "i":
				audioData.Data = tag[1]
			case "param":
				if len(tag) >= 3 && tag[1] == "format" {
					audioData.Format = tag[2]
				}
			}
		}
	}
	return &audioData
}