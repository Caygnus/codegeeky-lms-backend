# 🛒 **COMPREHENSIVE CHECKOUT & ENROLLMENT FLOW PRD**

## 📋 **Executive Summary**

This Product Requirements Document (PRD) outlines the implementation of a production-ready e-commerce checkout and enrollment system for the Interns Go-Backend platform. The system will enable users to seamlessly purchase internships through a modern, fault-tolerant checkout flow with payment integration, discount management, and automated enrollment processing.

### 🎯 **Vision Statement**

Create a world-class checkout experience that rivals platforms like Udemy, Coursera, and Internshala, providing users with a seamless journey from internship discovery to successful enrollment, while maintaining enterprise-grade reliability, security, and scalability.

---

## 🔍 **Market Research & Competitive Analysis**

### **Platform Benchmarking**

#### **Udemy's Checkout Flow**

- **Cart Management**: Persistent cart with session-based items
- **Payment Methods**: Credit cards, PayPal, Apple Pay, regional methods
- **Discount Strategy**: Coupon codes, flash sales, instructor discounts
- **User Experience**: Single-page checkout with real-time validation
- **Success Rate**: 95%+ payment success rate

#### **Coursera's Enrollment Process**

- **Subscription Model**: Monthly/yearly plans with course access
- **Enterprise Integration**: Corporate billing and team management
- **Financial Aid**: Need-based assistance program
- **Global Reach**: Multi-currency support with local payment methods

#### **Internshala's Approach**

- **Internship-Specific**: Tailored for internship and training programs
- **Payment Flexibility**: Multiple payment options including EMI
- **Corporate Partnerships**: Bulk enrollment for companies
- **Regional Focus**: India-specific payment methods (UPI, Net Banking)

### **Key Insights for Our Implementation**

1. **Cart Persistence**: Users expect their cart to persist across sessions
2. **Payment Method Diversity**: Support for local and international payment methods
3. **Discount Transparency**: Clear display of applied discounts and savings
4. **Mobile-First Design**: Optimized for mobile checkout experience
5. **Real-Time Validation**: Immediate feedback on payment and enrollment status

---

## 🎯 **Product Objectives**

### **Primary Goals**

1. **Seamless User Experience**: Reduce checkout abandonment to < 5%
2. **Payment Success Rate**: Achieve > 98% payment success rate
3. **Idempotent Operations**: Zero duplicate enrollments or charges
4. **Fault Tolerance**: 99.9% system uptime with graceful degradation
5. **Scalability**: Support 10,000+ concurrent checkout sessions

### **Success Metrics**

- **Conversion Rate**: > 85% cart-to-enrollment conversion
- **Payment Processing Time**: < 3 seconds average
- **Webhook Reliability**: 100% event processing with dead-letter handling
- **User Satisfaction**: > 4.5/5 rating on checkout experience
- **Technical Performance**: P95 API response time < 200ms

---

## 👥 **User Personas**

### **1. Individual Student (Primary)**

- **Demographics**: 18-25 years, tech-savvy, price-conscious
- **Payment Preferences**: UPI, credit cards, digital wallets
- **Pain Points**: Complex checkout, payment failures, unclear pricing
- **Goals**: Quick, secure, transparent enrollment process

### **2. Working Professional (Secondary)**

- **Demographics**: 25-35 years, career-focused, time-constrained
- **Payment Preferences**: Credit cards, EMI options, corporate billing
- **Pain Points**: Limited payment options, slow processing
- **Goals**: Fast checkout with flexible payment terms

### **3. Corporate Client (Tertiary)**

- **Demographics**: HR managers, L&D professionals
- **Payment Preferences**: Bank transfers, invoicing, bulk payments
- **Pain Points**: Manual processes, lack of reporting
- **Goals**: Bulk enrollments with detailed tracking

---

## 🔧 **Functional Requirements**

### **1. Cart Management System**

#### **Core Features**

