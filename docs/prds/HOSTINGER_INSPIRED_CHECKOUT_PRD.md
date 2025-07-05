# 🚀 **HOSTINGER-INSPIRED ROBUST CHECKOUT & ENROLLMENT FLOW PRD**

## 📋 **Executive Summary**

Based on detailed analysis of Hostinger's production checkout flow, this PRD outlines a robust, scalable, and error-resistant e-commerce checkout system for the Interns Go-Backend platform. The system incorporates proven patterns from Hostinger's architecture while adapting them for internship enrollment use cases.

### 🎯 **Vision Statement**

Create a bulletproof checkout experience that matches Hostinger's reliability and scalability, providing users with a seamless journey from internship discovery to successful enrollment, while maintaining enterprise-grade fault tolerance, security, and performance.

---

## 🔍 **Hostinger Flow Analysis & Key Insights**

### **Flow Breakdown Analysis**

#### **1. Cart Creation Flow**

```
POST /api-proxy/api/cart
├── Sale slug validation
├── Product configuration
├── Period selection (billing cycle)
├── Analytics data collection
└── Cart token generation
```

**Key Insights:**

- **Separate Cart Domain**: Dedicated cart.hostinger.com subdomain
- **Token-Based Cart**: UUID-based cart tokens for stateless operations
- **Analytics Integration**: Rich analytics data collection throughout flow
- **Sale Integration**: Dynamic sale/promotion system

#### **2. Authentication & Authorization**

```
POST /api/v1/cart/{cart_token}/auth/check
├── JWT token validation
├── User session verification
├── Cart ownership validation
└── Permission checks
```

**Key Insights:**

- **Cart-Level Auth**: Authentication at cart level, not just user level
- **Session Management**: Robust session handling with device tracking
- **Correlation IDs**: Request tracing for debugging and monitoring

#### **3. Pricing & Estimation**

```
POST /api/v1/cart/estimate
├── Cart token validation
├── Tax calculation (IGST 18%)
├── Discount application
├── Currency conversion
└── Final amount calculation
```

**Key Insights:**

- **Real-Time Pricing**: Dynamic pricing with tax and discount calculations
- **Currency Handling**: Multi-currency support with conversion rates
- **Tax Integration**: Automated tax calculation based on location
- **Discount Transparency**: Clear breakdown of applied discounts

#### **4. Order Creation & Payment**

```
POST /api/v1/order
├── Cart validation
├── Order token generation
├── Payment token creation
├── Razorpay integration
└── Payment session initiation
```

**Key Insights:**

- **Dual Token System**: Separate order and payment tokens
- **Payment Gateway Integration**: Direct Razorpay integration
- **Session Management**: Payment session with token-based security

---

## 🏗️ **Robust Architecture Design**

### **System Architecture**

```
┌─────────────────────────────────────────────────────────────────┐
│                        Load Balancer                            │
└─────────────────────┬───────────────────────────────────────────┘
                      │
         ┌────────────┼────────────┐
         │            │            │
┌────────▼─────┐ ┌────▼────┐ ┌────▼────┐
│  Main API    │ │  Cart   │ │ Payment │
│  Gateway     │ │ Service │ │ Gateway │
└──────────────┘ └─────────┘ └─────────┘
         │            │            │
         └────────────┼────────────┘
                      │
         ┌────────────▼────────────┐
         │     PostgreSQL DB      │
         │   (Primary + Replica)  │
         └─────────────────────────┘
                      │
         ┌────────────▼────────────┐
         │    Redis Cluster       │
         │   (Session + Cache)    │
         └─────────────────────────┘
                      │
         ┌────────────▼────────────┐
         │   Message Queue        │
         │   (RabbitMQ/Kafka)     │
         └─────────────────────────┘
```

### **Service Separation Strategy**

#### **1. Cart Service (cart.interns.com)**

- **Domain**: Dedicated subdomain for cart operations
- **Responsibilities**: Cart management, pricing, discounts
- **Data Store**: Redis for cart data, PostgreSQL for persistence
- **Scaling**: Horizontal scaling with load balancing

