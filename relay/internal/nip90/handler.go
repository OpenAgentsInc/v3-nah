package nip90

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/openagentsinc/v3/relay/internal/nostr"
	"github.com/openagentsinc/v3/relay/internal/groq"
	"github.com/openagentsinc/v3/relay/internal/nip01"
)

type AudioData struct {
	Data   string
	Format string
}

func HandleAudioMessage(conn *websocket.Conn, audioData *AudioData) {
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
	response := nip01.CreateEventMessage(responseEvent)
	err = conn.WriteJSON(response)
	if err != nil {
		log.Println("Error writing audio response to WebSocket:", err)
	}
}

func HandleAgentCommandRequest(conn *websocket.Conn, event *nostr.Event) {
	log.Printf("Received agent command request: %s", event.Content)

	// TODO: Implement agent command routing logic here
	// For now, we'll just echo the command back as a response

	responseEvent := &nostr.Event{
		Kind:      6838, // Updated event kind for agent command response
		Content:   "Acknowledged. Will respond shortly.",
		CreatedAt: time.Now(),
		Tags:      [][]string{},
	}

	// Send the response back to the client
	response := nip01.CreateEventMessage(responseEvent)
	err := conn.WriteJSON(response)
	if err != nil {
		log.Println("Error writing agent command response to WebSocket:", err)
	}
}

func HandleNIP90Event(conn *websocket.Conn, event *nostr.Event) {
	switch event.Kind {
	case 5252:
		var audioData struct {
			Audio  string `json:"audio"`
			Format string `json:"format"`
		}
		err := json.Unmarshal([]byte(event.Content), &audioData)
		if err != nil {
			log.Printf("Error unmarshaling audio data: %v", err)
			return
		}
		HandleAudioMessage(conn, &AudioData{
			Data:   audioData.Audio,
			Format: audioData.Format,
		})
	case 5838:
		HandleAgentCommandRequest(conn, event)
	default:
		log.Printf("Unhandled NIP-90 event kind: %d", event.Kind)
	}
}