- **Cart Creation**: Automatic cart creation for authenticated users
- **Item Management**: Add, remove, update internship quantities
- **Cart Persistence**: Session-based cart with database backup
- **Cart Expiration**: Automatic cleanup after 24 hours of inactivity
- **Cart Sharing**: Share cart via unique URL (optional)

#### **Technical Specifications**

```go
type Cart struct {
    ID            string          `json:"id"`
    UserID        string          `json:"user_id"`
    Status        CartStatus      `json:"status"`
    Items         []CartItem      `json:"items"`
    Subtotal      decimal.Decimal `json:"subtotal"`
    Total         decimal.Decimal `json:"total"`
    DiscountCode  *string         `json:"discount_code,omitempty"`
    DiscountAmount decimal.Decimal `json:"discount_amount"`
    ExpiresAt     time.Time       `json:"expires_at"`
    CreatedAt     time.Time       `json:"created_at"`
    UpdatedAt     time.Time       `json:"updated_at"`
}

type CartItem struct {
    InternshipID string          `json:"internship_id"`
    Title        string          `json:"title"`
    Price        decimal.Decimal `json:"price"`
    Quantity     int             `json:"quantity"`
    Metadata     map[string]any  `json:"metadata,omitempty"`
}
```

#### **API Endpoints**

```
POST   /api/v1/cart                    # Create/retrieve cart
PATCH  /api/v1/cart/items/:item_id     # Update cart item
DELETE /api/v1/cart/items/:item_id     # Remove cart item
POST   /api/v1/cart/discount           # Apply discount code
DELETE /api/v1/cart/discount           # Remove discount
GET    /api/v1/cart                    # Get cart details
```

### **2. Discount & Coupon System**

#### **Coupon Types**

- **Percentage Discount**: 10%, 20%, 50% off total
- **Fixed Amount**: ₹500, $50 off total
- **Free Enrollment**: 100% discount for specific internships
- **Bulk Discount**: Tiered pricing for multiple enrollments
- **Referral Codes**: Friend referral discounts

#### **Validation Rules**

- **Usage Limits**: Per user, total usage, time-based restrictions
- **Minimum Order Value**: Minimum cart value for coupon application
- **Combinability**: Multiple coupons per transaction (configurable)
- **Expiration**: Automatic expiration based on valid date range
- **Category Restrictions**: Internship-specific or category-specific coupons

#### **API Endpoints**

```
POST   /api/v1/cart/discount/validate  # Validate discount code
POST   /api/v1/cart/discount/apply     # Apply discount to cart
DELETE /api/v1/cart/discount           # Remove applied discount
```

### **3. Order Management System**

#### **Order Lifecycle**

1. **PENDING**: Order created, awaiting payment
2. **PAID**: Payment successful, enrollment processing
3. **FAILED**: Payment failed, order cancelled
4. **EXPIRED**: Order expired due to timeout
5. **REFUNDED**: Order refunded (future enhancement)

#### **Order Features**

- **Idempotency**: Prevent duplicate orders with idempotency keys
- **Payment Linking**: Direct link to payment attempts
- **Metadata Storage**: Flexible metadata for business logic
- **Audit Trail**: Complete order history and status changes

#### **Technical Specifications**

```go
type Order struct {
    ID                string          `json:"id"`
    UserID            string          `json:"user_id"`
    CartID            string          `json:"cart_id"`
    Status            OrderStatus     `json:"status"`
    Amount            decimal.Decimal `json:"amount"`
    Currency          string          `json:"currency"`
    PaymentAttemptID  *string         `json:"payment_attempt_id,omitempty"`
    IdempotencyKey    string          `json:"idempotency_key"`
    Metadata          map[string]any  `json:"metadata,omitempty"`
    CreatedAt         time.Time       `json:"created_at"`
    UpdatedAt         time.Time       `json:"updated_at"`
}
```

### **4. Payment Integration System**

#### **Supported Gateways**

