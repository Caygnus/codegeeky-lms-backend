# 🏗️ **SERVICE ARCHITECTURE & INTERACTIONS DIAGRAM**

## 🔄 **Service Layer Interactions & Data Flow**

```mermaid
graph TB
  %% API Layer
  subgraph "API Layer"
    API1["Cart API"]
    API2["Order API"]
    API3["Payment API"]
    API4["Enrollment API"]
    API5["Webhook API"]
  end

  %% Service Layer
  subgraph "Service Layer"
    CS["CartService"]
    OS["OrderService"]
    PS["PaymentService"]
    ES["EnrollmentService"]
    DS["DiscountService"]
    NS["NotificationService"]
    AS["AnalyticsService"]
    IS["IdempotencyService"]
    WS["WebhookService"]
  end

  %% Repository Layer
  subgraph "Repository Layer"
    CR["CartRepository"]
    OR["OrderRepository"]
    PR["PaymentRepository"]
    ER["EnrollmentRepository"]
    DR["DiscountRepository"]
    UR["UserRepository"]
    IR["InternshipRepository"]
  end

  %% External Services
  subgraph "External Services"
    RP["Razorpay Gateway"]
    ST["Stripe Gateway"]
    RD["Redis Cache"]
    PG["PostgreSQL DB"]
    MQ["Message Queue"]
    SMTP["SMTP Service"]
  end

  %% API to Service Connections
  API1 --> CS
  API2 --> OS
  API3 --> PS
  API4 --> ES
  API5 --> WS

  %% Service Dependencies
  CS --> CR
  CS --> RD
  CS --> DS
  CS --> AS

  OS --> OR
  OS --> CS
  OS --> PS
  OS --> IS
  OS --> AS

  PS --> PR
  PS --> RP
  PS --> ST
  PS --> MQ

  ES --> ER
  ES --> NS
  ES --> AS
  ES --> IS

  DS --> DR
  WS --> PS
  WS --> OS
  WS --> ES

  %% Repository to Database
  CR --> PG
  OR --> PG
  PR --> PG
  ER --> PG
  DR --> PG
  UR --> PG
  IR --> PG

  %% External Service Connections
  NS --> SMTP
  MQ --> ES
  MQ --> NS
  MQ --> AS

  %% Styling
  classDef api fill:#e3f2fd,stroke:#1976d2,stroke-width:3px
  classDef service fill:#f3e5f5,stroke:#7b1fa2,stroke-width:3px
  classDef repo fill:#e8f5e8,stroke:#388e3c,stroke-width:3px
  classDef external fill:#fff3e0,stroke:#f57c00,stroke-width:3px

  class API1,API2,API3,API4,API5 api
  class CS,OS,PS,ES,DS,NS,AS,IS,WS service
  class CR,OR,PR,ER,DR,UR,IR repo
  class RP,ST,RD,PG,MQ,SMTP external
```

## 🎯 **Scenario-Based Service Interactions**

### **1. Free Internship Flow**

```mermaid
sequenceDiagram
    participant U as User
    participant API as Order API
    participant OS as OrderService
    participant IS as IdempotencyService
    participant CS as CartService
    participant ES as EnrollmentService
    participant NS as NotificationService

    U->>API: POST /api/v1/order
    API->>OS: CreateFromCart(cart_token)
    OS->>IS: CheckKey(idempotency_key)
    IS-->>OS: Key not found
    OS->>CS: ValidateCart(cart_token)
    CS-->>OS: Cart valid, total = 0
    OS->>OS: CreateFreeOrder()
    OS->>ES: CreateEnrollment(order_data)
    ES-->>OS: Enrollment created
    OS->>NS: SendWelcomeEmail(user_id)
    NS-->>OS: Email sent
    OS-->>API: {status: 'completed', enrollment_id}
    API-->>U: Order completed
```

### **2. Paid Internship Flow**

```mermaid
sequenceDiagram
    participant U as User
    participant API as Order API
    participant OS as OrderService
    participant PS as PaymentService
    participant RP as Razorpay
    participant MQ as Message Queue

    U->>API: POST /api/v1/order
    API->>OS: CreateFromCart(cart_token)
    OS->>OS: CreatePaidOrder()
    OS->>PS: CreatePaymentSession(order_data)
    PS->>RP: CreateOrder(amount, currency)
    RP-->>PS: {payment_token, checkout_url}
    PS-->>OS: Payment session created
    OS->>API: {order_token, payment_token, checkout_url}
    API-->>U: Redirect to checkout

    Note over U,RP: User completes payment
    RP->>MQ: payment.captured event
    MQ->>OS: ProcessPaymentSuccess()
    OS->>OS: MarkAsPaid()
    OS->>MQ: order.paid event
```

### **3. Discount Application Flow**

```mermaid
sequenceDiagram
    participant U as User
    participant API as Cart API
    participant CS as CartService
    participant DS as DiscountService
    participant CR as CartRepository

    U->>API: POST /api/v1/cart/{token}/coupon
    API->>CS: ApplyDiscount(cart_token, coupon)
    CS->>DS: ValidateCoupon(coupon)
    DS-->>CS: Coupon valid, 20% off
    CS->>CR: UpdateCart(cart_id, discount_data)
    CR-->>CS: Cart updated
    CS->>CS: RecalculateTotals()
    CS-->>API: {success: true, new_total}
    API-->>U: Discount applied
```

### **4. Webhook Processing Flow**

