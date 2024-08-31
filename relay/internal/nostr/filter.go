package nostr

import (
	"time"
)

type Filter struct {
	IDs     []string  `json:"ids,omitempty"`
	Authors []string  `json:"authors,omitempty"`
	Kinds   []int     `json:"kinds,omitempty"`
	Since   time.Time `json:"since,omitempty"`
	Until   time.Time `json:"until,omitempty"`
	Limit   int       `json:"limit,omitempty"`
}

func (f *Filter) Match(e *Event) bool {
	if len(f.IDs) > 0 && !contains(f.IDs, e.ID) {
		return false
	}
	if len(f.Authors) > 0 && !contains(f.Authors, e.PubKey) {
		return false
	}
	if len(f.Kinds) > 0 && !containsInt(f.Kinds, e.Kind) {
		return false
	}
	if !f.Since.IsZero() && e.CreatedAt.Before(f.Since) {
		return false
	}
	if !f.Until.IsZero() && e.CreatedAt.After(f.Until) {
		return false
	}
	return true
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func containsInt(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}