- **Primary**: Razorpay (India-focused)
- **Secondary**: Stripe (International)
- **Future**: PayPal, Paytm, PhonePe

#### **Payment Methods**

- **Credit/Debit Cards**: Visa, Mastercard, American Express, RuPay
- **UPI**: All UPI-enabled apps (GPay, PhonePe, Paytm)
- **Net Banking**: Major Indian banks
- **Digital Wallets**: Paytm, PhonePe, Amazon Pay
- **EMI**: No-cost EMI options (Razorpay)

#### **Payment Flow**

1. **Session Creation**: Create payment session with gateway
2. **Redirect**: Redirect user to hosted checkout
3. **Payment Processing**: User completes payment on gateway
4. **Webhook Notification**: Gateway notifies success/failure
5. **Status Update**: Update order and trigger enrollment

### **5. Webhook & Event System**

#### **Webhook Processing**

- **Signature Verification**: HMAC verification for security
- **Idempotency**: Prevent duplicate webhook processing
- **Retry Logic**: Exponential backoff for failed webhooks
- **Dead Letter Queue**: Handle permanently failed events

#### **Event Publishing**

- **Order Events**: OrderPaid, OrderFailed, OrderExpired
- **Payment Events**: PaymentSuccess, PaymentFailed
- **Enrollment Events**: EnrollmentCreated, EnrollmentFailed

#### **Event Consumers**

- **Enrollment Worker**: Process successful payments and create enrollments
- **Notification Worker**: Send confirmation emails and push notifications
- **Analytics Worker**: Track conversion metrics and business intelligence

### **6. Enrollment Automation**

#### **Enrollment Creation**

- **Automatic Processing**: Create enrollments after successful payment
- **Idempotency**: Prevent duplicate enrollments
- **Batch Processing**: Handle multiple internships in single order
- **Status Tracking**: Track enrollment status and progress

#### **Enrollment Features**

- **Access Control**: Grant immediate access to internship content
- **Progress Tracking**: Initialize learning progress
- **Certificate Eligibility**: Mark user as eligible for certification
- **Support Integration**: Create support tickets if needed

---

## 🏗️ **Technical Architecture**

### **System Components**

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend UI   │    │   API Gateway   │    │   Load Balancer │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Gin Router    │
                    └─────────────────┘
                                 │
         ┌───────────────────────┼───────────────────────┐
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Cart Service   │    │  Order Service  │    │ Payment Service │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   PostgreSQL    │
                    └─────────────────┘
                                 │
         ┌───────────────────────┼───────────────────────┐
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Webhook Handler│    │  Event Consumer │    │ Payment Gateway │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### **Database Schema**

#### **Cart Table**

```sql
CREATE TABLE carts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    status cart_status DEFAULT 'draft',
    items JSONB NOT NULL DEFAULT '[]',
    subtotal DECIMAL(10,2) NOT NULL DEFAULT 0,
    total DECIMAL(10,2) NOT NULL DEFAULT 0,
    discount_code VARCHAR(50),
    discount_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    UNIQUE(user_id, status)
);

CREATE INDEX idx_carts_user_id ON carts(user_id);
CREATE INDEX idx_carts_expires_at ON carts(expires_at);
```

#### **Order Table**

```sql
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    cart_id UUID NOT NULL REFERENCES carts(id),
    status order_status DEFAULT 'pending',
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'INR',
    payment_attempt_id UUID REFERENCES payment_attempts(id),
    idempotency_key VARCHAR(255) UNIQUE NOT NULL,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_payment_attempt_id ON orders(payment_attempt_id);
```

#### **Idempotency Table**

```sql
CREATE TABLE idempotency_keys (
    key VARCHAR(255) PRIMARY KEY,
    request_hash VARCHAR(64) NOT NULL,
    response_body JSONB,
    status_code INTEGER,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_idempotency_expires_at ON idempotency_keys(expires_at);
```

### **Service Dependencies**

