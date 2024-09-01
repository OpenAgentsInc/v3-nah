package nip01

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/openagentsinc/v3/relay/internal/nostr"
)

type MessageType string

const (
	EventMessage  MessageType = "EVENT"
	ReqMessage    MessageType = "REQ"
	CloseMessage  MessageType = "CLOSE"
	NoticeMessage MessageType = "NOTICE"
	EoseMessage   MessageType = "EOSE"
	OkMessage     MessageType = "OK"
)

type Message struct {
	Type MessageType   `json:"type"`
	Data interface{}   `json:"data"`
	Raw  []interface{} `json:"-"`
}

type AudioData struct {
	Data   string `json:"data"`
	Format string `json:"format"`
}

func ParseMessage(data []byte) (*Message, error) {
	var raw []json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {
		// Try unmarshaling as a single JSON object
		var singleObject map[string]interface{}
		err = json.Unmarshal(data, &singleObject)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling message: %v", err)
		}
		// Convert single object to array format
		raw = []json.RawMessage{json.RawMessage(data)}
	}

	if len(raw) == 0 {
		return nil, fmt.Errorf("empty message")
	}

	var typeStr string
	err = json.Unmarshal(raw[0], &typeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid message type: %v", err)
	}

	msg := &Message{
		Type: MessageType(typeStr),
		Raw:  make([]interface{}, len(raw)),
	}

	for i, item := range raw {
		var v interface{}
		err = json.Unmarshal(item, &v)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling item %d: %v", i, err)
		}
		msg.Raw[i] = v
	}

	switch msg.Type {
	case EventMessage:
		return handleEventMessage(msg)
	case ReqMessage:
		return handleReqMessage(msg)
	case CloseMessage:
		return handleCloseMessage(msg)
	case NoticeMessage, EoseMessage, OkMessage:
		return handleSimpleMessage(msg)
	default:
		return nil, fmt.Errorf("unknown message type: %s", msg.Type)
	}
}

func handleEventMessage(msg *Message) (*Message, error) {
	if len(msg.Raw) < 2 {
		return nil, fmt.Errorf("invalid EVENT message")
	}

	eventData, ok := msg.Raw[1].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid EVENT data format")
	}

	event := &nostr.Event{}
	if id, ok := eventData["id"].(string); ok {
		event.ID = id
	}
	if pubkey, ok := eventData["pubkey"].(string); ok {
		event.PubKey = pubkey
	}
	if kind, ok := eventData["kind"].(float64); ok {
		event.Kind = int(kind)
	}
	if content, ok := eventData["content"].(string); ok {
		event.Content = content
	}
	if sig, ok := eventData["sig"].(string); ok {
		event.Sig = sig
	}

	// Handle created_at
	if createdAt, ok := eventData["created_at"].(float64); ok {
		event.CreatedAt = time.Unix(int64(createdAt), 0)
	} else if createdAtStr, ok := eventData["created_at"].(string); ok {
		if createdAtInt, err := strconv.ParseInt(createdAtStr, 10, 64); err == nil {
			event.CreatedAt = time.Unix(createdAtInt, 0)
		} else {
			// Try parsing as RFC3339 format
			if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
				event.CreatedAt = t
			} else {
				log.Printf("Error parsing created_at: %v", err)
			}
		}
	} else {
		log.Printf("created_at field not found or has unexpected type: %T", eventData["created_at"])
	}

	// Handle tags
	if tags, ok := eventData["tags"].([]interface{}); ok {
		for _, tag := range tags {
			if tagSlice, ok := tag.([]interface{}); ok {
				var stringSlice []string
				for _, item := range tagSlice {
					if str, ok := item.(string); ok {
						stringSlice = append(stringSlice, str)
					}
				}
				event.Tags = append(event.Tags, stringSlice)
			}
		}
	}

	msg.Data = event
	return msg, nil
}

func handleReqMessage(msg *Message) (*Message, error) {
	if len(msg.Raw) < 2 {
		return nil, fmt.Errorf("invalid REQ message")
	}
	msg.Data = msg.Raw[1:]
	return msg, nil
}

func handleCloseMessage(msg *Message) (*Message, error) {
	if len(msg.Raw) < 2 {
		return nil, fmt.Errorf("invalid CLOSE message")
	}
	msg.Data = msg.Raw[1]
	return msg, nil
}

func handleSimpleMessage(msg *Message) (*Message, error) {
	if len(msg.Raw) < 2 {
		return nil, fmt.Errorf("invalid %s message", msg.Type)
	}
	msg.Data = msg.Raw[1]
	return msg, nil
}

func CreateEventMessage(event *nostr.Event) *Message {
	return &Message{
		Type: EventMessage,
		Data: event,
		Raw:  []interface{}{"EVENT", event},
	}
}