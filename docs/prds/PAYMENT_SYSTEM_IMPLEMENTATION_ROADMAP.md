# 🗺️ Payment System Implementation Roadmap

## 📋 Overview

This document provides a detailed implementation roadmap for the comprehensive internship enrollment workflow with payment integration system. The roadmap is structured in phases with specific milestones, deliverables, and timelines.

---

## 🚀 Implementation Phases

### **Phase 1: Foundation & Core Architecture (Weeks 1-4)**

#### **Week 1: Project Setup & Architecture Design**

**Objectives:**

- Set up development environment and project structure
- Design system architecture and database schema
- Set up CI/CD pipeline and monitoring

**Deliverables:**

- [ ] Project repository setup with Clean Architecture structure
- [ ] Database schema design and migrations
- [ ] Docker development environment
- [ ] CI/CD pipeline configuration
- [ ] Monitoring and logging setup

**Technical Tasks:**

```bash
# Project structure
mkdir -p internal/{payment,enrollment,billing,notification,webhook}
mkdir -p internal/payment/{adapters,factory,manager}
mkdir -p internal/domain/{payment,enrollment,discount}
mkdir -p cmd/{server,migrate}
mkdir -p docs/{api,architecture}
```

**Database Schema Implementation:**

```sql
-- Core payment tables
CREATE TABLE payment_gateways (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,
    type VARCHAR(50) NOT NULL,
    status gateway_status DEFAULT 'active',
    config JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE payment_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gateway_id UUID REFERENCES payment_gateways(id),
    method_type VARCHAR(50) NOT NULL,
    method_name VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    config JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);
```

#### **Week 2: Payment Gateway Interface & Core Adapters**

**Objectives:**

- Implement payment gateway interface and core adapters
- Set up gateway factory pattern
- Implement basic payment processing flow

**Deliverables:**

- [ ] Payment gateway interface definition
- [ ] Stripe gateway adapter
- [ ] Razorpay gateway adapter
- [ ] Gateway factory implementation
- [ ] Basic payment processing service

**Code Implementation:**

```go
// Payment Gateway Interface
type PaymentGateway interface {
    Name() string
    ProcessPayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error)
    ProcessRefund(ctx context.Context, req *RefundRequest) (*RefundResponse, error)
    ValidateWebhook(payload []byte, signature string) (*WebhookEvent, error)
    GetSupportedMethods() []PaymentMethodType
}

// Gateway Manager
type GatewayManager struct {
    gateways map[string]PaymentGateway
    selector GatewaySelector
    monitor  HealthMonitor
}

func (m *GatewayManager) ProcessPayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
    gateway, err := m.selector.SelectGateway(req)
    if err != nil {
        return nil, fmt.Errorf("gateway selection failed: %w", err)
    }

    return gateway.ProcessPayment(ctx, req)
}
```

#### **Week 3: Enrollment Service & Payment Integration**

**Objectives:**

- Implement enrollment service
- Integrate payment processing with enrollments
- Set up webhook handling system

**Deliverables:**

- [ ] Enrollment service implementation
- [ ] Payment-enrollment integration
- [ ] Webhook processing system
- [ ] Transaction management
- [ ] Event publishing system

**Service Implementation:**

```go
type EnrollmentService struct {
    paymentManager *payment.GatewayManager
    enrollRepo     enrollment.Repository
    eventPublisher events.Publisher
    logger         logger.Logger
}

func (s *EnrollmentService) ProcessEnrollment(ctx context.Context, req *EnrollmentRequest) (*EnrollmentResponse, error) {
    // Create enrollment record
    enrollment, err := s.enrollRepo.Create(ctx, &enrollment.CreateRequest{
        UserID:       req.UserID,
        InternshipID: req.InternshipID,
        Status:       enrollment.StatusPending,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create enrollment: %w", err)
    }

    // Process payment
    paymentReq := &payment.PaymentRequest{
        Amount:      req.Amount,
        Currency:    req.Currency,
        UserID:      req.UserID,
        Description: fmt.Sprintf("Enrollment for internship %s", req.InternshipID),
        Metadata: map[string]string{
            "enrollment_id": enrollment.ID,
            "type":         "enrollment",
        },
    }

    paymentResp, err := s.paymentManager.ProcessPayment(ctx, paymentReq)
    if err != nil {
        // Update enrollment status to failed
        s.enrollRepo.UpdateStatus(ctx, enrollment.ID, enrollment.StatusFailed)
        return nil, fmt.Errorf("payment processing failed: %w", err)
    }

    // Update enrollment with payment info
    enrollment.PaymentID = &paymentResp.TransactionID
    enrollment.Status = enrollment.StatusPaid

    if err := s.enrollRepo.Update(ctx, enrollment); err != nil {
        return nil, fmt.Errorf("failed to update enrollment: %w", err)
    }

    // Publish enrollment event
    s.eventPublisher.Publish(ctx, &events.EnrollmentCompleted{
        EnrollmentID: enrollment.ID,
        UserID:      req.UserID,
        Amount:      req.Amount,
    })

    return &EnrollmentResponse{
        EnrollmentID:  enrollment.ID,
        PaymentID:     paymentResp.TransactionID,
        Status:        string(enrollment.Status),
        RedirectURL:   paymentResp.RedirectURL,
    }, nil
}
```

