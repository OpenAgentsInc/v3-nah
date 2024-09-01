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
	AudioMessage  MessageType = "AUDIO"
)

type Message struct {
	Type MessageType `json:"type"`
	Data interface{} `json:"data"`
}

type AudioData struct {
	Audio  string `json:"audio"`
	Format string `json:"format"`
}

func ParseMessage(data []byte) (*Message, error) {
	var arrayMsg []json.RawMessage
	err := json.Unmarshal(data, &arrayMsg)
	if err == nil && len(arrayMsg) > 0 {
		var msgType string
		err = json.Unmarshal(arrayMsg[0], &msgType)
		if err == nil {
			log.Printf("Received message type: %s", msgType)
			return handleArrayMessage(MessageType(msgType), arrayMsg[1:])
		}
	}

	var msg Message
	err = json.Unmarshal(data, &msg)
	if err != nil {
		log.Printf("Error unmarshaling as Message struct: %v", err)
		return nil, err
	}

	// log.Printf("Parsed message type: %s", msg.Type)

	switch msg.Type {
	case EventMessage:
		var event nostr.Event
		err = json.Unmarshal(msg.Data.([]byte), &event)
		if err != nil {
			log.Printf("Error unmarshaling EventMessage data: %v", err)
			return nil, err
		}
		msg.Data = &event
	case ReqMessage:
		var filter nostr.Filter
		err = json.Unmarshal(msg.Data.([]byte), &filter)
		if err != nil {
			log.Printf("Error unmarshaling ReqMessage data: %v", err)
			return nil, err
		}
		msg.Data = &filter
	case AudioMessage:
		var audioData AudioData
		err = json.Unmarshal(msg.Data.([]byte), &audioData)
		if err != nil {
			log.Printf("Error unmarshaling AudioMessage data: %v", err)
			return nil, err
		}
		msg.Data = &audioData
	}

	return &msg, nil
}

func handleArrayMessage(msgType MessageType, data []json.RawMessage) (*Message, error) {
	switch msgType {
	case EventMessage:
		if len(data) != 1 {
			return nil, fmt.Errorf("invalid EVENT message format")
		}
		var rawEvent map[string]interface{}
		err := json.Unmarshal(data[0], &rawEvent)
		if err != nil {
			return nil, err
		}

		// Manually construct the Event struct
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
			}
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

		return &Message{Type: EventMessage, Data: event}, nil
	case ReqMessage:
		if len(data) < 2 {
			return nil, fmt.Errorf("invalid REQ message format")
		}
		var subscriptionID string
		err := json.Unmarshal(data[0], &subscriptionID)
		if err != nil {
			return nil, err
		}
		var filter nostr.Filter
		err = json.Unmarshal(data[1], &filter)
		if err != nil {
			return nil, err
		}
		return &Message{Type: ReqMessage, Data: &nostr.ReqMessage{SubscriptionID: subscriptionID, Filter: filter}}, nil
	case AudioMessage:
		if len(data) != 1 {
			return nil, fmt.Errorf("invalid AUDIO message format")
		}
		var audioData AudioData
		err := json.Unmarshal(data[0], &audioData)
		if err != nil {
			return nil, err
		}
		return &Message{Type: AudioMessage, Data: &audioData}, nil
	default:
		return nil, fmt.Errorf("unsupported array message type: %s", msgType)
	}
}

func CreateEventMessage(event *nostr.Event) (*Message, error) {
	return &Message{
		Type: EventMessage,
		Data: event,
	}, nil
}

func CreateReqMessage(reqMsg *nostr.ReqMessage) (*Message, error) {
	return &Message{
		Type: ReqMessage,
		Data: reqMsg,
	}, nil
}

func CreateCloseMessage(subscriptionID string) (*Message, error) {
	return &Message{
		Type: CloseMessage,
		Data: subscriptionID,
	}, nil
}

func CreateNoticeMessage(message string) (*Message, error) {
	return &Message{
		Type: NoticeMessage,
		Data: message,
	}, nil
}

func CreateEoseMessage(subscriptionID string) (*Message, error) {
	return &Message{
		Type: EoseMessage,
		Data: subscriptionID,
	}, nil
}

func CreateAudioResponseMessage(transcription string) (*Message, error) {
	return &Message{
		Type: EventMessage,
		Data: &nostr.Event{
			Kind:    6252, // Custom event kind for transcription response
			Content: transcription,
		},
	}, nil
}
