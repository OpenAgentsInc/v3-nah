package nip01

import (
	"github.com/gorilla/websocket"
	"github.com/openagentsinc/v3/relay/internal/nostr"
	"log"
	"net/http"
	"sync"
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
	log.Printf("Handling message: %s", string(message))

	msg, err := ParseMessage(message)
	if err != nil {
		log.Println("Error parsing message:", err)
		return
	}

	log.Printf("Parsed message type: %s", msg.Type)

	switch msg.Type {
	case EventMessage:
		log.Println("Handling EventMessage")
		event, ok := msg.Data.(*nostr.Event)
		if !ok {
			log.Println("Error: EventMessage data is not of type *nostr.Event")
			return
		}
		r.handleEventMessage(conn, event)
	case ReqMessage:
		log.Println("Handling ReqMessage")
		filter, ok := msg.Data.(*nostr.Filter)
		if !ok {
			log.Println("Error: ReqMessage data is not of type *nostr.Filter")
			return
		}
		r.handleReqMessage(conn, filter)
	case CloseMessage:
		log.Println("Handling CloseMessage")
		subscriptionID, ok := msg.Data.(string)
		if !ok {
			log.Println("Error: CloseMessage data is not of type string")
			return
		}
		r.handleCloseMessage(conn, subscriptionID)
	case AudioMessage:
		log.Println("Handling AudioMessage")
		audioData, ok := msg.Data.(*AudioData)
		if !ok {
			log.Println("Error: AudioMessage data is not of type *AudioData")
			return
		}
		r.handleAudioMessage(conn, audioData)
	default:
		log.Println("Unknown message type:", msg.Type)
	}
}

func (r *Relay) handleEventMessage(conn *websocket.Conn, event *nostr.Event) {
	// TODO: Implement event validation and storage
	r.subscriptionManager.BroadcastEvent(event)
}

func (r *Relay) handleReqMessage(conn *websocket.Conn, filter *nostr.Filter) {
	// TODO: Implement subscription creation and event fetching
	// For now, just create a subscription
	sub := r.subscriptionManager.AddSubscription("temp_id", []*nostr.Filter{filter})
	go r.handleSubscription(conn, sub)
}

func (r *Relay) handleCloseMessage(conn *websocket.Conn, subscriptionID string) {
	r.subscriptionManager.RemoveSubscription(subscriptionID)
}

func (r *Relay) handleAudioMessage(conn *websocket.Conn, audioData *AudioData) {
	log.Printf("Received audio message. Format: %s, Length: %d\n", audioData.Format, len(audioData.Audio))

	// For now, just return a dummy transcription
	transcription := "This is a test transcription"
	response, err := CreateAudioResponseMessage(transcription)
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