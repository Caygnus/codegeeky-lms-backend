# 🛠️ Payment System Implementation Guide

## 📋 Architecture Overview

Based on research of platforms like Udemy, Coursera, and LinkedIn Learning, here's a comprehensive implementation guide for a flexible payment system.

---

## 🏗️ Core System Architecture

### **Design Patterns Implementation**

#### **1. Strategy Pattern - Payment Gateway Management**

```go
// File: internal/payment/gateway/interface.go
package gateway

import (
    "context"
    "github.com/omkar273/codegeeky/internal/types"
)

type PaymentGateway interface {
    Name() string
    ProcessPayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error)
    ProcessRefund(ctx context.Context, req *RefundRequest) (*RefundResponse, error)
    ValidateWebhook(payload []byte, signature string) (*WebhookEvent, error)
    GetSupportedMethods() []types.PaymentMethod
    HealthCheck(ctx context.Context) error
}

type PaymentRequest struct {
    Amount        types.Money           `json:"amount"`
    Currency      string               `json:"currency"`
    UserID        string               `json:"user_id"`
    OrderID       string               `json:"order_id"`
    PaymentMethod types.PaymentMethod  `json:"payment_method"`
    PaymentDetails map[string]interface{} `json:"payment_details"`
    ReturnURL     string               `json:"return_url"`
    CancelURL     string               `json:"cancel_url"`
    Metadata      map[string]string    `json:"metadata"`
}

type PaymentResponse struct {
    TransactionID string                 `json:"transaction_id"`
    Status        types.PaymentStatus    `json:"status"`
    RedirectURL   string                 `json:"redirect_url,omitempty"`
    PaymentMethod types.PaymentMethod    `json:"payment_method"`
    GatewayData   map[string]interface{} `json:"gateway_data,omitempty"`
    ExpiresAt     *time.Time            `json:"expires_at,omitempty"`
}
```

#### **2. Factory Pattern - Gateway Creation**

```go
// File: internal/payment/gateway/factory.go
package gateway

import (
    "fmt"
    "github.com/omkar273/codegeeky/internal/config"
)

type GatewayFactory interface {
    CreateGateway(gatewayType string, cfg *config.GatewayConfig) (PaymentGateway, error)
}

type ConcreteGatewayFactory struct{}

func (f *ConcreteGatewayFactory) CreateGateway(gatewayType string, cfg *config.GatewayConfig) (PaymentGateway, error) {
    switch gatewayType {
    case "stripe":
        return NewStripeGateway(cfg.Stripe)
    case "razorpay":
        return NewRazorpayGateway(cfg.Razorpay)
    case "paypal":
        return NewPayPalGateway(cfg.PayPal)
    default:
        return nil, fmt.Errorf("unsupported gateway type: %s", gatewayType)
    }
}
```

#### **3. Gateway Manager with Load Balancing**

```go
// File: internal/payment/manager.go
package payment

import (
    "context"
    "fmt"
    "sync"
    "github.com/omkar273/codegeeky/internal/payment/gateway"
    "github.com/omkar273/codegeeky/internal/logger"
)

type GatewayManager struct {
    gateways map[string]gateway.PaymentGateway
    selector GatewaySelector
    health   *HealthMonitor
    mu       sync.RWMutex
    logger   *logger.Logger
}

func NewGatewayManager(logger *logger.Logger) *GatewayManager {
    return &GatewayManager{
        gateways: make(map[string]gateway.PaymentGateway),
        selector: NewSmartGatewaySelector(),
        health:   NewHealthMonitor(logger),
        logger:   logger,
    }
}

func (m *GatewayManager) RegisterGateway(name string, gw gateway.PaymentGateway) {
    m.mu.Lock()
    defer m.mu.Unlock()

    m.gateways[name] = gw
    m.health.Monitor(name, gw)
    m.logger.Infow("Gateway registered", "gateway", name)
}

func (m *GatewayManager) ProcessPayment(ctx context.Context, req *gateway.PaymentRequest) (*gateway.PaymentResponse, error) {
    selectedGateway, err := m.selector.SelectGateway(req, m.gateways)
    if err != nil {
        return nil, fmt.Errorf("gateway selection failed: %w", err)
    }

    m.logger.Infow("Processing payment",
        "gateway", selectedGateway.Name(),
        "amount", req.Amount,
        "method", req.PaymentMethod)

    return selectedGateway.ProcessPayment(ctx, req)
}
```

---

## 🔧 Implementation Components

### **1. Enrollment Service with Payment Integration**

