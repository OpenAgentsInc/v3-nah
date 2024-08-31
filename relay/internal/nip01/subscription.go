package nip01

import (
	"github.com/openagentsinc/v3/relay/internal/nostr"
	"sync"
)

type Subscription struct {
	ID     string
	Filters []*nostr.Filter
	Events  chan *nostr.Event
}

type SubscriptionManager struct {
	subscriptions map[string]*Subscription
	mu            sync.RWMutex
}

func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		subscriptions: make(map[string]*Subscription),
	}
}

func (sm *SubscriptionManager) AddSubscription(id string, filters []*nostr.Filter) *Subscription {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sub := &Subscription{
		ID:     id,
		Filters: filters,
		Events:  make(chan *nostr.Event, 100), // Buffered channel to prevent blocking
	}
	sm.subscriptions[id] = sub
	return sub
}

func (sm *SubscriptionManager) RemoveSubscription(id string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sub, ok := sm.subscriptions[id]; ok {
		close(sub.Events)
		delete(sm.subscriptions, id)
	}
}

func (sm *SubscriptionManager) GetSubscription(id string) (*Subscription, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sub, ok := sm.subscriptions[id]
	return sub, ok
}

func (sm *SubscriptionManager) BroadcastEvent(event *nostr.Event) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	for _, sub := range sm.subscriptions {
		for _, filter := range sub.Filters {
			if filter.Match(event) {
				select {
				case sub.Events <- event:
					// Event sent successfully
				default:
					// Channel is full, consider handling this case (e.g., logging, dropping events)
				}
				break // Move to the next subscription once we've matched and sent the event
			}
		}
	}
}