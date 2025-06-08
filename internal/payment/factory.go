package gateway

import (
	"context"
	"fmt"
	"sync"

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

// RegisterProvider registers an instantiated provider under a name
func (r *gatewayRegistryService) RegisterProvider(name types.PaymentGatewayProvider, provider GatewayProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[name] = provider
}

// GetProvider returns a provider based on selection attributes (e.g., selected payment provider)
func (r *gatewayRegistryService) GetProvider(ctx context.Context, attrs *types.SelectionAttributes) (GatewayProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, ok := r.providers[attrs.PaymentGateway]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", attrs.PaymentGateway)
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
