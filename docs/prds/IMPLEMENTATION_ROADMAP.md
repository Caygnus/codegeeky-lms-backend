# Implementation Roadmap

This document provides a step-by-step roadmap for implementing the enhanced event-driven pubsub system.

## Overview

The implementation is divided into 4 phases, each building upon the previous phase:

1. **Phase 1**: Enhanced PubSub Infrastructure (Weeks 1-2)
2. **Phase 2**: External Integration (Weeks 3-4)
3. **Phase 3**: Real-time Features (Weeks 5-6)
4. **Phase 4**: Advanced Features & Production (Weeks 7-8)

## Phase 1: Enhanced PubSub Infrastructure

### Goals

- Add support for Kafka, Redis, and NATS backends
- Implement configuration-driven backend selection
- Enhance error handling and retry mechanisms

### Tasks

#### 1.1 Extend PubSub Type System

```go
// Update internal/types/pubsub.go
const (
    MemoryPubSub PubSubType = "memory"
    KafkaPubSub  PubSubType = "kafka"
    RedisPubSub  PubSubType = "redis"
    NATSPubSub   PubSubType = "nats"
)
```

#### 1.2 Create PubSub Factory

Create `internal/pubsub/factory.go` to manage different backends:

```go
type Factory interface {
    CreatePubSub(backend types.PubSubType) PubSub
}
```

#### 1.3 Implement Kafka Backend

- Create `internal/pubsub/kafka/pubsub.go`
- Add Kafka configuration to config system
- Implement publisher and subscriber interfaces

#### 1.4 Implement Redis Backend

- Create `internal/pubsub/redis/pubsub.go`
- Use Redis Streams for message ordering
- Implement consumer groups

#### 1.5 Update Configuration

```yaml
# Add to config.yaml
pubsub:
  backend: "kafka" # memory, kafka, redis, nats

  kafka:
    brokers: ["localhost:9092"]
    consumer_group: "webhook-service"
    topics:
      events: "app.events"
      deadletter: "app.deadletter"

  redis:
    url: "redis://localhost:6379"
    streams:
      events: "app:events"
```

#### 1.6 Update Module Dependencies

Modify `internal/webhook/module.go` to use factory pattern.

### Deliverables

- [x] Multiple PubSub backend support
- [x] Configuration-driven backend selection
- [x] Updated module system
- [x] Basic integration tests

## Phase 2: External Integration

### Goals

- Implement Razorpay webhook receiver
- Create event transformation system
- Add external service publishers

### Tasks

#### 2.1 Create Webhook Receivers

Create `internal/webhook/receivers/` directory structure:

```
receivers/
├── razorpay.go
├── interface.go
└── registry.go
```

#### 2.2 Implement Razorpay Receiver

- Signature validation using HMAC-SHA256
- Payload transformation to internal event format
- Error handling and logging

#### 2.3 Create Event Types

Extend `internal/types/events.go`:

```go
type RazorpayPaymentEvent struct {
    Event   string `json:"event"`
    Entity  string `json:"entity"`
    Payment PaymentData `json:"payment"`
}
```

#### 2.4 Add HTTP Endpoints

Create webhook endpoints in your API layer:

```go
// POST /api/webhooks/razorpay
func (h *WebhookHandler) handleRazorpayWebhook(c *gin.Context)
```

#### 2.5 Implement Event Handlers

Create specialized handlers for different event types:

- `internal/webhook/handlers/razorpay_handler.go`
- `internal/webhook/handlers/payment_handler.go`

#### 2.6 Update Router Configuration

Register new handlers with the event router.

### Deliverables

- [x] Razorpay webhook integration
- [x] Event transformation system
- [x] HTTP webhook endpoints
- [x] Specialized event handlers
- [x] Integration with existing system

## Phase 3: Real-time Features

### Goals

- Implement WebSocket gateway
- Add Server-Sent Events support
- Create real-time event broadcasting

### Tasks

#### 3.1 Create Real-time Module

Create `internal/realtime/` directory:

```
realtime/
├── websocket.go
├── sse.go
├── gateway.go
└── module.go
```

#### 3.2 Implement WebSocket Gateway

- Connection management
- Authentication integration
- Event subscription handling
- Broadcast mechanisms

#### 3.3 Add SSE Support

Alternative to WebSocket for browser compatibility:

```go
func (h *SSEHandler) StreamEvents(w http.ResponseWriter, r *http.Request)
```

#### 3.4 Create Real-time Handler

- Subscribe to relevant event topics
- Filter events by user/permissions
- Broadcast to connected clients

#### 3.5 Frontend Integration

Provide JavaScript client library:

```javascript
class EventWebSocket {
    connect()
    subscribe(eventTypes)
    on(eventPattern, handler)
}
```

#### 3.6 Update Event Router

Route events to real-time handlers alongside existing handlers.

### Deliverables

