package nip01

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/openagentsinc/v3/relay/internal/nostr"
)

// ... (keep the existing type definitions and constants)

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

// ... (keep the rest of the file unchanged)