#### **Week 4: API Layer & Testing**

**Objectives:**

- Implement REST API endpoints
- Set up comprehensive testing
- API documentation

**Deliverables:**

- [ ] REST API endpoints for enrollment
- [ ] API middleware (auth, validation, rate limiting)
- [ ] Unit tests for all services
- [ ] Integration tests
- [ ] API documentation

**API Implementation:**

```go
type EnrollmentHandler struct {
    enrollmentSvc *service.EnrollmentService
    validator     *validator.Validator
    logger        logger.Logger
}

// POST /api/v1/enrollments
func (h *EnrollmentHandler) CreateEnrollment(c *gin.Context) {
    var req dto.CreateEnrollmentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{
            Error:   "invalid_request",
            Message: err.Error(),
        })
        return
    }

    // Validate request
    if err := h.validator.Validate(&req); err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{
            Error:   "validation_failed",
            Message: err.Error(),
        })
        return
    }

    // Process enrollment
    resp, err := h.enrollmentSvc.ProcessEnrollment(c.Request.Context(), &service.EnrollmentRequest{
        UserID:       c.GetString("user_id"),
        InternshipID: req.InternshipID,
        Amount:       req.Amount,
        Currency:     req.Currency,
        PaymentMethod: req.PaymentMethod,
    })
    if err != nil {
        h.logger.Error("enrollment processing failed", "error", err)
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
            Error:   "enrollment_failed",
            Message: "Failed to process enrollment",
        })
        return
    }

    c.JSON(http.StatusCreated, resp)
}
```

---

### **Phase 2: Advanced Payment Features (Weeks 5-8)**

#### **Week 5: Discount & Coupon System**

**Objectives:**

- Implement comprehensive discount system
- Add coupon management
- Integrate discounts with payment flow

**Deliverables:**

- [ ] Discount service implementation
- [ ] Coupon management API
- [ ] Discount calculation engine
- [ ] Admin interface for discount management
- [ ] Bulk discount support

**Discount System:**

```go
type DiscountService struct {
    discountRepo discount.Repository
    calculator   *discount.Calculator
    validator    *discount.Validator
}

type DiscountCalculator struct {
    rules []DiscountRule
}

func (c *DiscountCalculator) CalculateDiscount(
    ctx context.Context,
    originalAmount decimal.Decimal,
    coupons []string,
    userID string,
) (*DiscountResult, error) {
    result := &DiscountResult{
        OriginalAmount: originalAmount,
        FinalAmount:    originalAmount,
        AppliedDiscounts: make([]*AppliedDiscount, 0),
    }

    for _, couponCode := range coupons {
        discount, err := c.discountRepo.GetByCouponCode(ctx, couponCode)
        if err != nil {
            continue // Skip invalid coupons
        }

        if !c.validator.IsValidForUser(discount, userID) {
            continue
        }

        applied := c.applyDiscount(result.FinalAmount, discount)
        if applied.Amount.GreaterThan(decimal.Zero) {
            result.AppliedDiscounts = append(result.AppliedDiscounts, applied)
            result.FinalAmount = result.FinalAmount.Sub(applied.Amount)
            result.TotalDiscount = result.TotalDiscount.Add(applied.Amount)
        }
    }

    return result, nil
}
```

#### **Week 6: Multiple Payment Methods**

**Objectives:**

