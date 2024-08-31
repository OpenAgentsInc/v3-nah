package nip01

import (
	"encoding/json"
	"github.com/openagentsinc/v3/relay/internal/nostr"
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
	var msg Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	switch msg.Type {
	case EventMessage:
		var event nostr.Event
		err = json.Unmarshal(msg.Data.([]byte), &event)
		if err != nil {
			return nil, err
		}
		msg.Data = &event
	case ReqMessage:
		var filter nostr.Filter
		err = json.Unmarshal(msg.Data.([]byte), &filter)
		if err != nil {
			return nil, err
		}
		msg.Data = &filter
	case AudioMessage:
		var audioData AudioData
		err = json.Unmarshal(msg.Data.([]byte), &audioData)
		if err != nil {
			return nil, err
		}
		msg.Data = &audioData
	}

	return &msg, nil
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