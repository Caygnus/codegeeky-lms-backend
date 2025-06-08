package payload

import (
	"fmt"
)

// PayloadBuilderFactory interface for getting event-specific payload builders
type PayloadBuilderFactory interface {
	GetBuilder(eventType string) (PayloadBuilder, error)
}

type payloadBuilderFactory struct {
	builders map[string]func() PayloadBuilder
	services *Services
}

// NewPayloadBuilderFactory creates a new factory with registered builders
func NewPayloadBuilderFactory(services *Services) PayloadBuilderFactory {
	f := &payloadBuilderFactory{
		builders: make(map[string]func() PayloadBuilder),
		services: services,
	}

	// f.builders[types.WebhookEventUserCreated] = func() PayloadBuilder {
	// 	return &userCreatedPayloadBuilder{services: f.services}
	// }

	return f
}

// GetBuilder returns a payload builder for the given event type
func (f *payloadBuilderFactory) GetBuilder(eventType string) (PayloadBuilder, error) {
	builderFn, ok := f.builders[eventType]
	if !ok {
		return nil, fmt.Errorf("no builder registered for event type: %s", eventType)
	}

	return builderFn(), nil
}