- Expand payment method support
- Implement UPI, wallets, and bank transfers
- Add gift card system

**Deliverables:**

- [ ] UPI payment integration
- [ ] Digital wallet support (PayPal, Apple Pay, Google Pay)
- [ ] Bank transfer support
- [ ] Gift card system
- [ ] Offline payment tracking

**Payment Method Expansion:**

```go
// UPI Payment Adapter
type UPIPaymentAdapter struct {
    client *upi.Client
    config *UPIConfig
}

func (a *UPIPaymentAdapter) ProcessPayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
    upiReq := &upi.PaymentRequest{
        Amount:      req.Amount,
        Currency:    req.Currency,
        VPA:         req.PaymentDetails["vpa"].(string),
        Description: req.Description,
        OrderID:     req.OrderID,
    }

    resp, err := a.client.CreatePaymentRequest(ctx, upiReq)
    if err != nil {
        return nil, fmt.Errorf("UPI payment failed: %w", err)
    }

    return &PaymentResponse{
        TransactionID: resp.TransactionID,
        Status:       PaymentStatusPending,
        RedirectURL:  resp.PaymentURL,
        PaymentMethod: PaymentMethodUPI,
    }, nil
}

// Gift Card System
type GiftCardService struct {
    cardRepo    giftcard.Repository
    paymentRepo payment.Repository
    generator   *giftcard.CodeGenerator
}

func (s *GiftCardService) ApplyGiftCard(ctx context.Context, req *ApplyGiftCardRequest) (*GiftCardApplication, error) {
    card, err := s.cardRepo.GetByCode(ctx, req.Code)
    if err != nil {
        return nil, fmt.Errorf("invalid gift card: %w", err)
    }

    if card.Status != giftcard.StatusActive {
        return nil, fmt.Errorf("gift card is not active")
    }

    if card.Balance.LessThan(req.Amount) {
        return nil, fmt.Errorf("insufficient gift card balance")
    }

    // Apply gift card
    newBalance := card.Balance.Sub(req.Amount)
    if err := s.cardRepo.UpdateBalance(ctx, card.ID, newBalance); err != nil {
        return nil, fmt.Errorf("failed to update gift card balance: %w", err)
    }

    return &GiftCardApplication{
        CardID:        card.ID,
        AmountApplied: req.Amount,
        RemainingBalance: newBalance,
    }, nil
}
```

#### **Week 7: Subscription Management**

**Objectives:**

- Implement subscription billing
- Add recurring payment support
- Create subscription lifecycle management

**Deliverables:**

- [ ] Subscription service
- [ ] Recurring billing system
- [ ] Subscription upgrade/downgrade
- [ ] Dunning management
- [ ] Subscription analytics

#### **Week 8: Refund & Dispute System**

**Objectives:**

- Implement comprehensive refund system
- Add dispute management
- Create automatic refund policies

**Deliverables:**

- [ ] Refund processing service
- [ ] Dispute management system
- [ ] Automatic refund rules
- [ ] Chargeback handling
- [ ] Reconciliation system

---

### **Phase 3: Enterprise Features (Weeks 9-12)**

#### **Week 9: Invoicing & Billing**

**Objectives:**

- Implement invoice generation
- Add tax calculation
- Create billing automation

**Deliverables:**

- [ ] Invoice generation service
- [ ] Tax calculation engine
- [ ] PDF invoice templates
- [ ] Automated billing cycles
- [ ] Multi-currency invoicing

#### **Week 10: Advanced Analytics & Reporting**

**Objectives:**

- Create comprehensive reporting system
- Add real-time analytics
- Implement business intelligence features

**Deliverables:**

- [ ] Payment analytics dashboard
- [ ] Financial reporting system
- [ ] Real-time monitoring
- [ ] Business intelligence queries
- [ ] Export functionality

#### **Week 11: Notification & Communication**

**Objectives:**

- Implement multi-channel notifications
- Add email template system
- Create SMS and push notification support

**Deliverables:**

- [ ] Notification service
- [ ] Email template engine
- [ ] SMS integration
- [ ] Push notification system
- [ ] Notification preferences

#### **Week 12: Webhook & Integration System**

**Objectives:**

- Create comprehensive webhook system
- Add third-party integrations
- Implement event-driven architecture

**Deliverables:**

