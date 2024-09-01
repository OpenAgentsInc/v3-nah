package nip01

import (
	"encoding/json"
	"fmt"

	"github.com/openagentsinc/v3/relay/internal/nostr"
)

type MessageType int

const (
	EventMessage MessageType = iota
	ReqMessage
	CloseMessage
)

type Message struct {
	Type MessageType
	Data interface{}
}

func ParseMessage(data []byte) (*Message, error) {
	var rawMessage []json.RawMessage
	err := json.Unmarshal(data, &rawMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to parse message: %v", err)
	}

	if len(rawMessage) < 2 {
		return nil, fmt.Errorf("invalid message format")
	}

	var messageType string
	err = json.Unmarshal(rawMessage[0], &messageType)
	if err != nil {
		return nil, fmt.Errorf("failed to parse message type: %v", err)
	}

	switch messageType {
	case "EVENT":
		var event nostr.Event
		err = json.Unmarshal(rawMessage[1], &event)
		if err != nil {
			return nil, fmt.Errorf("failed to parse event: %v", err)
		}
		return &Message{Type: EventMessage, Data: &event}, nil
	case "REQ":
		var reqData []interface{}
		err = json.Unmarshal(data, &reqData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse REQ message: %v", err)
		}
		return &Message{Type: ReqMessage, Data: reqData[1:]}, nil
	case "CLOSE":
		var closeData []interface{}
		err = json.Unmarshal(data, &closeData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CLOSE message: %v", err)
		}
		if len(closeData) < 2 {
			return nil, fmt.Errorf("invalid CLOSE message format")
		}
		subscriptionID, ok := closeData[1].(string)
		if !ok {
			return nil, fmt.Errorf("invalid subscription ID in CLOSE message")
		}
		return &Message{Type: CloseMessage, Data: subscriptionID}, nil
	default:
		return nil, fmt.Errorf("unknown message type: %s", messageType)
	}
}