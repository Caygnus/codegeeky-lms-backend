# Webhook System Technical Specification

## 1. Watermill Integration

### 1.1 Router Configuration
```go
// internal/pubsub/router/router.go
type Router struct {
    router *message.Router
    logger *logger.Logger
    config *config.Webhook
}

func NewRouter(cfg *config.Configuration, logger *logger.Logger) (*Router, error) {
    router, err := message.NewRouter(
        message.RouterConfig{},
        watermill.NewStdLogger(true, false),
    )
    
    // Add middleware
    router.AddMiddleware(
        middleware.PoisonQueue,
        middleware.Recoverer,
        middleware.CorrelationID,
        middleware.Retry{
            MaxRetries: cfg.Webhook.MaxRetries,
            InitialInterval: cfg.Webhook.InitialInterval,
            MaxInterval: cfg.Webhook.MaxInterval,
            Multiplier: cfg.Webhook.Multiplier,
            MaxElapsedTime: cfg.Webhook.MaxElapsedTime,
        },
    )
}
```

### 1.2 Message Processing
```go
// internal/webhook/handler/handler.go
func (h *handler) processMessage(msg *message.Message) error {
    var event types.WebhookEvent
    if err := json.Unmarshal(msg.Payload, &event); err != nil {
        return nil // Don't retry on unmarshal errors
    }

    // Get user config and validate
    userCfg, ok := h.config.Users[*event.UserID]
    if !ok || !userCfg.Enabled {
        return nil
    }

    // Build and send webhook
    builder, err := h.factory.GetBuilder(event.EventName)
    if err != nil {
        return err
    }

    webHookPayload, err := builder.BuildPayload(ctx, event.EventName, event.Payload)
    if err != nil {
        return err
    }

    // Send webhook
    req := &httpclient.Request{
        Method: "POST",
        URL: userCfg.Endpoint,
        Headers: userCfg.Headers,
        Body: webHookPayload,
    }

    resp, err := h.client.Send(ctx, req)
    if err != nil {
        return err
    }

    return nil
}
```

## 2. Payment Gateway Integration

### 2.1 Razorpay Webhook Handler
```go
// internal/webhook/handler/razorpay.go
type RazorpayHandler struct {
    service *service.PaymentService
    logger  *logger.Logger
}

func (h *RazorpayHandler) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
    // Verify signature
    if err := h.verifySignature(payload, signature); err != nil {
        return err
    }

    // Parse payload
    var event RazorpayEvent
    if err := json.Unmarshal(payload, &event); err != nil {
        return err
    }

    // Process based on event type
    switch event.Event {
    case "payment.authorized":
        return h.handlePaymentAuthorized(ctx, event)
    case "payment.failed":
        return h.handlePaymentFailed(ctx, event)
    // ... other event types
    }

    return nil
}
```

### 2.2 Stripe Webhook Handler
```go
// internal/webhook/handler/stripe.go
type StripeHandler struct {
    service *service.PaymentService
    logger  *logger.Logger
}

func (h *StripeHandler) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
    // Verify signature
    if err := h.verifySignature(payload, signature); err != nil {
        return err
    }

    // Parse payload
    var event stripe.Event
    if err := json.Unmarshal(payload, &event); err != nil {
        return err
    }

    // Process based on event type
    switch event.Type {
    case "payment_intent.succeeded":
        return h.handlePaymentSucceeded(ctx, event)
    case "payment_intent.payment_failed":
        return h.handlePaymentFailed(ctx, event)
    // ... other event types
    }

    return nil
}
```

## 3. Database Operations

### 3.1 Webhook Event Storage
```go
// internal/webhook/repository/event.go
type EventRepository interface {
    Create(ctx context.Context, event *types.WebhookEvent) error
    UpdateStatus(ctx context.Context, id string, status string) error
    GetByID(ctx context.Context, id string) (*types.WebhookEvent, error)
    ListByUserID(ctx context.Context, userID string, limit, offset int) ([]*types.WebhookEvent, error)
}

type eventRepository struct {
    db *sql.DB
}

func (r *eventRepository) Create(ctx context.Context, event *types.WebhookEvent) error {
    query := `
        INSERT INTO webhook_events (
            id, user_id, event_name, payload, status, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
    // ... implementation
}
```

### 3.2 Webhook Delivery Tracking
```go
// internal/webhook/repository/delivery.go
type DeliveryRepository interface {
    Create(ctx context.Context, delivery *types.WebhookDelivery) error
    UpdateStatus(ctx context.Context, id string, status string, responseCode int, responseBody string) error
    GetByEventID(ctx context.Context, eventID string) ([]*types.WebhookDelivery, error)
}

type deliveryRepository struct {
    db *sql.DB
}