#### **2. Payment Service (pay.interns.com)**

- **Domain**: Dedicated subdomain for payment operations
- **Responsibilities**: Payment processing, gateway integration
- **Security**: PCI DSS compliance, token-based security
- **Isolation**: Complete isolation from main application

#### **3. Main API Gateway**

- **Responsibilities**: Authentication, routing, rate limiting
- **Security**: JWT validation, CORS, request sanitization
- **Monitoring**: Request tracing, error tracking

---

## 🔧 **Robust Implementation Requirements**

### **1. Cart Management System**

#### **Cart Token Architecture**

```go
type Cart struct {
    Token           string          `json:"token"`           // UUID-based cart token
    UserID          *string         `json:"user_id"`         // Optional, for guest carts
    Status          CartStatus      `json:"status"`          // draft, active, expired, converted
    Items           []CartItem      `json:"items"`
    Subtotal        decimal.Decimal `json:"subtotal"`
    Total           decimal.Decimal `json:"total"`
    Currency        string          `json:"currency"`
    TaxRate         decimal.Decimal `json:"tax_rate"`
    TaxAmount       decimal.Decimal `json:"tax_amount"`
    DiscountCode    *string         `json:"discount_code"`
    DiscountAmount  decimal.Decimal `json:"discount_amount"`
    ExpiresAt       time.Time       `json:"expires_at"`
    DeviceID        string          `json:"device_id"`       // For guest tracking
    AnalyticsData   []AnalyticsData `json:"analytics_data"`
    CreatedAt       time.Time       `json:"created_at"`
    UpdatedAt       time.Time       `json:"updated_at"`
}

type CartItem struct {
    InternshipID    string          `json:"internship_id"`
    Title           string          `json:"title"`
    Price           decimal.Decimal `json:"price"`
    Quantity        int             `json:"quantity"`
    Period          *Period         `json:"period,omitempty"`     // Billing period
    Metadata        map[string]any  `json:"metadata,omitempty"`
}

type Period struct {
    Value int    `json:"value"` // Number of months/years
    Unit  string `json:"unit"`  // "month", "year"
}
```

#### **API Endpoints (Inspired by Hostinger)**

```
# Cart Creation & Management
POST   /api-proxy/api/cart                    # Create cart with products
GET    /api/v1/cart/{token}                   # Get cart details
PUT    /api/v1/cart/{token}/plan/{plan_id}    # Update cart plan
DELETE /api/v1/cart/{token}/item/{item_id}    # Remove item

# Authentication & Authorization
POST   /api/v1/cart/{token}/auth/check        # Check cart authentication
POST   /api/v1/cart/{token}/auth/login        # Login to cart

# Pricing & Estimation
POST   /api/v1/cart/estimate                  # Get cart pricing
GET    /api/v1/cart/{token}/discounts         # Get available discounts

# Discount Management
POST   /api/v1/cart/{token}/coupon            # Apply coupon code
DELETE /api/v1/cart/{token}/coupon            # Remove coupon

# Order Creation
POST   /api/v1/order                          # Create order from cart
GET    /api/v1/order/{order_token}            # Get order details
```

### **2. Robust Error Handling & Recovery**

#### **Error Categories & Handling**

```go
type ErrorResponse struct {
    Status  int    `json:"status"`
    Success bool   `json:"success"`
    Error   *Error `json:"error,omitempty"`
    Data    any    `json:"data,omitempty"`
}

type Error struct {
    Code                int      `json:"code"`
    Message             string   `json:"message"`
    ValidationMessages  []string `json:"validation_messages,omitempty"`
    CorrelationID       string   `json:"correlation_id"`
    Retryable           bool     `json:"retryable"`
    SuggestedAction     string   `json:"suggested_action,omitempty"`
}

// Error Codes (Inspired by Hostinger)
const (
    ErrCartNotFound           = 1001
    ErrCartExpired           = 1002
    ErrInvalidCoupon         = 2004
    ErrPaymentFailed         = 3001
    ErrInsufficientBalance   = 3002
    ErrGatewayTimeout        = 4001
    ErrRateLimitExceeded     = 429
)
```

