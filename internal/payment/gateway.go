package gateway

import (
	"context"

	"github.com/omkar273/codegeeky/internal/api/dto"
	"github.com/omkar273/codegeeky/internal/types"
)

// GatewayProvider defines a base interface for any external payment provider integration
type GatewayProvider interface {
	// Identifies the provider (e.g., "stripe", "razorpay", etc.)
	ProviderName() types.PaymentGatewayProvider

	// Returns a list of supported capabilities like "refunds", "webhooks", "upi", etc.
	SupportedFeatures() []types.PaymentGatewayFeatures

	// Initialize or prepare the provider (optional)
	Initialize(ctx context.Context) error

	// ProcessWebhook is a generic handler for incoming webhooks from the provider
	ProcessWebhook(ctx context.Context, payload []byte, headers map[string]string) (*dto.WebhookResult, error)

	// CreatePaymentOrder is a generic abstraction to initiate payments
	CreatePaymentOrder(ctx context.Context, input *dto.PaymentRequest) (*dto.PaymentResponse, error)

	// VerifyPaymentStatus checks status of a payment from provider
	VerifyPaymentStatus(ctx context.Context, providerPaymentID string) (*dto.PaymentStatus, error)
}

type GatewayRegistryService interface {
	GetProviderByName(ctx context.Context, name types.PaymentGatewayProvider) (GatewayProvider, error)
	ListAvailableProviders(ctx context.Context) ([]types.PaymentGatewayProvider, error)
	RegisterProvider(name types.PaymentGatewayProvider, provider GatewayProvider)
}