func (r *deliveryRepository) Create(ctx context.Context, delivery *types.WebhookDelivery) error {
    query := `
        INSERT INTO webhook_deliveries (
            id, event_id, endpoint, status, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6)
    `
    // ... implementation
}
```

## 4. Client Response Handling

### 4.1 WebSocket Response
```go
// internal/webhook/response/websocket.go
type WebSocketResponse struct {
    hub *websocket.Hub
}

func (r *WebSocketResponse) Send(ctx context.Context, userID string, event *types.WebhookEvent) error {
    message := &websocket.Message{
        Type: "webhook",
        Data: event,
    }
    return r.hub.SendToUser(userID, message)
}
```

### 4.2 PubSub Response
```go
// internal/webhook/response/pubsub.go
type PubSubResponse struct {
    publisher publisher.WebhookPublisher
}

func (r *PubSubResponse) Send(ctx context.Context, event *types.WebhookEvent) error {
    return r.publisher.Publish(ctx, event)
}
```

## 5. Error Handling

### 5.1 Retry Mechanism
```go
// internal/webhook/retry/retry.go
type RetryConfig struct {
    MaxRetries      int
    InitialInterval time.Duration
    MaxInterval     time.Duration
    Multiplier      float64
}

func WithRetry(config RetryConfig, fn func() error) error {
    var err error
    interval := config.InitialInterval

    for i := 0; i < config.MaxRetries; i++ {
        if err = fn(); err == nil {
            return nil
        }

        time.Sleep(interval)
        interval = time.Duration(float64(interval) * config.Multiplier)
        if interval > config.MaxInterval {
            interval = config.MaxInterval
        }
    }

    return err
}
```

### 5.2 Dead Letter Queue
```go
// internal/webhook/dlq/dlq.go
type DeadLetterQueue struct {
    pubsub pubsub.PubSub
    logger *logger.Logger
}

func (q *DeadLetterQueue) HandleFailedMessage(ctx context.Context, msg *message.Message, err error) error {
    // Store failed message in DLQ
    failedMsg := &types.FailedMessage{
        MessageID: msg.UUID,
        Payload:   msg.Payload,
        Error:     err.Error(),
        Timestamp: time.Now(),
    }

    return q.pubsub.Publish(ctx, "dlq", failedMsg)
}
```

## 6. Monitoring

### 6.1 Metrics Collection
```go
// internal/webhook/metrics/metrics.go
type Metrics struct {
    deliverySuccess prometheus.Counter
    deliveryFailure prometheus.Counter
    deliveryLatency prometheus.Histogram
    retryCount      prometheus.Counter
}

func (m *Metrics) RecordDelivery(success bool, latency time.Duration) {
    if success {
        m.deliverySuccess.Inc()
    } else {
        m.deliveryFailure.Inc()
    }
    m.deliveryLatency.Observe(latency.Seconds())
}

func (m *Metrics) RecordRetry() {
    m.retryCount.Inc()
}
```

### 6.2 Health Checks
```go
// internal/webhook/health/health.go
type HealthChecker struct {
    db        *sql.DB
    pubsub    pubsub.PubSub
    logger    *logger.Logger
}

func (h *HealthChecker) Check(ctx context.Context) error {
    // Check database connection
    if err := h.db.PingContext(ctx); err != nil {
        return err
    }

    // Check pubsub connection
    if err := h.pubsub.Ping(ctx); err != nil {
        return err
    }

    return nil
}
```

## 7. Testing

### 7.1 Unit Tests
```go
// internal/webhook/handler/handler_test.go
func TestWebhookHandler_ProcessMessage(t *testing.T) {
    tests := []struct {
        name    string
        payload []byte
        wantErr bool
    }{
        {
            name: "valid payload",
            payload: []byte(`{"event":"payment.success","user_id":"123"}`),
            wantErr: false,
        },
        {
            name: "invalid payload",
            payload: []byte(`invalid json`),
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            h := NewHandler(mockPubSub, mockConfig, mockFactory, mockClient, mockLogger)
            msg := message.NewMessage(watermill.NewUUID(), tt.payload)
            err := h.processMessage(msg)
            if (err != nil) != tt.wantErr {
                t.Errorf("processMessage() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### 7.2 Integration Tests
```go
// internal/webhook/integration_test.go
func TestWebhookFlow(t *testing.T) {
    // Setup test environment
    ctx := context.Background()
    db := setupTestDB(t)
    pubsub := setupTestPubSub(t)
    handler := setupTestHandler(t, db, pubsub)

    // Create test event
    event := &types.WebhookEvent{
        ID:        uuid.New(),
        UserID:    "test-user",
        EventName: "payment.success",
        Payload:   []byte(`{"amount":1000}`),
    }

    // Process event
    err := handler.ProcessEvent(ctx, event)
    require.NoError(t, err)

    // Verify database state
    stored, err := db.GetWebhookEvent(ctx, event.ID)
    require.NoError(t, err)
    assert.Equal(t, event.UserID, stored.UserID)

    // Verify pubsub message
    msg := <-pubsub.Messages
    assert.Equal(t, event.EventName, msg.EventName)
} 