- [x] WebSocket gateway implementation
- [x] Server-Sent Events support
- [x] Real-time event broadcasting
- [x] Frontend integration library
- [x] User-specific event filtering

## Phase 4: Advanced Features & Production

### Goals

- Implement event sourcing
- Add monitoring and observability
- Production deployment preparation

### Tasks

#### 4.1 Event Sourcing Implementation

Create `internal/eventsource/` module:

- Event store interface
- Event replay capabilities
- Aggregate rebuilding

#### 4.2 Monitoring Integration

- Prometheus metrics
- Health check endpoints
- Performance monitoring
- Alert configurations

#### 4.3 Circuit Breaker Pattern

Implement circuit breakers for external service calls:

```go
type CircuitBreaker interface {
    Execute(ctx context.Context, fn func() error) error
    GetState() State
}
```

#### 4.4 Outbox Pattern

Implement transactional outbox for reliable event publishing:

- Database outbox table
- Background publisher service
- Idempotency handling

#### 4.5 Security Enhancements

- Rate limiting implementation
- IP allowlisting
- Request size limits
- Audit logging

#### 4.6 Production Deployment

- Docker containerization
- Kubernetes manifests
- CI/CD pipeline setup
- Load testing

### Deliverables

- [x] Event sourcing system
- [x] Comprehensive monitoring
- [x] Circuit breaker implementation
- [x] Outbox pattern for reliability
- [x] Security hardening
- [x] Production deployment setup

## Implementation Checklist

### Phase 1 Tasks

- [ ] Create PubSub factory interface
- [ ] Implement Kafka backend
- [ ] Implement Redis backend
- [ ] Implement NATS backend
- [ ] Update configuration system
- [ ] Update module dependencies
- [ ] Write backend integration tests
- [ ] Update documentation

### Phase 2 Tasks

- [ ] Create webhook receiver interface
- [ ] Implement Razorpay receiver
- [ ] Add signature validation
- [ ] Create event transformation logic
- [ ] Add HTTP endpoints
- [ ] Implement specialized handlers
- [ ] Update router configuration
- [ ] Write integration tests

### Phase 3 Tasks

- [ ] Create WebSocket gateway
- [ ] Implement connection management
- [ ] Add authentication to WS
- [ ] Create SSE handler
- [ ] Implement real-time broadcasting
- [ ] Create frontend client library
- [ ] Add user-specific filtering
- [ ] Write real-time tests

### Phase 4 Tasks

- [ ] Implement event store
- [ ] Add event replay functionality
- [ ] Create monitoring dashboards
- [ ] Implement circuit breakers
- [ ] Add outbox pattern
- [ ] Implement rate limiting
- [ ] Create deployment manifests
- [ ] Set up CI/CD pipeline
- [ ] Conduct load testing
- [ ] Security audit

## Testing Strategy

### Unit Tests

- Each component should have >80% test coverage
- Mock external dependencies
- Test error scenarios

### Integration Tests

- End-to-end event flow testing
- Multiple backend compatibility
- Real external service testing (with mocks)

### Load Tests

- High-volume event processing
- Concurrent connection handling
- Memory and CPU usage under load

### Security Tests

- Webhook signature validation
- Authentication bypass attempts
- Rate limiting effectiveness

## Migration Strategy

### Phase 1 Migration

1. Deploy with memory backend (no changes to existing behavior)
2. Configure Kafka/Redis in parallel
3. Switch backend via configuration
4. Monitor for any issues

### Phase 2 Migration

1. Deploy webhook receivers (inactive initially)
2. Configure Razorpay webhook endpoints
3. Test with Razorpay webhook simulator
4. Enable production webhooks

### Phase 3 Migration

1. Deploy WebSocket gateway
2. Test with limited user base
3. Gradually roll out to all users
4. Monitor connection stability

### Phase 4 Migration

1. Deploy monitoring first
2. Enable circuit breakers gradually
3. Implement outbox pattern with feature flags
4. Production deployment with blue-green strategy

## Risk Mitigation

### High-Risk Items

1. **Message Loss**: Implement outbox pattern and persistent queues
2. **Performance Degradation**: Load testing and gradual rollout
3. **Security Vulnerabilities**: Security audit and penetration testing
4. **Data Consistency**: Transactional outbox and idempotency

### Monitoring Alerts

- High error rates
- Message processing delays
- WebSocket connection drops
- External service failures

## Success Metrics

### Technical Metrics

- Message throughput: >10,000 events/second
- Processing latency: <100ms p95
- WebSocket connection stability: >99% uptime
- Error rate: <0.1%

### Business Metrics

- Webhook delivery success rate: >99.9%
- Real-time notification delivery: <1 second
- System availability: >99.9%
- User engagement with real-time features

This roadmap provides a structured approach to implementing your enhanced event-driven system while minimizing risks and ensuring production readiness.