#### **Retry Mechanisms**

```go
type RetryConfig struct {
    MaxAttempts     int           `json:"max_attempts"`
    InitialDelay    time.Duration `json:"initial_delay"`
    MaxDelay        time.Duration `json:"max_delay"`
    BackoffFactor   float64       `json:"backoff_factor"`
    RetryableErrors []int         `json:"retryable_errors"`
}

// Retry Strategy
func (s *CartService) CreateCartWithRetry(ctx context.Context, req *CreateCartRequest) (*CartResponse, error) {
    var lastErr error
    for attempt := 1; attempt <= s.config.MaxAttempts; attempt++ {
        cart, err := s.createCart(ctx, req)
        if err == nil {
            return cart, nil
        }

        if !s.isRetryableError(err) {
            return nil, err
        }

        lastErr = err
        delay := s.calculateBackoffDelay(attempt)
        time.Sleep(delay)
    }

    return nil, fmt.Errorf("max retry attempts exceeded: %w", lastErr)
}
```

### **3. Scalable Data Architecture**

#### **Database Design**

```sql
-- Cart Table (Optimized for performance)
CREATE TABLE carts (
    token VARCHAR(36) PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    status cart_status DEFAULT 'draft',
    items JSONB NOT NULL DEFAULT '[]',
    subtotal DECIMAL(10,2) NOT NULL DEFAULT 0,
    total DECIMAL(10,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'INR',
    tax_rate DECIMAL(5,2) NOT NULL DEFAULT 0,
    tax_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    discount_code VARCHAR(50),
    discount_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    expires_at TIMESTAMP NOT NULL,
    device_id VARCHAR(255),
    analytics_data JSONB DEFAULT '[]',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    -- Indexes for performance
    INDEX idx_carts_user_id (user_id),
    INDEX idx_carts_status (status),
    INDEX idx_carts_expires_at (expires_at),
    INDEX idx_carts_device_id (device_id)
);

-- Order Table (Separate from cart)
CREATE TABLE orders (
    order_token VARCHAR(36) PRIMARY KEY,
    cart_token VARCHAR(36) NOT NULL REFERENCES carts(token),
    user_id UUID NOT NULL REFERENCES users(id),
    status order_status DEFAULT 'pending',
    currency VARCHAR(3) NOT NULL DEFAULT 'INR',
    subtotal DECIMAL(10,2) NOT NULL,
    total DECIMAL(10,2) NOT NULL,
    tax_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    discount_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    payment_token VARCHAR(36),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    INDEX idx_orders_cart_token (cart_token),
    INDEX idx_orders_user_id (user_id),
    INDEX idx_orders_status (status),
    INDEX idx_orders_payment_token (payment_token)
);

-- Payment Tokens Table
CREATE TABLE payment_tokens (
    token VARCHAR(36) PRIMARY KEY,
    order_token VARCHAR(36) NOT NULL REFERENCES orders(order_token),
    gateway VARCHAR(50) NOT NULL,
    gateway_order_id VARCHAR(255),
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    status payment_status DEFAULT 'pending',
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    INDEX idx_payment_tokens_order_token (order_token),
    INDEX idx_payment_tokens_status (status),
    INDEX idx_payment_tokens_expires_at (expires_at)
);
```

#### **Caching Strategy**

```go
type CacheConfig struct {
    CartTTL        time.Duration `json:"cart_ttl"`        // 24 hours
    PricingTTL     time.Duration `json:"pricing_ttl"`     // 5 minutes
    UserSessionTTL time.Duration `json:"user_session_ttl"` // 30 minutes
    DiscountTTL    time.Duration `json:"discount_ttl"`    // 1 hour
}

// Redis Cache Keys
const (
    CartKeyPrefix     = "cart:"
    PricingKeyPrefix  = "pricing:"
    SessionKeyPrefix  = "session:"
    DiscountKeyPrefix = "discount:"
)

func (s *CartService) GetCartWithCache(ctx context.Context, token string) (*Cart, error) {
    // Try cache first
    cacheKey := CartKeyPrefix + token
    if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
        var cart Cart
        if err := json.Unmarshal([]byte(cached), &cart); err == nil {
            return &cart, nil
        }
    }

    // Fallback to database
    cart, err := s.repo.GetByToken(ctx, token)
    if err != nil {
        return nil, err
    }

    // Cache the result
    if cartData, err := json.Marshal(cart); err == nil {
        s.cache.Set(ctx, cacheKey, string(cartData), s.config.CartTTL)
    }

    return cart, nil
}
```

