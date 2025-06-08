package gateway

import (
	"context"
	"fmt"
	"sync"

	"github.com/omkar273/codegeeky/internal/types"
)

// GatewayRegistryService manages registered payment gateway providers
type gatewayRegistryService struct {
	providers map[types.PaymentProvider]GatewayProvider
	mu        sync.RWMutex
}

// NewGatewayRegistryService creates a new registry
func NewGatewayRegistryService() GatewayRegistryService {
	return &gatewayRegistryService{
		providers: make(map[types.PaymentProvider]GatewayProvider),
	}
}

// RegisterProvider registers an instantiated provider under a name
func (r *gatewayRegistryService) RegisterProvider(name types.PaymentProvider, provider GatewayProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[name] = provider
}

// GetProvider returns a provider based on selection attributes (e.g., selected payment provider)
func (r *gatewayRegistryService) GetProvider(ctx context.Context, attrs *types.SelectionAttributes) (GatewayProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, ok := r.providers[attrs.PaymentProvider]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", attrs.PaymentProvider)
	}

	return provider, nil
}

func (r *gatewayRegistryService) ListAvailableProviders(ctx context.Context) ([]types.PaymentProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	providers := make([]types.PaymentProvider, 0, len(r.providers))
	for provider := range r.providers {
		providers = append(providers, provider)
	}
	return providers, nil
}