```go
// File: internal/service/enrollment.go
package service

import (
    "context"
    "fmt"
    "github.com/omkar273/codegeeky/internal/api/dto"
    "github.com/omkar273/codegeeky/internal/domain/enrollment"
    "github.com/omkar273/codegeeky/internal/payment"
    "github.com/omkar273/codegeeky/internal/types"
)

type EnrollmentService struct {
    paymentManager *payment.GatewayManager
    enrollRepo     enrollment.Repository
    discountSvc    *DiscountService
    notificationSvc *NotificationService
    eventPublisher  *EventPublisher
    logger         *logger.Logger
}

func (s *EnrollmentService) ProcessEnrollment(ctx context.Context, req *dto.EnrollmentRequest) (*dto.EnrollmentResponse, error) {
    // 1. Validate request
    if err := s.validateEnrollmentRequest(ctx, req); err != nil {
        return nil, err
    }

    // 2. Apply discounts
    finalAmount, discounts, err := s.discountSvc.CalculateDiscount(ctx, &discount.CalculateRequest{
        OriginalAmount: req.Amount,
        CouponCodes:    req.CouponCodes,
        UserID:         req.UserID,
        InternshipID:   req.InternshipID,
    })
    if err != nil {
        return nil, fmt.Errorf("discount calculation failed: %w", err)
    }

    // 3. Create enrollment record
    enrollment, err := s.enrollRepo.Create(ctx, &enrollment.CreateRequest{
        UserID:         req.UserID,
        InternshipID:   req.InternshipID,
        OriginalAmount: req.Amount,
        DiscountAmount: req.Amount.Sub(finalAmount),
        FinalAmount:    finalAmount,
        Status:         types.EnrollmentStatusPending,
        AppliedDiscounts: discounts,
    })
    if err != nil {
        return nil, fmt.Errorf("enrollment creation failed: %w", err)
    }

    // 4. Process payment
    paymentReq := &gateway.PaymentRequest{
        Amount:         finalAmount,
        Currency:       req.Currency,
        UserID:         req.UserID,
        OrderID:        enrollment.ID,
        PaymentMethod:  req.PaymentMethod,
        PaymentDetails: req.PaymentDetails,
        ReturnURL:      req.ReturnURL,
        CancelURL:      req.CancelURL,
        Metadata: map[string]string{
            "enrollment_id": enrollment.ID,
            "internship_id": req.InternshipID,
            "type":         "enrollment",
        },
    }

    paymentResp, err := s.paymentManager.ProcessPayment(ctx, paymentReq)
    if err != nil {
        // Update enrollment to failed
        s.enrollRepo.UpdateStatus(ctx, enrollment.ID, types.EnrollmentStatusFailed)
        return nil, fmt.Errorf("payment processing failed: %w", err)
    }

    // 5. Update enrollment with payment info
    enrollment.PaymentID = &paymentResp.TransactionID
    enrollment.Status = types.EnrollmentStatusPaymentPending
    if err := s.enrollRepo.Update(ctx, enrollment); err != nil {
        return nil, fmt.Errorf("enrollment update failed: %w", err)
    }

    // 6. Send notifications
    go s.sendEnrollmentNotifications(ctx, enrollment, paymentResp)

    // 7. Publish event
    s.eventPublisher.PublishEnrollmentCreated(ctx, enrollment)

    return &dto.EnrollmentResponse{
        EnrollmentID:  enrollment.ID,
        PaymentID:     paymentResp.TransactionID,
        Status:        string(enrollment.Status),
        RedirectURL:   paymentResp.RedirectURL,
        FinalAmount:   finalAmount,
        DiscountAmount: req.Amount.Sub(finalAmount),
    }, nil
}
```

### **2. Flexible Discount System**