- [ ] Webhook management system
- [ ] Event processing pipeline
- [ ] Third-party integrations
- [ ] API gateway enhancements
- [ ] Integration testing

---

### **Phase 4: Optimization & Launch (Weeks 13-16)**

#### **Week 13: Performance Optimization**

**Objectives:**

- Optimize system performance
- Add caching layers
- Implement load balancing

**Deliverables:**

- [ ] Database query optimization
- [ ] Redis caching implementation
- [ ] Connection pooling
- [ ] Load balancer configuration
- [ ] Performance monitoring

#### **Week 14: Security & Compliance**

**Objectives:**

- Ensure PCI DSS compliance
- Implement security best practices
- Add fraud detection

**Deliverables:**

- [ ] Security audit
- [ ] PCI DSS compliance verification
- [ ] Fraud detection system
- [ ] Security monitoring
- [ ] Compliance documentation

#### **Week 15: Testing & Quality Assurance**

**Objectives:**

- Comprehensive testing
- Load testing
- Security testing

**Deliverables:**

- [ ] End-to-end testing
- [ ] Load testing results
- [ ] Security penetration testing
- [ ] User acceptance testing
- [ ] Bug fixes and optimizations

#### **Week 16: Production Deployment**

**Objectives:**

- Deploy to production
- Monitor system health
- Handle post-launch issues

**Deliverables:**

- [ ] Production deployment
- [ ] Monitoring setup
- [ ] Alerting configuration
- [ ] Documentation completion
- [ ] Team training

---

## 🛠️ Technical Implementation Details

### **Design Patterns to Implement**

#### **1. Strategy Pattern - Payment Gateway Selection**

```go
type PaymentStrategy interface {
    ProcessPayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error)
    GetName() string
    GetSupportedMethods() []PaymentMethod
}

type PaymentContext struct {
    strategy PaymentStrategy
}

func (c *PaymentContext) SetStrategy(strategy PaymentStrategy) {
    c.strategy = strategy
}

func (c *PaymentContext) ExecutePayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
    return c.strategy.ProcessPayment(ctx, req)
}
```

#### **2. Factory Pattern - Gateway Creation**

```go
type GatewayFactory interface {
    CreateGateway(gatewayType string, config GatewayConfig) (PaymentGateway, error)
}

type ConcreteGatewayFactory struct {
    configManager *config.Manager
}

func (f *ConcreteGatewayFactory) CreateGateway(gatewayType string, config GatewayConfig) (PaymentGateway, error) {
    switch gatewayType {
    case "stripe":
        return NewStripeGateway(config.StripeConfig)
    case "razorpay":
        return NewRazorpayGateway(config.RazorpayConfig)
    case "paypal":
        return NewPayPalGateway(config.PayPalConfig)
    default:
        return nil, fmt.Errorf("unsupported gateway type: %s", gatewayType)
    }
}
```

#### **3. Observer Pattern - Event Handling**

```go
type PaymentEventObserver interface {
    OnPaymentSuccess(event *PaymentSuccessEvent)
    OnPaymentFailure(event *PaymentFailureEvent)
    OnRefundProcessed(event *RefundEvent)
}

type PaymentEventPublisher struct {
    observers []PaymentEventObserver
    mu        sync.RWMutex
}

func (p *PaymentEventPublisher) Subscribe(observer PaymentEventObserver) {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.observers = append(p.observers, observer)
}

func (p *PaymentEventPublisher) NotifyPaymentSuccess(event *PaymentSuccessEvent) {
    p.mu.RLock()
    defer p.mu.RUnlock()
    for _, observer := range p.observers {
        go observer.OnPaymentSuccess(event)
    }
}
```

#### **4. Decorator Pattern - Payment Enhancement**