### **4. Payment Integration (Razorpay-Inspired)**

#### **Payment Flow Architecture**

```go
type PaymentService struct {
    razorpayClient *razorpay.Client
    stripeClient   *stripe.Client
    cache          cache.Cache
    logger         *logger.Logger
}

type PaymentRequest struct {
    OrderToken     string          `json:"order_token"`
    Amount         decimal.Decimal `json:"amount"`
    Currency       string          `json:"currency"`
    PaymentMethod  string          `json:"payment_method"`
    UPIID          *string         `json:"upi_id,omitempty"`
    CardToken      *string         `json:"card_token,omitempty"`
    ReturnURL      string          `json:"return_url"`
    CancelURL      string          `json:"cancel_url"`
}

type PaymentResponse struct {
    PaymentToken   string          `json:"payment_token"`
    GatewayOrderID string          `json:"gateway_order_id"`
    CheckoutURL    string          `json:"checkout_url"`
    Status         string          `json:"status"`
    ExpiresAt      time.Time       `json:"expires_at"`
}

func (s *PaymentService) CreatePaymentSession(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
    // Create payment token
    paymentToken := uuid.New().String()

    // Create gateway order
    gatewayOrder, err := s.createGatewayOrder(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("gateway order creation failed: %w", err)
    }

    // Store payment token
    paymentTokenData := &PaymentToken{
        Token:          paymentToken,
        OrderToken:     req.OrderToken,
        Gateway:        "razorpay",
        GatewayOrderID: gatewayOrder.ID,
        Amount:         req.Amount,
        Currency:       req.Currency,
        Status:         "pending",
        ExpiresAt:      time.Now().Add(30 * time.Minute),
    }

    if err := s.repo.CreatePaymentToken(ctx, paymentTokenData); err != nil {
        return nil, err
    }

    return &PaymentResponse{
        PaymentToken:   paymentToken,
        GatewayOrderID: gatewayOrder.ID,
        CheckoutURL:    gatewayOrder.CheckoutURL,
        Status:         "pending",
        ExpiresAt:      paymentTokenData.ExpiresAt,
    }, nil
}
```

### **5. Webhook & Event System**

#### **Robust Webhook Handling**

```go
type WebhookHandler struct {
    signatureVerifier SignatureVerifier
    eventProcessor    EventProcessor
    retryQueue        queue.Queue
    logger            *logger.Logger
}

func (h *WebhookHandler) HandleRazorpayWebhook(ctx context.Context, payload []byte, headers map[string]string) error {
    // Verify signature
    if err := h.signatureVerifier.Verify(payload, headers["X-Razorpay-Signature"]); err != nil {
        h.logger.Errorw("webhook signature verification failed", "error", err)
        return http.StatusBadRequest
    }

    // Parse event
    var event RazorpayEvent
    if err := json.Unmarshal(payload, &event); err != nil {
        h.logger.Errorw("webhook payload parsing failed", "error", err)
        return http.StatusBadRequest
    }

    // Process event with retry
    if err := h.processEventWithRetry(ctx, &event); err != nil {
        h.logger.Errorw("webhook event processing failed", "event_id", event.ID, "error", err)

        // Queue for retry
        h.retryQueue.Publish(ctx, &RetryMessage{
            EventID:   event.ID,
            Payload:   payload,
            Attempts:  1,
            NextRetry: time.Now().Add(5 * time.Minute),
        })

        return http.StatusInternalServerError
    }

    return http.StatusOK
}

func (h *WebhookHandler) processEventWithRetry(ctx context.Context, event *RazorpayEvent) error {
    // Check if already processed
    if h.isEventProcessed(ctx, event.ID) {
        return nil // Idempotent
    }

    // Process based on event type
    switch event.Event {
    case "payment.captured":
        return h.handlePaymentSuccess(ctx, event)
    case "payment.failed":
        return h.handlePaymentFailure(ctx, event)
    case "order.paid":
        return h.handleOrderPaid(ctx, event)
    default:
        h.logger.Warnw("unknown webhook event type", "event_type", event.Event)
        return nil
    }
}
```

