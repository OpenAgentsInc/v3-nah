package nip90

import (
	"fmt"
	"log"
	"strings"

	"github.com/openagentsinc/v3/relay/internal/nostr"
)

func LogEventDetails(event *nostr.Event) {
	var sb strings.Builder

	sb.WriteString("Event Details:\n")
	sb.WriteString(fmt.Sprintf("  ID: %s\n", event.ID))
	sb.WriteString(fmt.Sprintf("  PubKey: %s\n", event.PubKey))
	sb.WriteString(fmt.Sprintf("  CreatedAt: %s\n", event.CreatedAt))
	sb.WriteString(fmt.Sprintf("  Kind: %d\n", event.Kind))
	sb.WriteString("  Tags:\n")
	for _, tag := range event.Tags {
		sb.WriteString(fmt.Sprintf("    - %v\n", tag))
	}
	sb.WriteString(fmt.Sprintf("  Content: %s\n", event.Content))
	sb.WriteString(fmt.Sprintf("  Sig: %s\n", event.Sig))

	log.Print(sb.String())
}