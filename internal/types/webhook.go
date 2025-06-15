package types

import (
	"encoding/json"
	"time"
)

// WebhookEvent represents a webhook event to be delivered
type WebhookEvent struct {
	ID        string          `json:"id"`
	EventName string          `json:"event_name"`
	UserID    *string         `json:"user_id,omitempty"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp time.Time       `json:"timestamp"`
}

// user events
const (
	WebhookEventUserCreated = "user.created"
	WebhookEventUserUpdated = "user.updated"
	WebhookEventUserDeleted = "user.deleted"
	WebhookEventUserLogin   = "user.login"
	WebhookEventUserLogout  = "user.logout"
)

// EventSource defines the source of an event
type EventSource string

const (
	EventSourceInternal EventSource = "internal"
	EventSourceRazorpay EventSource = "razorpay"
)