```go
type PaymentProcessor interface {
    ProcessPayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error)
}

type PaymentProcessorDecorator struct {
    processor PaymentProcessor
}

type LoggingDecorator struct {
    PaymentProcessorDecorator
    logger logger.Logger
}

func (d *LoggingDecorator) ProcessPayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
    start := time.Now()
    d.logger.Info("Processing payment",
        "amount", req.Amount,
        "method", req.PaymentMethod,
        "user_id", req.UserID)

    resp, err := d.processor.ProcessPayment(ctx, req)

    duration := time.Since(start)
    if err != nil {
        d.logger.Error("Payment processing failed",
            "error", err,
            "duration", duration)
    } else {
        d.logger.Info("Payment processed successfully",
            "transaction_id", resp.TransactionID,
            "duration", duration)
    }

    return resp, err
}

type MetricsDecorator struct {
    PaymentProcessorDecorator
    metrics metrics.Collector
}

func (d *MetricsDecorator) ProcessPayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
    start := time.Now()

    resp, err := d.processor.ProcessPayment(ctx, req)

    duration := time.Since(start)
    d.metrics.RecordPaymentProcessingTime(duration)

    if err != nil {
        d.metrics.IncrementPaymentFailures(req.PaymentMethod)
    } else {
        d.metrics.IncrementPaymentSuccesses(req.PaymentMethod)
    }

    return resp, err
}
```

### **Database Schema Design**

#### **Core Tables**

