package nip01

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/openagentsinc/v3/relay/internal/nostr"
)

// ... (keep the existing code unchanged)

func handleObjectMessage(msg *Message) (*Message, error) {
	log.Printf("Handling object message of type: %s", msg.Type)
	log.Printf("Message Data type: %T", msg.Data)
	
	switch msg.Type {
	case EventMessage:
		log.Printf("Handling EventMessage")
		var rawEvent map[string]interface{}
		data, err := json.Marshal(msg.Data)
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
	// ... (keep the rest of the function unchanged)
	}

	return msg, nil
}

// ... (keep the rest of the file unchanged)