```go
// File: internal/service/discount.go
package service

import (
    "context"
    "github.com/omkar273/codegeeky/internal/types"
)

type DiscountService struct {
    discountRepo discount.Repository
    calculator   *DiscountCalculator
    validator    *DiscountValidator
}

type DiscountCalculator struct {
    rules []DiscountRule
}

func (c *DiscountCalculator) CalculateDiscount(ctx context.Context, req *CalculateRequest) (types.Money, []*AppliedDiscount, error) {
    result := &DiscountResult{
        OriginalAmount:   req.OriginalAmount,
        FinalAmount:     req.OriginalAmount,
        AppliedDiscounts: make([]*AppliedDiscount, 0),
    }

    // Get applicable discounts
    discounts, err := c.discountRepo.GetApplicableDiscounts(ctx, &discount.Query{
        CouponCodes:   req.CouponCodes,
        UserID:        req.UserID,
        InternshipID:  req.InternshipID,
        Amount:        req.OriginalAmount,
    })
    if err != nil {
        return req.OriginalAmount, nil, err
    }

    // Apply discounts in order of priority
    for _, discount := range discounts {
        if !c.validator.IsValidForApplication(discount, req) {
            continue
        }

        applied := c.applyDiscount(result.FinalAmount, discount)
        if applied.Amount.GreaterThan(types.ZeroMoney) {
            result.AppliedDiscounts = append(result.AppliedDiscounts, applied)
            result.FinalAmount = result.FinalAmount.Sub(applied.Amount)

            // Stop if not stackable
            if !discount.Stackable {
                break
            }
        }
    }

    return result.FinalAmount, result.AppliedDiscounts, nil
}

func (c *DiscountCalculator) applyDiscount(amount types.Money, discount *types.Discount) *AppliedDiscount {
    var discountAmount types.Money

    switch discount.Type {
    case types.DiscountTypePercentage:
        discountAmount = amount.Mul(discount.Value).Div(types.NewMoney(100))
    case types.DiscountTypeFixedAmount:
        discountAmount = discount.ValueMoney
    case types.DiscountTypeBOGO:
        // Buy one get one logic
        discountAmount = amount.Div(types.NewMoney(2))
    }

    // Apply maximum discount limit
    if discount.MaxDiscountAmount != nil && discountAmount.GreaterThan(*discount.MaxDiscountAmount) {
        discountAmount = *discount.MaxDiscountAmount
    }

    // Ensure we don't discount more than the amount
    if discountAmount.GreaterThan(amount) {
        discountAmount = amount
    }

    return &AppliedDiscount{
        DiscountID: discount.ID,
        Code:       discount.Code,
        Type:       discount.Type,
        Amount:     discountAmount,
    }
}
```

### **3. Multi-Gateway Payment Adapters**

#### **Stripe Adapter**

```go
// File: internal/payment/adapters/stripe.go
package adapters

import (
    "context"
    "fmt"
    "github.com/stripe/stripe-go/v72"
    "github.com/stripe/stripe-go/v72/paymentintent"
)

type StripeAdapter struct {
    client *stripe.Client
    config *StripeConfig
}

func (s *StripeAdapter) ProcessPayment(ctx context.Context, req *gateway.PaymentRequest) (*gateway.PaymentResponse, error) {
    params := &stripe.PaymentIntentParams{
        Amount:   stripe.Int64(req.Amount.ToCents()),
        Currency: stripe.String(strings.ToLower(req.Currency)),
        AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
            Enabled: stripe.Bool(true),
        },
        Metadata: req.Metadata,
    }

    // Add customer if exists
    if req.UserID != "" {
        customer, err := s.getOrCreateCustomer(ctx, req.UserID)
        if err == nil {
            params.Customer = stripe.String(customer.ID)
        }
    }

    // Handle different payment methods
    switch req.PaymentMethod {
    case types.PaymentMethodCard:
        // Card payments handled by automatic payment methods
    case types.PaymentMethodUPI:
        params.PaymentMethodTypes = stripe.StringSlice([]string{"upi"})
    case types.PaymentMethodWallet:
        params.PaymentMethodTypes = stripe.StringSlice([]string{"link", "paypal"})
    }

    pi, err := paymentintent.New(params)
    if err != nil {
        return nil, fmt.Errorf("stripe payment intent creation failed: %w", err)
    }

    return &gateway.PaymentResponse{
        TransactionID: pi.ID,
        Status:        s.mapStripeStatus(pi.Status),
        PaymentMethod: req.PaymentMethod,
        GatewayData: map[string]interface{}{
            "client_secret": pi.ClientSecret,
            "payment_intent_id": pi.ID,
        },
    }, nil
}
```

#### **Razorpay Adapter**