```sql
-- Payment Gateways Configuration
CREATE TABLE payment_gateways (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,
    type VARCHAR(50) NOT NULL,
    status gateway_status DEFAULT 'active',
    priority INTEGER DEFAULT 0,
    config JSONB NOT NULL,
    health_check_url VARCHAR(255),
    last_health_check TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Payment Methods
CREATE TABLE payment_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gateway_id UUID REFERENCES payment_gateways(id),
    method_type payment_method_type NOT NULL,
    method_name VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    min_amount DECIMAL(10,2),
    max_amount DECIMAL(10,2),
    supported_currencies TEXT[],
    config JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Enrollments
CREATE TABLE enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    internship_id UUID NOT NULL REFERENCES internships(id),
    status enrollment_status NOT NULL DEFAULT 'pending',
    payment_id UUID REFERENCES payments(id),
    discount_amount DECIMAL(10,2) DEFAULT 0,
    final_amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    enrolled_at TIMESTAMP,
    expires_at TIMESTAMP,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, internship_id)
);

-- Payments
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    enrollment_id UUID REFERENCES enrollments(id),
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    status payment_status NOT NULL DEFAULT 'pending',
    gateway_id UUID NOT NULL REFERENCES payment_gateways(id),
    gateway_transaction_id VARCHAR(255),
    payment_method payment_method_type,
    payment_details JSONB,
    failure_reason TEXT,
    retry_count INTEGER DEFAULT 0,
    expires_at TIMESTAMP,
    processed_at TIMESTAMP,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Payment Attempts (for retry logic)
CREATE TABLE payment_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id UUID NOT NULL REFERENCES payments(id),
    gateway_id UUID NOT NULL REFERENCES payment_gateways(id),
    attempt_number INTEGER NOT NULL,
    status payment_status NOT NULL,
    gateway_response JSONB,
    error_message TEXT,
    processing_time_ms INTEGER,
    processed_at TIMESTAMP DEFAULT NOW(),
    INDEX idx_payment_attempts_payment_id (payment_id),
    INDEX idx_payment_attempts_status (status)
);

-- Discounts and Coupons
CREATE TABLE discounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type discount_type NOT NULL,
    value DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3),
    min_order_amount DECIMAL(10,2),
    max_discount_amount DECIMAL(10,2),
    usage_limit INTEGER,
    usage_limit_per_user INTEGER DEFAULT 1,
    used_count INTEGER DEFAULT 0,
    valid_from TIMESTAMP NOT NULL,
    valid_until TIMESTAMP NOT NULL,
    applicable_to JSONB, -- {courses: [], categories: [], users: []}
    stackable BOOLEAN DEFAULT false,
    auto_apply BOOLEAN DEFAULT false,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    INDEX idx_discounts_code (code),
    INDEX idx_discounts_valid_period (valid_from, valid_until),
    INDEX idx_discounts_auto_apply (auto_apply)
);

-- Discount Usage Tracking
CREATE TABLE discount_usages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    discount_id UUID NOT NULL REFERENCES discounts(id),
    user_id UUID NOT NULL REFERENCES users(id),
    payment_id UUID NOT NULL REFERENCES payments(id),
    discount_amount DECIMAL(10,2) NOT NULL,
    used_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(discount_id, payment_id)
);

-- Refunds
CREATE TABLE refunds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id UUID NOT NULL REFERENCES payments(id),
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    reason refund_reason NOT NULL,
    reason_description TEXT,
    status refund_status NOT NULL DEFAULT 'pending',
    gateway_refund_id VARCHAR(255),
    gateway_response JSONB,
    initiated_by UUID REFERENCES users(id),
    processed_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    INDEX idx_refunds_payment_id (payment_id),
    INDEX idx_refunds_status (status)
);

-- Subscriptions
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    plan_id UUID NOT NULL REFERENCES subscription_plans(id),
    status subscription_status NOT NULL DEFAULT 'active',
    billing_cycle billing_cycle NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    current_period_start TIMESTAMP NOT NULL,
    current_period_end TIMESTAMP NOT NULL,
    trial_end TIMESTAMP,
    cancel_at_period_end BOOLEAN DEFAULT FALSE,
    canceled_at TIMESTAMP,
    gateway_id UUID NOT NULL REFERENCES payment_gateways(id),
    gateway_subscription_id VARCHAR(255),
    gateway_customer_id VARCHAR(255),
    next_billing_date TIMESTAMP,
    failure_count INTEGER DEFAULT 0,
    last_payment_attempt TIMESTAMP,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    INDEX idx_subscriptions_user_id (user_id),
    INDEX idx_subscriptions_status (status),
    INDEX idx_subscriptions_next_billing (next_billing_date)
);

-- Subscription Plans
CREATE TABLE subscription_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    billing_cycle billing_cycle NOT NULL,
    trial_days INTEGER DEFAULT 0,
    max_enrollments INTEGER, -- NULL for unlimited
    features JSONB,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Invoices
CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_number VARCHAR(50) UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    subscription_id UUID REFERENCES subscriptions(id),
    payment_id UUID REFERENCES payments(id),
    subtotal DECIMAL(10,2) NOT NULL,
    tax_amount DECIMAL(10,2) DEFAULT 0,
    discount_amount DECIMAL(10,2) DEFAULT 0,
    total_amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    status invoice_status NOT NULL DEFAULT 'draft',
    due_date TIMESTAMP,
    paid_at TIMESTAMP,
    billing_address JSONB,
    line_items JSONB NOT NULL,
    tax_details JSONB,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    INDEX idx_invoices_user_id (user_id),
    INDEX idx_invoices_status (status),
    INDEX idx_invoices_due_date (due_date)
);

-- Gift Cards
CREATE TABLE gift_cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL,
    initial_amount DECIMAL(10,2) NOT NULL,
    current_balance DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    status gift_card_status NOT NULL DEFAULT 'active',
    issued_to_user_id UUID REFERENCES users(id),
    issued_by_user_id UUID REFERENCES users(id),
    expires_at TIMESTAMP,
    last_used_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    INDEX idx_gift_cards_code (code),
    INDEX idx_gift_cards_status (status)
);

-- Gift Card Transactions
CREATE TABLE gift_card_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gift_card_id UUID NOT NULL REFERENCES gift_cards(id),
    payment_id UUID REFERENCES payments(id),
    transaction_type gift_card_transaction_type NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    balance_before DECIMAL(10,2) NOT NULL,
    balance_after DECIMAL(10,2) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Webhooks
CREATE TABLE webhooks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type VARCHAR(100) NOT NULL,
    gateway_id UUID REFERENCES payment_gateways(id),
    payload JSONB NOT NULL,
    signature VARCHAR(255),
    status webhook_status NOT NULL DEFAULT 'pending',
    retry_count INTEGER DEFAULT 0,
    last_retry_at TIMESTAMP,
    processed_at TIMESTAMP,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    INDEX idx_webhooks_status (status),
    INDEX idx_webhooks_event_type (event_type),
    INDEX idx_webhooks_created_at (created_at)
);

-- Webhook Endpoints (for outgoing webhooks)
CREATE TABLE webhook_endpoints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    url VARCHAR(500) NOT NULL,
    secret VARCHAR(255) NOT NULL,
    events TEXT[] NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Audit Logs
CREATE TABLE payment_audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    old_values JSONB,
    new_values JSONB,
    performed_by UUID REFERENCES users(id),
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    INDEX idx_audit_logs_entity (entity_type, entity_id),
    INDEX idx_audit_logs_created_at (created_at)
);
```

#### **Enums**

