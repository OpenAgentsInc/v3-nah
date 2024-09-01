package nip01

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/openagentsinc/v3/relay/internal/nostr"
)

type Relay struct {
	upgrader            websocket.Upgrader
	subscriptionManager *SubscriptionManager
	mu                  sync.Mutex
}

func NewRelay() *Relay {
	return &Relay{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for now
			},
		},
		subscriptionManager: NewSubscriptionManager(),
	}
}

func (r *Relay) HandleWebSocket(w http.ResponseWriter, req *http.Request) {
	conn, err := r.upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		r.handleMessage(conn, message)
	}
}

func (r *Relay) handleMessage(conn *websocket.Conn, message []byte) {
	msg, err := ParseMessage(message)
	if err != nil {
		log.Println("Error parsing message:", err)
		return
	}

	switch msg.Type {
	case EventMessage:
		event, ok := msg.Data.(*nostr.Event)
		if !ok {
			log.Println("Error: EventMessage data is not of type *nostr.Event")
			return
		}
		r.handleEventMessage(conn, event)
	case ReqMessage:
		reqMsg, ok := msg.Data.(*nostr.ReqMessage)
		if !ok {
			log.Println("Error: ReqMessage data is not of type *nostr.ReqMessage")
			return
		}
		r.handleReqMessage(conn, reqMsg)
	case CloseMessage:
		subscriptionID, ok := msg.Data.(string)
		if !ok {
			log.Println("Error: CloseMessage data is not of type string")
			return
		}
		r.handleCloseMessage(conn, subscriptionID)
	default:
		log.Println("Unknown message type:", msg.Type)
	}
}

func (r *Relay) handleEventMessage(conn *websocket.Conn, event *nostr.Event) {
	log.Printf("Handling event with kind: %d", event.Kind)

	if event.Kind == 5252 {
		var audioData AudioData
		err := json.Unmarshal([]byte(event.Content), &audioData)
		if err != nil {
			log.Printf("Error unmarshaling audio data: %v", err)
			return
		}
		r.handleAudioMessage(conn, &audioData)
	} else {
		// Handle other event types or broadcast to subscribers
		r.subscriptionManager.BroadcastEvent(event)
	}
}

func (r *Relay) handleReqMessage(conn *websocket.Conn, reqMsg *nostr.ReqMessage) {
	// TODO: Implement subscription creation and event fetching
	// For now, just create a subscription
	sub := r.subscriptionManager.AddSubscription(reqMsg.SubscriptionID, []*nostr.Filter{&reqMsg.Filter})
	go r.handleSubscription(conn, sub)
}

func (r *Relay) handleCloseMessage(conn *websocket.Conn, subscriptionID string) {
	r.subscriptionManager.RemoveSubscription(subscriptionID)
}

func (r *Relay) handleAudioMessage(conn *websocket.Conn, audioData *AudioData) {
	log.Printf("Received audio message. Format: %s, Length: %d\n", audioData.Format, len(audioData.Audio))

	// Transcribe the audio using Groq API
	transcription, err := TranscribeAudio(audioData.Audio, audioData.Format)
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
	response, err := CreateEventMessage(responseEvent)
	if err != nil {
		log.Println("Error creating audio response message:", err)
		return
	}

	err = conn.WriteJSON(response)
	if err != nil {
		log.Println("Error writing audio response to WebSocket:", err)
	}
}

func (r *Relay) handleSubscription(conn *websocket.Conn, sub *Subscription) {
	for event := range sub.Events {
		msg, err := CreateEventMessage(event)
		if err != nil {
			log.Println("Error creating event message:", err)
			continue
		}

		err = conn.WriteJSON(msg)
		if err != nil {
			log.Println("Error writing event to WebSocket:", err)
			break
		}
	}
}

func (r *Relay) Start(addr string) error {
	http.HandleFunc("/", r.HandleWebSocket)
	return http.ListenAndServe(addr, nil)
}
