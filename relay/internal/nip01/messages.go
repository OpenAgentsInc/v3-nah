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
	var raw []interface{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling message: %v", err)
	}

	if len(raw) == 0 {
		return nil, fmt.Errorf("empty message")
	}

	typeStr, ok := raw[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid message type")
	}

	msg := &Message{
		Type: MessageType(typeStr),
		Raw:  raw,
	}

	switch msg.Type {
	case EventMessage, ReqMessage, CloseMessage:
		return handleObjectMessage(msg)
	case NoticeMessage, EoseMessage, OkMessage:
		return handleSimpleMessage(msg)
	default:
		return nil, fmt.Errorf("unknown message type: %s", msg.Type)
	}
}

func handleObjectMessage(msg *Message) (*Message, error) {
	log.Printf("Handling object message of type: %s", msg.Type)
	log.Printf("Message Data type: %T", msg.Data)
	
	switch msg.Type {
	case EventMessage:
		log.Printf("Handling EventMessage")
		var rawEvent map[string]interface{}
		data, err := json.Marshal(msg.Raw[1])
		if err != nil {
			log.Printf("Error marshaling EventMessage data: %v", err)
			return nil, err
		}
		log.Printf("Marshaled EventMessage data: %s", string(data))
		err = json.Unmarshal(data, &rawEvent)
		if err != nil {
			log.Printf("Error unmarshaling EventMessage data to map: %v", err)
			return nil, err
		}

		log.Printf("Raw event data: %+v", rawEvent)

		event := &nostr.Event{}
		if id, ok := rawEvent["id"].(string); ok {
			event.ID = id
		}
		if pubkey, ok := rawEvent["pubkey"].(string); ok {
			event.PubKey = pubkey
		}
		if kind, ok := rawEvent["kind"].(float64); ok {
			event.Kind = int(kind)
		}
		if content, ok := rawEvent["content"].(string); ok {
			event.Content = content
		}
		if sig, ok := rawEvent["sig"].(string); ok {
			event.Sig = sig
		}

		// Handle created_at
		if createdAt, ok := rawEvent["created_at"].(float64); ok {
			event.CreatedAt = time.Unix(int64(createdAt), 0)
		} else if createdAtStr, ok := rawEvent["created_at"].(string); ok {
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
		} else if createdAtNum, ok := rawEvent["created_at"].(json.Number); ok {
			if createdAtInt, err := createdAtNum.Int64(); err == nil {
				event.CreatedAt = time.Unix(createdAtInt, 0)
			} else {
				log.Printf("Error parsing created_at as json.Number: %v", err)
			}
		} else {
			log.Printf("created_at field not found or has unexpected type: %T", rawEvent["created_at"])
		}

		// Handle tags
		if tags, ok := rawEvent["tags"].([]interface{}); ok {
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
		log.Printf("Successfully handled EventMessage")
	case ReqMessage:
		// Handle REQ message
		if len(msg.Raw) < 2 {
			return nil, fmt.Errorf("invalid REQ message")
		}
		msg.Data = msg.Raw[1:]
	case CloseMessage:
		// Handle CLOSE message
		if len(msg.Raw) < 2 {
			return nil, fmt.Errorf("invalid CLOSE message")
		}
		msg.Data = msg.Raw[1]
	}

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