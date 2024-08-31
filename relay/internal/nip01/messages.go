package nip01

import (
	"encoding/json"
	"fmt"
	"github.com/openagentsinc/v3/relay/internal/nostr"
	"log"
	"strconv"
)

type MessageType string

const (
	EventMessage   MessageType = "EVENT"
	ReqMessage     MessageType = "REQ"
	CloseMessage   MessageType = "CLOSE"
	NoticeMessage  MessageType = "NOTICE"
	EoseMessage    MessageType = "EOSE"
	AudioMessage   MessageType = "AUDIO"
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

	var arrayMsg []json.RawMessage
	err := json.Unmarshal(data, &arrayMsg)
	if err == nil && len(arrayMsg) > 0 {
		log.Printf("Message is an array with %d elements", len(arrayMsg))
		var msgType string
		err = json.Unmarshal(arrayMsg[0], &msgType)
		if err == nil {
			log.Printf("Message type: %s", msgType)
			return handleArrayMessage(MessageType(msgType), arrayMsg[1:])
		}
	}

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
		var rawEvent map[string]interface{}
		err := json.Unmarshal(data[0], &rawEvent)
		if err != nil {
			return nil, err
		}
		
		// Convert created_at to int64
		switch createdAt := rawEvent["created_at"].(type) {
		case float64:
			rawEvent["created_at"] = int64(createdAt)
		case string:
			createdAtInt, err := strconv.ParseInt(createdAt, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid created_at value: %v", err)
			}
			rawEvent["created_at"] = createdAtInt
		}
		
		eventJSON, err := json.Marshal(rawEvent)
		if err != nil {
			return nil, err
		}
		
		var event nostr.Event
		err = json.Unmarshal(eventJSON, &event)
		if err != nil {
			return nil, err
		}
		return &Message{Type: EventMessage, Data: &event}, nil
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
			Kind:    1235, // Custom event kind for transcription response
			Content: transcription,
		},
	}, nil
}