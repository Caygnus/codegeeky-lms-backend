package config

import (
	"time"

	"github.com/omkar273/codegeeky/internal/types"
)

// WebhookConfig represents the configuration for the webhook system
type WebhookConfig struct {
	Enabled         bool                         `mapstructure:"enabled"`
	Topic           string                       `mapstructure:"topic" default:"webhooks"`
	PubSub          types.PubSubType             `mapstructure:"pubsub" default:"memory"`
	MaxRetries      int                          `mapstructure:"max_retries" default:"3"`
	InitialInterval time.Duration                `mapstructure:"initial_interval" default:"1s"`
	MaxInterval     time.Duration                `mapstructure:"max_interval" default:"10s"`
	Multiplier      float64                      `mapstructure:"multiplier" default:"2.0"`
	MaxElapsedTime  time.Duration                `mapstructure:"max_elapsed_time" default:"2m"`
	Users           map[string]UserWebhookConfig `mapstructure:"users"`
}

// UserWebhookConfig represents webhook configuration for a specific user
type UserWebhookConfig struct {
	Endpoint       string            `mapstructure:"endpoint"`
	Headers        map[string]string `mapstructure:"headers"`
	Enabled        bool              `mapstructure:"enabled"`
	ExcludedEvents []string          `mapstructure:"excluded_events"`
}
