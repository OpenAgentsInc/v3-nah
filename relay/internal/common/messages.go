package common

import (
	"github.com/openagentsinc/v3/relay/internal/nostr"
)

func CreateEventMessage(event *nostr.Event) []interface{} {
	return []interface{}{"EVENT", event}
}