```sql
-- Enrollment statuses
CREATE TYPE enrollment_status AS ENUM (
    'pending',
    'payment_pending',
    'enrolled',
    'cancelled',
    'expired',
    'failed',
    'refunded'
);

-- Payment statuses
CREATE TYPE payment_status AS ENUM (
    'pending',
    'processing',
    'succeeded',
    'failed',
    'cancelled',
    'expired',
    'refunded',
    'partially_refunded'
);

-- Payment method types
CREATE TYPE payment_method_type AS ENUM (
    'credit_card',
    'debit_card',
    'upi',
    'net_banking',
    'wallet',
    'bank_transfer',
    'gift_card',
    'cryptocurrency',
    'offline'
);

-- Gateway statuses
CREATE TYPE gateway_status AS ENUM (
    'active',
    'inactive',
    'maintenance',
    'deprecated'
);

-- Discount types
CREATE TYPE discount_type AS ENUM (
    'percentage',
    'fixed_amount',
    'free_shipping',
    'bogo',
    'tiered'
);

-- Refund reasons
CREATE TYPE refund_reason AS ENUM (
    'customer_request',
    'duplicate_payment',
    'fraudulent_transaction',
    'service_not_delivered',
    'technical_error',
    'policy_violation',
    'other'
);

-- Refund statuses
CREATE TYPE refund_status AS ENUM (
    'pending',
    'processing',
    'succeeded',
    'failed',
    'cancelled'
);

-- Subscription statuses
CREATE TYPE subscription_status AS ENUM (
    'active',
    'cancelled',
    'expired',
    'past_due',
    'unpaid',
    'trialing'
);

-- Billing cycles
CREATE TYPE billing_cycle AS ENUM (
    'weekly',
    'monthly',
    'quarterly',
    'yearly'
);

-- Invoice statuses
CREATE TYPE invoice_status AS ENUM (
    'draft',
    'sent',
    'paid',
    'overdue',
    'cancelled',
    'refunded'
);

-- Gift card statuses
CREATE TYPE gift_card_status AS ENUM (
    'active',
    'expired',
    'used',
    'cancelled'
);

-- Gift card transaction types
CREATE TYPE gift_card_transaction_type AS ENUM (
    'issued',
    'used',
    'refunded',
    'expired'
);

-- Webhook statuses
CREATE TYPE webhook_status AS ENUM (
    'pending',
    'processing',
    'succeeded',
    'failed',
    'expired'
);
```

---

## 📊 Monitoring & Metrics

### **Key Performance Indicators (KPIs)**

#### **Technical Metrics**

- **Payment Success Rate**: Target > 95%
- **API Response Time**: 95th percentile < 500ms
- **Database Query Performance**: 95th percentile < 100ms
- **System Uptime**: Target 99.9%
- **Error Rate**: < 0.1% for critical paths

#### **Business Metrics**

- **Conversion Rate**: % of users completing enrollment
- **Cart Abandonment Rate**: % abandoning at payment step
- **Average Order Value**: Revenue per successful transaction
- **Refund Rate**: % of payments refunded
- **Customer Lifetime Value**: Long-term revenue per user

#### **Operational Metrics**

- **Gateway Performance**: Success rates by gateway
- **Payment Method Adoption**: Usage distribution
- **Geographic Performance**: Success rates by region
- **Fraud Detection**: False positive/negative rates
- **Support Ticket Volume**: Payment-related issues

### **Monitoring Implementation**

```go
// Metrics Collection
type PaymentMetrics struct {
    successCounter    prometheus.CounterVec
    failureCounter    prometheus.CounterVec
    processingTime    prometheus.HistogramVec
    activePayments    prometheus.GaugeVec
    gatewayHealth     prometheus.GaugeVec
}

func NewPaymentMetrics() *PaymentMetrics {
    return &PaymentMetrics{
        successCounter: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "payments_successful_total",
                Help: "Total number of successful payments",
            },
            []string{"gateway", "method", "currency"},
        ),
        failureCounter: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "payments_failed_total",
                Help: "Total number of failed payments",
            },
            []string{"gateway", "method", "reason"},
        ),
        processingTime: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "payment_processing_duration_seconds",
                Help: "Time taken to process payments",
                Buckets: prometheus.DefBuckets,
            },
            []string{"gateway", "method"},
        ),
    }
}

func (m *PaymentMetrics) RecordSuccess(gateway, method, currency string) {
    m.successCounter.WithLabelValues(gateway, method, currency).Inc()
}

func (m *PaymentMetrics) RecordFailure(gateway, method, reason string) {
    m.failureCounter.WithLabelValues(gateway, method, reason).Inc()
}

func (m *PaymentMetrics) RecordProcessingTime(gateway, method string, duration time.Duration) {
    m.processingTime.WithLabelValues(gateway, method).Observe(duration.Seconds())
}
```

