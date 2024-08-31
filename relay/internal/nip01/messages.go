package nip01

import (
	"encoding/json"
	"fmt"
	"github.com/openagentsinc/v3/relay/internal/nostr"
	"log"
)

type MessageType string

const (
	EventMessage   MessageType = "EVENT"
	ReqMessage     MessageType = "REQ"
	CloseMessage   MessageType = "CLOSE"
	NoticeMessage  MessageType = "NOTICE"
	EoseMessage    MessageType = "EOSE"
	AudioMessage   MessageType = "AUDIO" // New message type for audio
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
	log.Printf("Received raw message: %s", string(data))

	// Try to unmarshal as an array first
	var arrayMsg []json.RawMessage
	err := json.Unmarshal(data, &arrayMsg)
	if err == nil {
		log.Printf("Message is an array with %d elements", len(arrayMsg))
		if len(arrayMsg) > 0 {
			// Check if the first element is a string (message type)
			var msgType string
			err = json.Unmarshal(arrayMsg[0], &msgType)
			if err == nil {
				log.Printf("Message type: %s", msgType)
				// Handle array-style messages
				return handleArrayMessage(MessageType(msgType), arrayMsg[1:])
			}
		}
	}

	// If not an array, try to unmarshal as a regular Message struct
	var msg Message
	err = json.Unmarshal(data, &msg)
	if err != nil {
		log.Printf("Error unmarshaling as Message struct: %v", err)
		return nil, err
	}

	log.Printf("Parsed message type: %s", msg.Type)

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
		var event nostr.Event
		err := json.Unmarshal(data[0], &event)
		if err != nil {
			return nil, err
		}
		return &Message{Type: EventMessage, Data: &event}, nil
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
	// Add other message types as needed
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

func CreateReqMessage(filter *nostr.Filter) (*Message, error) {
	return &Message{
		Type: ReqMessage,
		Data: filter,
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
			Kind:    1235, // Custom event kind for transcription response
			Content: transcription,
		},
	}, nil
}