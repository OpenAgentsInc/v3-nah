package nostr

import (
	"encoding/json"
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