### **Alerting Rules**

```yaml
# Prometheus Alerting Rules
groups:
  - name: payment_system
    rules:
      - alert: HighPaymentFailureRate
        expr: rate(payments_failed_total[5m]) / rate(payments_total[5m]) > 0.05
        for: 2m
        labels:
          severity: critical
        annotations:
          summary: "High payment failure rate detected"
          description: "Payment failure rate is {{ $value }}% over the last 5 minutes"

      - alert: PaymentProcessingLatency
        expr: histogram_quantile(0.95, rate(payment_processing_duration_seconds_bucket[5m])) > 5
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "High payment processing latency"
          description: "95th percentile latency is {{ $value }} seconds"

      - alert: GatewayDown
        expr: gateway_health_check == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Payment gateway is down"
          description: "Gateway {{ $labels.gateway }} is not responding"

      - alert: HighRefundRate
        expr: rate(refunds_total[1h]) / rate(payments_successful_total[1h]) > 0.10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High refund rate detected"
          description: "Refund rate is {{ $value }}% over the last hour"
```

---

## 🔐 Security Implementation

### **Security Measures**

#### **Data Protection**

```go
// Encryption Service
type EncryptionService struct {
    key []byte
}

func (s *EncryptionService) EncryptPII(data string) (string, error) {
    block, err := aes.NewCipher(s.key)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }

    ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *EncryptionService) DecryptPII(encryptedData string) (string, error) {
    data, err := base64.StdEncoding.DecodeString(encryptedData)
    if err != nil {
        return "", err
    }

    block, err := aes.NewCipher(s.key)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonceSize := gcm.NonceSize()
    if len(data) < nonceSize {
        return "", fmt.Errorf("ciphertext too short")
    }

    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}
```

#### **Fraud Detection**

```go
// Fraud Detection Service
type FraudDetectionService struct {
    rules []FraudRule
    ml    *MLModel
}

type FraudAnalysis struct {
    RiskScore    float64
    IsBlocked    bool
    Reasons      []string
    Recommended  FraudAction
}

func (s *FraudDetectionService) AnalyzePayment(ctx context.Context, req *PaymentRequest) (*FraudAnalysis, error) {
    analysis := &FraudAnalysis{
        RiskScore: 0.0,
        Reasons:   make([]string, 0),
    }

    // Rule-based analysis
    for _, rule := range s.rules {
        if rule.Matches(req) {
            analysis.RiskScore += rule.RiskWeight
            analysis.Reasons = append(analysis.Reasons, rule.Description)
        }
    }

    // ML-based analysis
    if s.ml != nil {
        mlScore := s.ml.Predict(req)
        analysis.RiskScore = (analysis.RiskScore + mlScore) / 2
    }

    // Determine action
    switch {
    case analysis.RiskScore > 0.8:
        analysis.IsBlocked = true
        analysis.Recommended = FraudActionBlock
    case analysis.RiskScore > 0.5:
        analysis.Recommended = FraudActionReview
    default:
        analysis.Recommended = FraudActionAllow
    }

    return analysis, nil
}
```

#### **Rate Limiting**

```go
// Rate Limiter
type RateLimiter struct {
    redis  *redis.Client
    limits map[string]RateLimit
}

type RateLimit struct {
    Requests int
    Window   time.Duration
}

func (r *RateLimiter) CheckLimit(ctx context.Context, key string, limitType string) (bool, error) {
    limit, exists := r.limits[limitType]
    if !exists {
        return true, nil // No limit configured
    }

    pipe := r.redis.Pipeline()
    incr := pipe.Incr(ctx, key)
    pipe.Expire(ctx, key, limit.Window)

    _, err := pipe.Exec(ctx)
    if err != nil {
        return false, err
    }

    count := incr.Val()
    return count <= int64(limit.Requests), nil
}
```

---

This comprehensive implementation roadmap provides a structured approach to building a world-class payment system for internship enrollments. The roadmap includes detailed technical specifications, code examples, database schemas, and monitoring strategies to ensure successful delivery.
