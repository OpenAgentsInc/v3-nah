package nostr

// ReqMessage represents a subscription request message
type ReqMessage struct {
	SubscriptionID string
	Filter         Filter
}