```go
// File: internal/payment/adapters/razorpay.go
package adapters

import (
    "context"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
)

type RazorpayAdapter struct {
    client *razorpay.Client
    config *RazorpayConfig
}

func (r *RazorpayAdapter) ProcessPayment(ctx context.Context, req *gateway.PaymentRequest) (*gateway.PaymentResponse, error) {
    orderReq := map[string]interface{}{
        "amount":   req.Amount.ToCents(), // Razorpay uses paise
        "currency": req.Currency,
        "receipt":  req.OrderID,
        "notes":    req.Metadata,
    }

    order, err := r.client.Order.Create(orderReq, nil)
    if err != nil {
        return nil, fmt.Errorf("razorpay order creation failed: %w", err)
    }

    checkoutOptions := map[string]interface{}{
        "key":         r.config.KeyID,
        "amount":      req.Amount.ToCents(),
        "currency":    req.Currency,
        "order_id":    order["id"],
        "callback_url": req.ReturnURL,
        "cancel_url":   req.CancelURL,
        "prefill": map[string]interface{}{
            "contact": req.PaymentDetails["phone"],
            "email":   req.PaymentDetails["email"],
        },
        "theme": map[string]interface{}{
            "color": "#3399cc",
        },
    }

    // Handle UPI payments
    if req.PaymentMethod == types.PaymentMethodUPI {
        checkoutOptions["method"] = map[string]interface{}{
            "upi": true,
        }
    }

    checkoutURL := r.generateCheckoutURL(checkoutOptions)

    return &gateway.PaymentResponse{
        TransactionID: order["id"].(string),
        Status:        types.PaymentStatusPending,
        RedirectURL:   checkoutURL,
        PaymentMethod: req.PaymentMethod,
        GatewayData: map[string]interface{}{
            "order_id": order["id"],
            "key_id":   r.config.KeyID,
        },
    }, nil
}

func (r *RazorpayAdapter) ValidateWebhook(payload []byte, signature string) (*gateway.WebhookEvent, error) {
    expectedSignature := r.generateWebhookSignature(payload)
    if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
        return nil, fmt.Errorf("invalid webhook signature")
    }

    var event map[string]interface{}
    if err := json.Unmarshal(payload, &event); err != nil {
        return nil, fmt.Errorf("invalid webhook payload: %w", err)
    }

    return &gateway.WebhookEvent{
        EventType: event["event"].(string),
        Data:      event,
    }, nil
}

func (r *RazorpayAdapter) generateWebhookSignature(payload []byte) string {
    h := hmac.New(sha256.New, []byte(r.config.WebhookSecret))
    h.Write(payload)
    return hex.EncodeToString(h.Sum(nil))
}
```

### **4. Comprehensive Webhook System**

```go
// File: internal/webhook/processor.go
package webhook

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/omkar273/codegeeky/internal/payment/gateway"
)

type WebhookProcessor struct {
    gatewayManager *payment.GatewayManager
    enrollmentSvc  *service.EnrollmentService
    paymentRepo    payment.Repository
    logger         *logger.Logger
}

func (p *WebhookProcessor) ProcessWebhook(ctx context.Context, req *WebhookRequest) error {
    // 1. Identify gateway
    gateway, err := p.gatewayManager.GetGateway(req.GatewayName)
    if err != nil {
        return fmt.Errorf("unknown gateway: %w", err)
    }

    // 2. Validate webhook
    event, err := gateway.ValidateWebhook(req.Payload, req.Signature)
    if err != nil {
        return fmt.Errorf("webhook validation failed: %w", err)
    }

    // 3. Process based on event type
    switch event.EventType {
    case "payment.captured", "payment_intent.succeeded":
        return p.processPaymentSuccess(ctx, event)
    case "payment.failed", "payment_intent.payment_failed":
        return p.processPaymentFailure(ctx, event)
    case "refund.processed":
        return p.processRefund(ctx, event)
    default:
        p.logger.Warnw("Unhandled webhook event", "event_type", event.EventType)
        return nil
    }
}

func (p *WebhookProcessor) processPaymentSuccess(ctx context.Context, event *gateway.WebhookEvent) error {
    transactionID := p.extractTransactionID(event)

    // Update payment status
    payment, err := p.paymentRepo.GetByTransactionID(ctx, transactionID)
    if err != nil {
        return fmt.Errorf("payment not found: %w", err)
    }

    payment.Status = types.PaymentStatusSucceeded
    payment.ProcessedAt = time.Now()

    if err := p.paymentRepo.Update(ctx, payment); err != nil {
        return fmt.Errorf("payment update failed: %w", err)
    }

    // Update enrollment status
    if payment.EnrollmentID != nil {
        err := p.enrollmentSvc.CompleteEnrollment(ctx, *payment.EnrollmentID)
        if err != nil {
            p.logger.Errorw("Failed to complete enrollment",
                "enrollment_id", *payment.EnrollmentID,
                "error", err)
        }
    }

    return nil
}
```

### **5. Advanced Reporting System**