```
CartService
├── CartRepository
├── DiscountService
└── PricingService

OrderService
├── OrderRepository
├── CartService
├── PaymentService
└── IdempotencyService

CheckoutService
├── OrderService
├── PaymentService
├── CartService
└── EventPublisher

PaymentService
├── PaymentRepository
├── PaymentGateway
└── WebhookHandler

EnrollmentWorker
├── InternshipEnrollmentRepository
├── OrderRepository
└── NotificationService
```

---

## 🔒 **Security Requirements**

### **Data Protection**

- **PCI DSS Compliance**: Secure handling of payment data
- **Data Encryption**: AES-256 encryption for sensitive data
- **Tokenization**: Payment method tokenization for security
- **Audit Logging**: Complete audit trail for all transactions

### **Authentication & Authorization**

- **JWT Authentication**: Secure API access with JWT tokens
- **ABAC Authorization**: Attribute-based access control
- **Rate Limiting**: Prevent abuse with rate limiting
- **Input Validation**: Comprehensive input sanitization

### **Payment Security**

- **Webhook Verification**: HMAC signature verification
- **Idempotency**: Prevent duplicate charges
- **Fraud Detection**: ML-based fraud prevention (future)
- **Secure Redirects**: HTTPS-only payment redirects

---

## 📊 **Performance Requirements**

### **Response Time Targets**

- **Cart Operations**: < 100ms
- **Order Creation**: < 200ms
- **Payment Session**: < 500ms
- **Webhook Processing**: < 1s
- **Enrollment Creation**: < 2s

### **Scalability Targets**

- **Concurrent Users**: 10,000+ simultaneous checkout sessions
- **Throughput**: 1,000+ orders per minute
- **Database**: 100,000+ orders per day
- **Storage**: 1TB+ data storage capacity

### **Availability Targets**

- **System Uptime**: 99.9% availability
- **Payment Gateway**: 99.95% uptime
- **Database**: 99.99% uptime
- **Recovery Time**: < 5 minutes for critical failures

---

## 🧪 **Testing Strategy**

### **Unit Testing**

- **Service Layer**: 90%+ code coverage
- **Repository Layer**: Database operation testing
- **Utility Functions**: Helper function testing
- **Mock Integration**: External service mocking

### **Integration Testing**

- **API Endpoints**: End-to-end API testing
- **Database Operations**: Transaction testing
- **Payment Gateway**: Sandbox environment testing
- **Webhook Processing**: Event flow testing

### **Load Testing**

- **Stress Testing**: Maximum capacity testing
- **Performance Testing**: Response time validation
- **Concurrency Testing**: Race condition testing
- **Failover Testing**: System resilience testing

### **Security Testing**

- **Penetration Testing**: Vulnerability assessment
- **OWASP Testing**: OWASP Top 10 compliance
- **Payment Security**: PCI DSS compliance testing
- **Data Protection**: GDPR compliance validation

---

## 📈 **Monitoring & Observability**

### **Application Monitoring**

- **Health Checks**: Service health monitoring
- **Performance Metrics**: Response time tracking
- **Error Tracking**: Error rate monitoring
- **Business Metrics**: Conversion rate tracking

### **Infrastructure Monitoring**

- **Server Metrics**: CPU, memory, disk usage
- **Database Metrics**: Query performance, connection pools
- **Network Metrics**: Latency, throughput, errors
- **External Services**: Payment gateway status

### **Business Intelligence**

- **Conversion Funnel**: Cart-to-enrollment conversion
- **Revenue Tracking**: Real-time revenue metrics
- **User Behavior**: Checkout flow analysis
- **Payment Analytics**: Success rate, method preferences

---

## 🚀 **Deployment Strategy**

### **Environment Setup**

- **Development**: Local development environment
- **Staging**: Production-like testing environment
- **Production**: Live production environment
- **Disaster Recovery**: Backup and recovery procedures

### **CI/CD Pipeline**

