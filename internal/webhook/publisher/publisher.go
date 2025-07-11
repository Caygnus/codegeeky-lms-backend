package publisher

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/omkar273/codegeeky/internal/config"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/pubsub"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
)

// WebhookPublisher interface for producing webhook events
type WebhookPublisher interface {
	PublishWebhook(ctx context.Context, event *types.WebhookEvent) error
	Close() error
}

// Handler implements handler.Handler using watermill's gochannel
type webhookPublisher struct {
	pubSub pubsub.PubSub
	config *config.WebhookConfig
	logger *logger.Logger
}

// NewHandler creates a new memory-based handler
func NewPublisher(
	pubSub pubsub.PubSub,
	cfg *config.Configuration,
	logger *logger.Logger,
) (WebhookPublisher, error) {
	return &webhookPublisher{
		pubSub: pubSub,
		config: &cfg.Webhook,
		logger: logger,
	}, nil
}

func (p *webhookPublisher) PublishWebhook(ctx context.Context, event *types.WebhookEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	messageID := event.ID
	if messageID == "" {
		messageID = watermill.NewUUID()
	}

	msg := message.NewMessage(messageID, payload)
	msg.Metadata.Set("user_id", *event.UserID)

	p.logger.Debugw("publishing webhook event",
		"event_id", event.ID,
		"event_name", event.EventName,
		"user_id", lo.FromPtr(event.UserID),
		"topic", p.config.Topic,
		"payload", string(payload),
	)

	if err := p.pubSub.Publish(ctx, p.config.Topic, msg); err != nil {
		p.logger.Errorw("failed to publish webhook event",
			"error", err,
			"event_id", event.ID,
			"event_name", event.EventName,
			"user_id", lo.FromPtr(event.UserID),
		)
		return err
	}

	p.logger.Infow("successfully published webhook event",
		"event_id", event.ID,
		"event_name", event.EventName,
		"user_id", lo.FromPtr(event.UserID),
	)

	return nil
}

// Close closes the publisher
func (p *webhookPublisher) Close() error {
	return p.pubSub.Close()
}
