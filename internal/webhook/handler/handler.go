package handler

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/omkar273/codegeeky/internal/config"
	"github.com/omkar273/codegeeky/internal/httpclient"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/pubsub"
	pubsubRouter "github.com/omkar273/codegeeky/internal/pubsub/router"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/omkar273/codegeeky/internal/webhook/payload"
)

type Handler interface {
	RegisterHandler(router *pubsubRouter.Router) error
}

// handler implements handler.Handler using watermill's gochannel
type handler struct {
	pubSub   pubsub.PubSub
	config   *config.Webhook
	factory  payload.PayloadBuilderFactory
	client   httpclient.Client
	logger   *logger.Logger
	services *payload.Services
}

func NewHandler(
	pubSub pubsub.PubSub,
	config *config.Webhook,
	factory payload.PayloadBuilderFactory,
	client httpclient.Client,
	logger *logger.Logger,
	services *payload.Services,
) Handler {
	return &handler{
		pubSub:   pubSub,
		config:   config,
		factory:  factory,
		client:   client,
		logger:   logger,
		services: services,
	}
}

func (h *handler) RegisterHandler(router *pubsubRouter.Router) error {
	router.AddNoPublishHandler(
		"webhook_handler",
		h.config.Topic,
		h.pubSub,
		h.processMessage,
	)
	return nil
}

// processMessage processes a single webhook message
func (h *handler) processMessage(msg *message.Message) error {
	ctx := msg.Context()

	// log the context fields like user_id, event_name, etc
	h.logger.Debugw("context",
		"user_id", types.GetUserID(ctx),
		"event_name", types.GetRequestID(ctx),
	)

	var event types.WebhookEvent
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		h.logger.Errorw("failed to unmarshal webhook event",
			"error", err,
			"message_uuid", msg.UUID,
		)
		return nil // Don't retry on unmarshal errors
	}

	// Get user config
	userCfg, ok := h.config.Users[*event.UserID]
	if !ok {
		h.logger.Warnw("user config not found",
			"user_id", event.UserID,
			"message_uuid", msg.UUID,
		)
		// Don't retry if user not found
		return nil
	}

	// Check if user webhooks are enabled
	if !userCfg.Enabled {
		h.logger.Debugw("webhooks disabled for user",
			"user_id", event.UserID,
			"message_uuid", msg.UUID,
		)
		return nil
	}

	// Check if event is excluded
	for _, excludedEvent := range userCfg.ExcludedEvents {
		if excludedEvent == event.EventName {
			h.logger.Debugw("event excluded for user",
				"user_id", event.UserID,
				"event", event.EventName,
			)
			return nil
		}
	}

	// Build event payload
	builder, err := h.factory.GetBuilder(event.EventName)
	if err != nil {
		return err
	}

	h.logger.Debugw("building webhook payload",
		"event_name", event.EventName,
		"builder", builder,
	)

	// set user_id in context
	ctx = context.WithValue(ctx, types.CtxUserID, *event.UserID)
	webHookPayload, err := builder.BuildPayload(ctx, event.EventName, event.Payload)
	if err != nil {
		return err
	}

	h.logger.Debugw("built webhook payload",
		"event_name", event.EventName,
		"payload", string(webHookPayload),
	)

	// Send webhook
	req := &httpclient.Request{
		Method:  "POST",
		URL:     userCfg.Endpoint,
		Headers: userCfg.Headers,
		Body:    webHookPayload,
	}

	resp, err := h.client.Send(ctx, req)
	if err != nil {
		h.logger.Errorw("failed to send webhook",
			"error", err,
			"message_uuid", msg.UUID,
			"user_id", event.UserID,
			"event", event.EventName,
		)
		return err
	}

	h.logger.Infow("webhook sent successfully",
		"message_uuid", msg.UUID,
		"user_id", event.UserID,
		"event", event.EventName,
		"status_code", resp.StatusCode,
	)

	return nil
}