- **Code Quality**: Automated code review and testing
- **Security Scanning**: Vulnerability scanning
- **Performance Testing**: Automated performance validation
- **Deployment**: Blue-green deployment strategy

### **Rollback Strategy**

- **Feature Flags**: Gradual feature rollout
- **Database Migrations**: Backward-compatible migrations
- **Service Rollback**: Quick service rollback capability
- **Data Recovery**: Point-in-time data recovery

---

## 📋 **Implementation Timeline**

### **Phase 1: Foundation (Day 1 - Hours 1-4)**

- [ ] Database schema design and implementation
- [ ] Domain models and repository interfaces
- [ ] Basic service layer structure
- [ ] API endpoint scaffolding

### **Phase 2: Core Services (Day 1 - Hours 5-8)**

- [ ] Cart service implementation
- [ ] Order service implementation
- [ ] Payment service integration
- [ ] Basic API functionality

### **Phase 3: Payment Integration (Day 1 - Hours 9-12)**

- [ ] Razorpay provider completion
- [ ] Stripe provider implementation
- [ ] Webhook handling system
- [ ] Payment flow testing

### **Phase 4: Event System (Day 1 - Hours 13-16)**

- [ ] Event publishing implementation
- [ ] Enrollment worker implementation
- [ ] Idempotency middleware
- [ ] End-to-end testing

### **Phase 5: Production Readiness (Day 2 - Hours 1-8)**

- [ ] Comprehensive testing
- [ ] Performance optimization
- [ ] Security hardening
- [ ] Documentation completion

---

## 🎯 **Success Criteria**

### **Functional Success**

- [ ] Users can successfully add internships to cart
- [ ] Discount codes apply correctly to cart totals
- [ ] Orders are created with proper idempotency
- [ ] Payment processing works with multiple gateways
- [ ] Webhooks trigger enrollment creation
- [ ] No duplicate enrollments or charges occur

### **Performance Success**

- [ ] All API endpoints respond within target times
- [ ] System handles 10,000+ concurrent users
- [ ] Payment success rate exceeds 98%
- [ ] Webhook processing achieves 100% reliability

### **Business Success**

- [ ] Checkout abandonment rate below 5%
- [ ] Cart-to-enrollment conversion above 85%
- [ ] User satisfaction score above 4.5/5
- [ ] Revenue tracking accuracy 100%

---

## 🚨 **Risk Mitigation**

### **Technical Risks**

| Risk                     | Impact | Mitigation                        |
| ------------------------ | ------ | --------------------------------- |
| Payment Gateway Downtime | High   | Multiple gateway fallback         |
| Database Performance     | Medium | Query optimization, indexing      |
| Webhook Failures         | High   | Dead letter queue, retry logic    |
| Idempotency Issues       | High   | Comprehensive testing, monitoring |

### **Business Risks**

| Risk               | Impact   | Mitigation                   |
| ------------------ | -------- | ---------------------------- |
| User Abandonment   | High     | UX optimization, A/B testing |
| Payment Failures   | High     | Multiple payment methods     |
| Security Breaches  | Critical | Security audits, compliance  |
| Scalability Issues | Medium   | Load testing, auto-scaling   |

---

## 📚 **Documentation Requirements**

### **Technical Documentation**

- [ ] API documentation with OpenAPI/Swagger
- [ ] Database schema documentation
- [ ] Service architecture documentation
- [ ] Deployment and operations guide

### **User Documentation**

- [ ] Checkout flow user guide
- [ ] Payment method instructions
- [ ] Troubleshooting guide
- [ ] FAQ and support documentation

### **Business Documentation**

- [ ] Business process documentation
- [ ] Revenue tracking procedures
- [ ] Compliance documentation
- [ ] Audit trail procedures

---

This PRD provides a comprehensive blueprint for implementing a world-class checkout and enrollment system that will drive user satisfaction, increase conversion rates, and provide a solid foundation for business growth.
