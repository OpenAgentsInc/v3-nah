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
	log.Printf("Parsing message: %s", string(data))

	var msg Message
	err := json.Unmarshal(data, &msg)
	if err == nil {
		log.Printf("Successfully parsed object message. Type: %s", msg.Type)
		log.Printf("Message Data type: %T", msg.Data)
		return handleObjectMessage(&msg)
	}
	log.Printf("Failed to parse as object: %v", err)

	// If object parsing fails, try array parsing
	var arrayMsg []json.RawMessage
	err = json.Unmarshal(data, &arrayMsg)
	if err == nil && len(arrayMsg) > 0 {
		log.Printf("Successfully parsed as array. Length: %d", len(arrayMsg))
		var msgType string
		err = json.Unmarshal(arrayMsg[0], &msgType)
		if err == nil {
			log.Printf("Array message type: %s", msgType)
			return handleArrayMessage(MessageType(msgType), arrayMsg[1:])
		}
		log.Printf("Failed to parse array message type: %v", err)
	} else {
		log.Printf("Failed to parse as array: %v", err)
	}

	log.Printf("Failed to parse message as either object or array")
	return nil, fmt.Errorf("failed to parse message: %v", err)
}

func handleObjectMessage(msg *Message) (*Message, error) {
	log.Printf("Handling object message of type: %s", msg.Type)
	log.Printf("Message Data type: %T", msg.Data)
	
	switch msg.Type {
	case EventMessage:
		log.Printf("Handling EventMessage")
		var event nostr.Event
		data, err := json.Marshal(msg.Data)
		if err != nil {
			log.Printf("Error marshaling EventMessage data: %v", err)
			return nil, err
		}
		log.Printf("Marshaled EventMessage data: %s", string(data))
		err = json.Unmarshal(data, &event)
		if err != nil {
			log.Printf("Error unmarshaling EventMessage data: %v", err)
			return nil, err
		}
		msg.Data = &event
		log.Printf("Successfully handled EventMessage")
	case ReqMessage:
		log.Printf("Handling ReqMessage")
		var filter nostr.Filter
		data, err := json.Marshal(msg.Data)
		if err != nil {
			log.Printf("Error marshaling ReqMessage data: %v", err)
			return nil, err
		}
		log.Printf("Marshaled ReqMessage data: %s", string(data))
		err = json.Unmarshal(data, &filter)
		if err != nil {
			log.Printf("Error unmarshaling ReqMessage data: %v", err)
			return nil, err
		}
		msg.Data = &filter
		log.Printf("Successfully handled ReqMessage")
	case AudioMessage:
		log.Printf("Handling AudioMessage")
		var audioData AudioData
		data, err := json.Marshal(msg.Data)
		if err != nil {
			log.Printf("Error marshaling AudioMessage data: %v", err)
			return nil, err
		}
		log.Printf("Marshaled AudioMessage data: %s", string(data))
		err = json.Unmarshal(data, &audioData)
		if err != nil {
			log.Printf("Error unmarshaling AudioMessage data: %v", err)
			return nil, err
		}
		msg.Data = &audioData
		log.Printf("Successfully handled AudioMessage")
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}

	return msg, nil
}

func handleArrayMessage(msgType MessageType, data []json.RawMessage) (*Message, error) {
	log.Printf("Handling array message of type: %s", msgType)
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

		log.Printf("Successfully handled array EventMessage")
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
		log.Printf("Successfully handled array ReqMessage")
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
		log.Printf("Successfully handled array AudioMessage")
		return &Message{Type: AudioMessage, Data: &audioData}, nil
	default:
		log.Printf("Unsupported array message type: %s", msgType)
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