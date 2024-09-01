package nip01

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/openagentsinc/v3/relay/internal/nostr"
	"github.com/openagentsinc/v3/relay/internal/nip90"
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
		r.handleReqMessage(conn, msg)
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

	switch {
	case event.Kind == 5252 || event.Kind == 5838:
		nip90.HandleNIP90Event(conn, event)
	default:
		// Handle other event types or broadcast to subscribers
		r.subscriptionManager.BroadcastEvent(event)
	}
}

func (r *Relay) handleReqMessage(conn *websocket.Conn, msg *Message) {
	log.Printf("Handling REQ message: %+v", msg)

	reqData, ok := msg.Data.([]interface{})
	if !ok {
		log.Println("Error: REQ message data is not of type []interface{}")
		return
	}

	if len(reqData) < 2 {
		log.Println("Invalid REQ message format")
		return
	}

	subscriptionID, ok := reqData[0].(string)
	if !ok {
		log.Println("Invalid subscription ID in REQ message")
		return
	}

	filters := make([]*nostr.Filter, 0)
	for _, filterData := range reqData[1:] {
		filterMap, ok := filterData.(map[string]interface{})
		if !ok {
			log.Println("Invalid filter format in REQ message")
			continue
		}
		filter := &nostr.Filter{}
		// Parse filter data and populate the filter object
		// This is a simplified version, you may need to add more fields
		if ids, ok := filterMap["ids"].([]interface{}); ok {
			filter.IDs = make([]string, len(ids))
			for i, id := range ids {
				filter.IDs[i], _ = id.(string)
			}
		}
		filters = append(filters, filter)
	}

	sub := r.subscriptionManager.AddSubscription(subscriptionID, filters)
	go r.handleSubscription(conn, sub)
}

func (r *Relay) handleCloseMessage(conn *websocket.Conn, subscriptionID string) {
	r.subscriptionManager.RemoveSubscription(subscriptionID)
}

func (r *Relay) handleSubscription(conn *websocket.Conn, sub *Subscription) {
	for event := range sub.Events {
		msg := CreateEventMessage(event)
		err := conn.WriteJSON(msg)
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

func CreateEventMessage(event *nostr.Event) []interface{} {
	return []interface{}{"EVENT", event}
}