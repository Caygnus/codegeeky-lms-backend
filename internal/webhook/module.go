package webhook

import (
	"github.com/omkar273/codegeeky/internal/config"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/pubsub"
	"github.com/omkar273/codegeeky/internal/pubsub/memory"
	"github.com/omkar273/codegeeky/internal/service"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/omkar273/codegeeky/internal/webhook/handler"
	"github.com/omkar273/codegeeky/internal/webhook/payload"
	"github.com/omkar273/codegeeky/internal/webhook/publisher"
	"go.uber.org/fx"
)

// Module provides all webhook-related dependencies
var Module = fx.Options(
	// Core dependencies
	fx.Provide(
		// PubSub for sending webhook events
		providePubSub,
	),

	// Webhook components
	fx.Provide(
		// Publisher for sending webhook events
		publisher.NewPublisher,

		// Handler for processing webhook events
		handler.NewHandler,

		// Payload builder factory and services
		providePayloadBuilderFactory,

		// Main webhook service
		NewWebhookService,
	),
)

// providePayloadBuilderFactory creates a new payload builder factory with all required services
func providePayloadBuilderFactory(
	userService service.UserService,
	authService service.AuthService,
	categoryService service.CategoryService,
	onboardingService service.OnboardingService,
	internshipService service.InternshipService,
) payload.PayloadBuilderFactory {
	services := payload.NewServices(
		userService,
		authService,
		categoryService,
		onboardingService,
		internshipService,
	)
	return payload.NewPayloadBuilderFactory(services)
}

func providePubSub(
	cfg *config.Configuration,
	logger *logger.Logger,
) pubsub.PubSub {
	switch cfg.Webhook.PubSub {
	case types.MemoryPubSub:
		return memory.NewPubSub(cfg, logger)
	case types.KafkaPubSub:
		// TODO: implement
	}
	panic("unsupported pubsub type")
}