```go
// File: internal/service/analytics.go
package service

import (
    "context"
    "time"
)

type AnalyticsService struct {
    paymentRepo payment.Repository
    cache       cache.Cache
    logger      *logger.Logger
}

func (s *AnalyticsService) GetPaymentAnalytics(ctx context.Context, req *AnalyticsRequest) (*PaymentAnalytics, error) {
    cacheKey := fmt.Sprintf("analytics:%s:%s:%s", req.StartDate, req.EndDate, req.GroupBy)

    // Try cache first
    if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
        var analytics PaymentAnalytics
        if json.Unmarshal(cached, &analytics) == nil {
            return &analytics, nil
        }
    }

    // Query database
    analytics := &PaymentAnalytics{
        Period: Period{
            Start: req.StartDate,
            End:   req.EndDate,
        },
        Metrics: make(map[string]*Metric),
    }

    // Revenue metrics
    revenue, err := s.paymentRepo.GetRevenue(ctx, req.StartDate, req.EndDate)
    if err != nil {
        return nil, err
    }
    analytics.Metrics["total_revenue"] = &Metric{
        Value: revenue.Total.ToFloat64(),
        Currency: revenue.Currency,
    }

    // Success rate
    successRate, err := s.paymentRepo.GetSuccessRate(ctx, req.StartDate, req.EndDate)
    if err != nil {
        return nil, err
    }
    analytics.Metrics["success_rate"] = &Metric{
        Value: successRate,
        Unit:  "percentage",
    }

    // Gateway performance
    gatewayStats, err := s.paymentRepo.GetGatewayStats(ctx, req.StartDate, req.EndDate)
    if err != nil {
        return nil, err
    }
    analytics.GatewayPerformance = gatewayStats

    // Cache results
    if data, err := json.Marshal(analytics); err == nil {
        s.cache.Set(ctx, cacheKey, data, 15*time.Minute)
    }

    return analytics, nil
}
```

---

## 🚀 Quick Start Implementation

### **1. Set up the basic project structure:**

```bash
# Create directories
mkdir -p internal/payment/{adapters,gateway,manager}
mkdir -p internal/service
mkdir -p internal/webhook
mkdir -p internal/api/v1

# Initialize Go modules
go mod init github.com/omkar273/codegeeky
go mod tidy
```

### **2. Update your service factory:**

```go
// File: internal/service/factory.go
package service

func NewServiceFactory(
    entClient *ent.Client,
    logger *logger.Logger,
    cfg *config.Configuration,
) *ServiceFactory {
    // Initialize payment gateway manager
    gatewayManager := payment.NewGatewayManager(logger)

    // Register gateways
    if cfg.Payment.Stripe.Enabled {
        stripeGateway := adapters.NewStripeAdapter(cfg.Payment.Stripe)
        gatewayManager.RegisterGateway("stripe", stripeGateway)
    }

    if cfg.Payment.Razorpay.Enabled {
        razorpayGateway := adapters.NewRazorpayAdapter(cfg.Payment.Razorpay)
        gatewayManager.RegisterGateway("razorpay", razorpayGateway)
    }

    return &ServiceFactory{
        EnrollmentService: NewEnrollmentService(
            gatewayManager,
            repository.NewEnrollmentRepository(entClient),
            NewDiscountService(repository.NewDiscountRepository(entClient)),
            NewNotificationService(cfg.Notification),
            logger,
        ),
        PaymentService: NewPaymentService(gatewayManager, logger),
        AnalyticsService: NewAnalyticsService(
            repository.NewPaymentRepository(entClient),
            cache.NewRedisCache(cfg.Redis),
            logger,
        ),
    }
}
```

### **3. Add the API routes:**

```go
// File: internal/api/router.go
func NewRouter(services *service.ServiceFactory, cfg *config.Configuration, logger *logger.Logger) *gin.Engine {
    router := gin.New()

    // Middleware
    router.Use(middleware.CORS())
    router.Use(middleware.RequestLogger(logger))
    router.Use(middleware.RateLimiter(cfg.RateLimit))

    // API v1 routes
    v1 := router.Group("/api/v1")
    {
        // Enrollment routes
        enrollments := v1.Group("/enrollments")
        enrollments.Use(middleware.AuthRequired())
        {
            enrollments.POST("", handlers.CreateEnrollment)
            enrollments.GET("/user/:user_id", handlers.GetUserEnrollments)
        }

        // Payment routes
        payments := v1.Group("/payments")
        {
            payments.POST("/webhooks/:gateway", handlers.ProcessWebhook)
            payments.GET("/analytics", middleware.AdminRequired(), handlers.GetPaymentAnalytics)
        }

        // Discount routes
        discounts := v1.Group("/discounts")
        discounts.Use(middleware.AuthRequired())
        {
            discounts.POST("/validate", handlers.ValidateDiscount)
            discounts.POST("/apply", handlers.ApplyDiscount)
        }
    }

    return router
}
```

This implementation provides a solid foundation for your internship enrollment workflow with comprehensive payment integration. The system is designed to be flexible, scalable, and maintainable while supporting multiple payment gateways and methods.