### **6. Monitoring & Observability**

#### **Comprehensive Monitoring**

```go
type Metrics struct {
    CartCreated        prometheus.Counter
    CartConverted      prometheus.Counter
    PaymentSuccess     prometheus.Counter
    PaymentFailed      prometheus.Counter
    WebhookProcessed   prometheus.Counter
    WebhookFailed      prometheus.Counter
    ResponseTime       prometheus.Histogram
    ErrorRate          prometheus.Counter
}

type Tracing struct {
    tracer trace.Tracer
}

func (s *CartService) CreateCart(ctx context.Context, req *CreateCartRequest) (*CartResponse, error) {
    start := time.Now()
    defer func() {
        s.metrics.ResponseTime.Observe(time.Since(start).Seconds())
    }()

    // Create span for tracing
    ctx, span := s.tracing.tracer.Start(ctx, "cart.create")
    defer span.End()

    // Add correlation ID
    correlationID := uuid.New().String()
    ctx = context.WithValue(ctx, "correlation_id", correlationID)

    // Process request
    cart, err := s.processCreateCart(ctx, req)
    if err != nil {
        s.metrics.ErrorRate.Inc()
        s.logger.Errorw("cart creation failed",
            "correlation_id", correlationID,
            "error", err,
            "user_id", req.UserID)
        return nil, err
    }

    s.metrics.CartCreated.Inc()
    s.logger.Infow("cart created successfully",
        "correlation_id", correlationID,
        "cart_token", cart.Token,
        "user_id", req.UserID)

    return cart, nil
}
```

---

## 🚀 **Implementation Roadmap (1 Day)**

### **Phase 1: Foundation (Hours 1-4)**

#### **Hour 1: Database Schema & Migrations**

- [ ] Create cart, order, payment_tokens tables
- [ ] Implement database indexes for performance
- [ ] Set up database connection pooling
- [ ] Create database migration scripts

#### **Hour 2: Domain Models & Repositories**

- [ ] Implement Cart, Order, PaymentToken domain models
- [ ] Create repository interfaces and implementations
- [ ] Add validation logic for domain models
- [ ] Implement caching layer with Redis

#### **Hour 3: Service Layer Foundation**

- [ ] Create CartService with basic CRUD operations
- [ ] Implement OrderService with order lifecycle
- [ ] Add PaymentService with gateway integration
- [ ] Set up error handling and retry mechanisms

#### **Hour 4: API Layer Scaffolding**

- [ ] Create API routes for cart operations
- [ ] Implement authentication middleware
- [ ] Add request validation and sanitization
- [ ] Set up CORS and security headers

### **Phase 2: Core Functionality (Hours 5-8)**

#### **Hour 5: Cart Management**

- [ ] Implement cart creation with token generation
- [ ] Add cart item management (add, remove, update)
- [ ] Implement cart pricing calculation
- [ ] Add cart expiration and cleanup

#### **Hour 6: Pricing & Discounts**

- [ ] Implement real-time pricing calculation
- [ ] Add tax calculation based on location
- [ ] Implement coupon code validation and application
- [ ] Add discount transparency and breakdown

#### **Hour 7: Order Creation**

- [ ] Implement order creation from cart
- [ ] Add order token generation
- [ ] Implement order status management
- [ ] Add order validation and security checks

#### **Hour 8: Payment Integration**

- [ ] Complete Razorpay provider implementation
- [ ] Add payment session creation
- [ ] Implement payment token management
- [ ] Add payment status tracking

### **Phase 3: Robustness & Reliability (Hours 9-12)**

#### **Hour 9: Error Handling & Recovery**