```mermaid
sequenceDiagram
    participant RP as Razorpay
    participant WS as WebhookService
    participant PS as PaymentService
    participant OS as OrderService
    participant ES as EnrollmentService
    participant MQ as Message Queue

    RP->>WS: POST /webhooks/razorpay
    WS->>WS: VerifySignature(payload)
    WS->>PS: ProcessPaymentWebhook(payload)
    PS->>PS: UpdatePaymentStatus()
    PS->>OS: MarkOrderAsPaid(order_id)
    OS->>MQ: Publish('order.paid', order_data)
    MQ->>ES: CreateEnrollment(order_data)
    ES-->>MQ: Enrollment created
    MQ->>ES: SendConfirmationEmail()
```

## 🔧 **Service Implementation Details**

### **CartService Implementation**

```go
type CartService struct {
    cartRepo    CartRepository
    cache       CacheProvider
    discountSvc DiscountService
    analyticsSvc AnalyticsService
}

func (cs *CartService) CreateCart(req CreateCartRequest) (*Cart, error) {
    // 1. Validate internship exists and is active
    // 2. Generate unique cart token
    // 3. Calculate initial pricing
    // 4. Store in Redis (TTL: 24h)
    // 5. Store in PostgreSQL
    // 6. Track analytics event
}

func (cs *CartService) ApplyDiscount(cartToken, coupon string) error {
    // 1. Validate cart exists and not expired
    // 2. Validate coupon with DiscountService
    // 3. Apply discount calculation
    // 4. Update cart in cache and DB
    // 5. Recalculate totals
}
```

### **OrderService Implementation**

```go
type OrderService struct {
    orderRepo      OrderRepository
    cartSvc        CartService
    paymentSvc     PaymentService
    idempotencySvc IdempotencyService
    analyticsSvc   AnalyticsService
}

func (os *OrderService) CreateFromCart(cartToken, idempotencyKey string) (*Order, error) {
    // 1. Check idempotency key
    // 2. Validate cart and get final pricing
    // 3. Determine if payment required
    // 4. Create order record
    // 5. Handle free vs paid flow
    // 6. Store idempotency response
}

func (os *OrderService) CreateFreeOrder(cartData *Cart) (*Order, error) {
    // 1. Create order with status 'completed'
    // 2. Create enrollment immediately
    // 3. Send welcome notification
    // 4. Track conversion
}
```

### **PaymentService Implementation**

```go
type PaymentService struct {
    paymentRepo PaymentRepository
    razorpay    RazorpayProvider
    stripe      StripeProvider
    messageQueue MessageQueue
}

func (ps *PaymentService) CreatePaymentSession(order *Order) (*PaymentSession, error) {
    // 1. Determine payment gateway based on amount/currency
    // 2. Create payment intent with gateway
    // 3. Store payment session
    // 4. Return checkout URL
}

func (ps *PaymentService) ProcessWebhook(payload []byte, signature string) error {
    // 1. Verify webhook signature
    // 2. Parse payment status
    // 3. Update payment record
    // 4. Trigger order status update
    // 5. Publish event for enrollment
}
```

## 🛡️ **Error Handling & Resilience**

### **Circuit Breaker Implementation**

```go
type CircuitBreaker struct {
    failureThreshold int
    recoveryTimeout  time.Duration
    state           CircuitState
    failureCount    int
    lastFailureTime time.Time
}

func (cb *CircuitBreaker) Execute(command func() error) error {
    if cb.state == Open {
        if time.Since(cb.lastFailureTime) > cb.recoveryTimeout {
            cb.state = HalfOpen
        } else {
            return ErrCircuitBreakerOpen
        }
    }

    err := command()
    if err != nil {
        cb.recordFailure()
        return err
    }

    cb.recordSuccess()
    return nil
}
```

### **Retry Mechanism**

```go
type RetryConfig struct {
    MaxAttempts     int
    InitialDelay    time.Duration
    MaxDelay        time.Duration
    BackoffFactor   float64
    RetryableErrors []int
}

func (rc *RetryConfig) Execute(command func() error) error {
    var lastErr error
    delay := rc.InitialDelay

    for attempt := 1; attempt <= rc.MaxAttempts; attempt++ {
        err := command()
        if err == nil {
            return nil
        }

        lastErr = err
        if !rc.isRetryable(err) {
            return err
        }

        if attempt < rc.MaxAttempts {
            time.Sleep(delay)
            delay = time.Duration(float64(delay) * rc.BackoffFactor)
            if delay > rc.MaxDelay {
                delay = rc.MaxDelay
            }
        }
    }

    return lastErr
}
```

## 📊 **Data Models & Relationships**

### **Cart Entity**

```go
type Cart struct {
    ID              string    `json:"id"`
    Token           string    `json:"token"`
    UserID          *string   `json:"user_id,omitempty"`
    InternshipID    string    `json:"internship_id"`
    Period          Period    `json:"period"`
    Subtotal        int64     `json:"subtotal"`
    Total           int64     `json:"total"`
    TaxAmount       int64     `json:"tax_amount"`
    DiscountAmount  int64     `json:"discount_amount"`
    CouponCode      *string   `json:"coupon_code,omitempty"`
    Status          string    `json:"status"`
    ExpiresAt       time.Time `json:"expires_at"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}
```

### **Order Entity**

```go
type Order struct {
    ID              string    `json:"id"`
    Token           string    `json:"token"`
    CartID          string    `json:"cart_id"`
    UserID          string    `json:"user_id"`
    InternshipID    string    `json:"internship_id"`
    Status          string    `json:"status"`
    Subtotal        int64     `json:"subtotal"`
    Total           int64     `json:"total"`
    TaxAmount       int64     `json:"tax_amount"`
    DiscountAmount  int64     `json:"discount_amount"`
    PaymentToken    *string   `json:"payment_token,omitempty"`
    EnrollmentID    *string   `json:"enrollment_id,omitempty"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}
```

This architecture ensures clean separation of concerns, robust error handling, and scalable service interactions while maintaining data consistency and system reliability.
