package nip01

import (
	"encoding/json"
	"fmt"
	"github.com/openagentsinc/v3/relay/internal/nostr"
	"log"
	"strconv"
)

// ... (keep the existing code until the handleArrayMessage function)

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

// ... (keep the rest of the existing code)