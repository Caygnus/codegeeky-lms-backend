package service

import (
	"github.com/omkar273/codegeeky/internal/config"
	"github.com/omkar273/codegeeky/internal/domain/discount"
	"github.com/omkar273/codegeeky/internal/domain/payment"
	"github.com/omkar273/codegeeky/internal/domain/user"
	"github.com/omkar273/codegeeky/internal/httpclient"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/postgres"
	"github.com/omkar273/codegeeky/internal/webhook/publisher"
	"go.uber.org/fx"
)

type ServiceParams struct {
	fx.In

	// Core dependencies
	Logger *logger.Logger
	Config *config.Configuration
	DB     postgres.IClient

	// Repository dependencies
	UserRepo     user.Repository
	DiscountRepo discount.Repository
	PaymentRepo  payment.Repository

	// Service dependencies
	WebhookPublisher publisher.WebhookPublisher

	// http client
	HTTPClient httpclient.Client
}
