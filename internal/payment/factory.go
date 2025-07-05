package gateway

import (
	"context"
	"sync"

	"github.com/omkar273/codegeeky/internal/config"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/httpclient"
	razorpay "github.com/omkar273/codegeeky/internal/payment/providers"
	"github.com/omkar273/codegeeky/internal/types"
)

// GatewayRegistryService manages registered payment gateway providers
type gatewayRegistryService struct {
	providers map[types.PaymentGatewayProvider]GatewayProvider
	mu        sync.RWMutex
}

// NewGatewayRegistryService creates a new registry
func NewGatewayRegistryService() GatewayRegistryService {
	return &gatewayRegistryService{
		providers: make(map[types.PaymentGatewayProvider]GatewayProvider),
	}
}

func InitializeProviders(ctx context.Context, client *httpclient.Client, config *config.Configuration) GatewayRegistryService {
	registry := NewGatewayRegistryService()

	// razorpay
	registry.RegisterProvider(
		types.PaymentGatewayProviderRazorpay,
		razorpay.NewRazorpayProvider(client, config),
	)

	return registry
}

// RegisterProvider registers an instantiated provider under a name
func (r *gatewayRegistryService) RegisterProvider(name types.PaymentGatewayProvider, provider GatewayProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[name] = provider
}

// GetProviderByName returns a provider based on name (e.g., selected payment provider)
func (r *gatewayRegistryService) GetProviderByName(ctx context.Context, name types.PaymentGatewayProvider) (GatewayProvider, error) {
	if name == "" {
		return nil, ierr.NewError("provider name is required").
			WithHint("Please provide a valid payment gateway provider name").
			Mark(ierr.ErrValidation)
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, ok := r.providers[name]
	if !ok {
		return nil, ierr.NewError("provider not found").
			WithHint("Please provide a valid payment gateway provider name").
			WithReportableDetails(map[string]any{
				"name": name,
			}).
			Mark(ierr.ErrValidation)
	}

	return provider, nil
}

func (r *gatewayRegistryService) ListAvailableProviders(ctx context.Context) ([]types.PaymentGatewayProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	providers := make([]types.PaymentGatewayProvider, 0, len(r.providers))
	for provider := range r.providers {
		providers = append(providers, provider)
	}
	return providers, nil
}
