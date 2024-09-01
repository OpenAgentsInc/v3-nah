package nostr

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Event struct {
	ID        string    `json:"id"`
	PubKey    string    `json:"pubkey"`
	CreatedAt time.Time `json:"created_at"`
	Kind      int       `json:"kind"`
	Tags      [][]string `json:"tags"`
	Content   string    `json:"content"`
	Sig       string    `json:"sig"`
}

func (e *Event) UnmarshalJSON(data []byte) error {
	type Alias Event
	aux := &struct {
		CreatedAt interface{} `json:"created_at"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	switch v := aux.CreatedAt.(type) {
	case float64:
		e.CreatedAt = time.Unix(int64(v), 0)
	case string:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			e.CreatedAt = time.Unix(i, 0)
		} else if t, err := time.Parse(time.RFC3339, v); err == nil {
			e.CreatedAt = t
		} else {
			return fmt.Errorf("invalid timestamp format: %s", v)
		}
	default:
		return fmt.Errorf("invalid timestamp type: %T", v)
	}
	return nil
}

func (e *Event) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func DeserializeEvent(data []byte) (*Event, error) {
	var e Event
	err := json.Unmarshal(data, &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}