- [ ] Implement comprehensive error handling
- [ ] Add retry mechanisms with exponential backoff
- [ ] Implement circuit breaker pattern
- [ ] Add graceful degradation

#### **Hour 10: Webhook System**

- [ ] Implement webhook signature verification
- [ ] Add webhook event processing
- [ ] Implement idempotency for webhooks
- [ ] Add webhook retry queue

#### **Hour 11: Event-Driven Architecture**

- [ ] Set up message queue (RabbitMQ/Kafka)
- [ ] Implement event publishing
- [ ] Add event consumers for enrollment
- [ ] Implement dead letter queue

#### **Hour 12: Monitoring & Observability**

- [ ] Add comprehensive logging
- [ ] Implement metrics collection
- [ ] Add distributed tracing
- [ ] Set up health checks

### **Phase 4: Production Readiness (Hours 13-16)**

#### **Hour 13: Security Hardening**

- [ ] Implement rate limiting
- [ ] Add input validation and sanitization
- [ ] Implement CORS policies
- [ ] Add security headers

#### **Hour 14: Performance Optimization**

- [ ] Optimize database queries
- [ ] Implement connection pooling
- [ ] Add caching strategies
- [ ] Optimize API response times

#### **Hour 15: Testing & Validation**

- [ ] Write unit tests for core services
- [ ] Implement integration tests
- [ ] Add load testing scripts
- [ ] Test error scenarios

#### **Hour 16: Documentation & Deployment**

- [ ] Create API documentation
- [ ] Write deployment guides
- [ ] Add monitoring dashboards
- [ ] Prepare rollback procedures

---

## 🎯 **Success Criteria**

### **Functional Success**

- [ ] Cart creation and management works flawlessly
- [ ] Pricing calculation is accurate and real-time
- [ ] Payment processing succeeds with multiple gateways
- [ ] Webhook processing is reliable and idempotent
- [ ] Order-to-enrollment flow is automated

### **Performance Success**

- [ ] API response times < 200ms (P95)
- [ ] System handles 10,000+ concurrent users
- [ ] Payment success rate > 98%
- [ ] Webhook processing < 1s average

### **Reliability Success**

- [ ] 99.9% system uptime
- [ ] Zero data loss in payment processing
- [ ] Graceful handling of gateway failures
- [ ] Automatic recovery from transient errors

### **Security Success**

- [ ] PCI DSS compliance for payment data
- [ ] Secure webhook signature verification
- [ ] Protection against common attacks
- [ ] Audit trail for all transactions

---

## 🚨 **Risk Mitigation Strategies**

### **Technical Risks**

| Risk                     | Impact | Mitigation Strategy                                   |
| ------------------------ | ------ | ----------------------------------------------------- |
| Payment Gateway Downtime | High   | Multiple gateway fallback, circuit breaker            |
| Database Performance     | Medium | Query optimization, read replicas, caching            |
| Webhook Failures         | High   | Retry queue, dead letter queue, manual reconciliation |
| Memory Leaks             | Medium | Connection pooling, resource cleanup, monitoring      |

### **Business Risks**

| Risk               | Impact   | Mitigation Strategy                               |
| ------------------ | -------- | ------------------------------------------------- |
| Cart Abandonment   | High     | Cart persistence, email reminders, A/B testing    |
| Payment Failures   | High     | Multiple payment methods, clear error messages    |
| Data Loss          | Critical | Database backups, transaction logging, monitoring |
| Scalability Issues | Medium   | Horizontal scaling, load balancing, auto-scaling  |

---

## 📊 **Monitoring & Alerting**

### **Key Metrics to Track**

- Cart creation rate and conversion
- Payment success/failure rates
- API response times and error rates
- Webhook processing latency
- Database performance metrics
- Cache hit/miss ratios

### **Alerting Rules**

- Payment success rate < 95%
- API response time > 500ms
- Webhook processing failure rate > 5%
- Database connection pool exhaustion
- Cache miss rate > 20%

---

This PRD provides a comprehensive blueprint for implementing a robust, scalable, and error-resistant checkout system inspired by Hostinger's proven architecture while adapting it for internship enrollment use cases.
