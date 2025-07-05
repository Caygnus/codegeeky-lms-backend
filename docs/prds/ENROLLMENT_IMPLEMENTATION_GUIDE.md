# 🛠️ LMS Enrollment Workflow - Implementation Guide

## 📋 Overview

This guide provides step-by-step implementation instructions for the LMS enrollment workflow with Razorpay payment integration, building upon your existing architecture.

## 🏗️ Implementation Phases

### Phase 1: Database Schema Updates

#### 1.1 Update Enrollment Schema

Your existing enrollment schema is mostly complete. Add these enhancements:

```sql
-- Add indexes for performance
CREATE INDEX CONCURRENTLY idx_enrollments_user_course ON enrollments(user_id, internship_id) WHERE status != 'deleted';
CREATE INDEX CONCURRENTLY idx_enrollments_payment_lookup ON enrollments(payment_id) WHERE payment_id IS NOT NULL;

-- Add enrollment expiry support
ALTER TABLE enrollments ADD COLUMN enrollment_expires_at TIMESTAMP;
```

### Phase 2: Enhanced Service Implementation

#### 2.1 Create Enrollment Service

```go
// File: internal/service/enrollment.go
package service

import (
    "context"
    "fmt"
    "time"

    "github.com/omkar273/codegeeky/internal/api/dto"
    "github.com/omkar273/codegeeky/internal/domain/enrollment"
    ierr "github.com/omkar273/codegeeky/internal/errors"
    "github.com/omkar273/codegeeky/internal/types"
    "github.com/samber/lo"
    "github.com/shopspring/decimal"
)

type EnrollmentService interface {
    InitializeEnrollment(ctx context.Context, req *dto.InitializeEnrollmentRequest) (*dto.EnrollmentResponse, error)
    CompleteEnrollment(ctx context.Context, enrollmentID string) error
    GetEnrollmentStatus(ctx context.Context, enrollmentID string) (*dto.EnrollmentStatusResponse, error)
    CancelEnrollment(ctx context.Context, enrollmentID string, reason string) error
}

type enrollmentService struct {
    ServiceParams  ServiceParams
    PaymentService PaymentService
}

func NewEnrollmentService(params ServiceParams, paymentSvc PaymentService) EnrollmentService {
    return &enrollmentService{
        ServiceParams:  params,
        PaymentService: paymentSvc,
    }
}

func (s *enrollmentService) InitializeEnrollment(ctx context.Context, req *dto.InitializeEnrollmentRequest) (*dto.EnrollmentResponse, error) {
    // 1. Validate request
    if err := s.validateEnrollmentRequest(ctx, req); err != nil {
        return nil, err
    }

    // 2. Check existing enrollment
    existingEnrollment, err := s.checkExistingEnrollment(ctx, req.UserID, req.CourseID)
    if err != nil {
        return nil, err
    }
    if existingEnrollment != nil {
        return s.handleExistingEnrollment(ctx, existingEnrollment)
    }

    // 3. Get course details
    course, err := s.ServiceParams.InternshipRepo.Get(ctx, req.CourseID)
    if err != nil {
        return nil, ierr.WithError(err).
            WithHint("Course not found").
            Mark(ierr.ErrNotFound)
    }

    // 4. Calculate pricing
    pricing := s.calculatePricing(course, req.CouponCode)

    var response *dto.EnrollmentResponse

    // 5. Execute in transaction
    err = s.ServiceParams.DB.WithTx(ctx, func(ctx context.Context) error {
        // Create enrollment record
        enrollment, err := s.createEnrollmentRecord(ctx, req.UserID, req.CourseID, pricing)
        if err != nil {
            return err
        }

        // Handle free courses
        if pricing.FinalAmount.IsZero() {
            return s.handleFreeCourse(ctx, enrollment)
        }

        // Create payment for paid courses
        paymentResp, err := s.createPaymentOrder(ctx, enrollment, pricing, req)
        if err != nil {
            return err
        }

        response = &dto.EnrollmentResponse{
            EnrollmentID:    enrollment.ID,
            Status:          string(enrollment.EnrollmentStatus),
            PaymentRequired: true,
            Pricing:         pricing,
            PaymentSession:  paymentResp.PaymentSession,
        }

        return nil
    })

    return response, err
}
```

#### 2.2 Create DTOs

```go
// File: internal/api/dto/enrollment.go
package dto

import (
    "time"
    "github.com/shopspring/decimal"
)

type InitializeEnrollmentRequest struct {
    CourseID    string `json:"course_id" validate:"required"`
    UserID      string `json:"-"` // Set from context
    CouponCode  string `json:"coupon_code,omitempty"`
    SuccessURL  string `json:"success_url" validate:"required,url"`
    CancelURL   string `json:"cancel_url" validate:"required,url"`
    Metadata    map[string]string `json:"metadata,omitempty"`
}

type EnrollmentResponse struct {
    EnrollmentID    string              `json:"enrollment_id"`
    Status          string              `json:"status"`
    PaymentRequired bool                `json:"payment_required"`
    Pricing         *PricingInfo        `json:"pricing"`
    PaymentSession  *PaymentSessionInfo `json:"payment_session,omitempty"`
}

type PricingInfo struct {
    OriginalAmount  decimal.Decimal `json:"original_amount"`
    DiscountAmount  decimal.Decimal `json:"discount_amount"`
    FinalAmount     decimal.Decimal `json:"final_amount"`
    Currency        string          `json:"currency"`
    TaxAmount       decimal.Decimal `json:"tax_amount"`
    TotalPayable    decimal.Decimal `json:"total_payable"`
}

type PaymentSessionInfo struct {
    PaymentID       string    `json:"payment_id"`
    RazorpayOrderID string    `json:"razorpay_order_id"`
    RazorpayKey     string    `json:"razorpay_key"`
    PaymentURL      string    `json:"payment_url,omitempty"`
    ExpiresAt       time.Time `json:"expires_at"`
}

type EnrollmentStatusResponse struct {
    EnrollmentID     string    `json:"enrollment_id"`
    EnrollmentStatus string    `json:"enrollment_status"`
    PaymentStatus    string    `json:"payment_status"`
    PaymentID        string    `json:"payment_id,omitempty"`
    RazorpayPaymentID string   `json:"razorpay_payment_id,omitempty"`
    CompletedAt      *time.Time `json:"completed_at,omitempty"`
    CourseAccessURL  string    `json:"course_access_url,omitempty"`
}
```

### Phase 3: Enhanced Razorpay Integration

#### 3.1 Extend Razorpay Provider

```go
// File: internal/payment/providers/razorpay.go
// Add these methods to existing RazorpayProvider

func (r *RazorpayProvider) CreateEnrollmentOrder(ctx context.Context, req *dto.EnrollmentPaymentRequest) (*dto.PaymentResponse, error) {
    // Convert amount to paise
    amountInPaise := int(req.Amount.Mul(decimal.NewFromInt(100)).IntPart())

    orderData := map[string]interface{}{
        "amount":          amountInPaise,
        "currency":        req.Currency,
        "receipt":         req.IdempotencyKey,
        "payment_capture": 1,
        "notes": map[string]string{
            "enrollment_id": req.EnrollmentID,
            "course_id":     req.CourseID,
            "user_id":       req.UserID,
        },
    }

    order, err := r.razorpayClient.Order.Create(orderData, nil)
    if err != nil {
        return nil, fmt.Errorf("razorpay order creation failed: %w", err)
    }

    return &dto.PaymentResponse{
        Payment: payment.Payment{
            ID: order["id"].(string),
        },
        GatewayResponse: &dto.PaymentGatewayResponse{
            ProviderPaymentID: order["id"].(string),
            RedirectURL:       r.generateCheckoutURL(order["id"].(string)),
            Status:            "created",
            Raw:               order,
        },
    }, nil
}

func (r *RazorpayProvider) generateCheckoutURL(orderID string) string {
    // Return hosted checkout URL or return empty for custom integration
    return fmt.Sprintf("https://checkout.razorpay.com/v1/checkout.js?order_id=%s", orderID)
}
```

#### 3.2 Enhanced Webhook Processing

```go
// File: internal/webhook/subscriber/enrollment.go
package subscriber

import (
    "context"
    "encoding/json"

    "github.com/omkar273/codegeeky/internal/api/dto"
    "github.com/omkar273/codegeeky/internal/service"
    "github.com/omkar273/codegeeky/internal/types"
)

type EnrollmentWebhookSubscriber struct {
    enrollmentSvc service.EnrollmentService
    paymentSvc    service.PaymentService
}

func NewEnrollmentWebhookSubscriber(enrollmentSvc service.EnrollmentService, paymentSvc service.PaymentService) *EnrollmentWebhookSubscriber {
    return &EnrollmentWebhookSubscriber{
        enrollmentSvc: enrollmentSvc,
        paymentSvc:    paymentSvc,
    }
}

func (s *EnrollmentWebhookSubscriber) HandlePaymentSuccess(ctx context.Context, event *dto.WebhookResult) error {
    // Extract enrollment ID from payment metadata
    enrollmentID, ok := event.Payload["enrollment_id"].(string)
    if !ok {
        return fmt.Errorf("enrollment_id not found in payment webhook")
    }

    // Update payment status
    updateReq := &dto.UpdatePaymentRequest{
        PaymentStatus:     &types.PaymentStatusSuccess,
        GatewayPaymentID:  &event.PaymentID,
    }

    if _, err := s.paymentSvc.Update(ctx, event.PaymentID, updateReq); err != nil {
        return fmt.Errorf("failed to update payment: %w", err)
    }

    // Complete enrollment
    return s.enrollmentSvc.CompleteEnrollment(ctx, enrollmentID)
}

func (s *EnrollmentWebhookSubscriber) HandlePaymentFailure(ctx context.Context, event *dto.WebhookResult) error {
    enrollmentID, ok := event.Payload["enrollment_id"].(string)
    if !ok {
        return fmt.Errorf("enrollment_id not found in payment webhook")
    }

    // Update payment status
    updateReq := &dto.UpdatePaymentRequest{
        PaymentStatus:    &types.PaymentStatusFailed,
        GatewayPaymentID: &event.PaymentID,
        ErrorMessage:     &event.Reason,
    }

    if _, err := s.paymentSvc.Update(ctx, event.PaymentID, updateReq); err != nil {
        return fmt.Errorf("failed to update payment: %w", err)
    }

    // Cancel enrollment
    return s.enrollmentSvc.CancelEnrollment(ctx, enrollmentID, "payment_failed")
}
```

### Phase 4: API Endpoints

#### 4.1 Enrollment Controller

```go
// File: internal/api/v1/enrollment.go
package v1

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/omkar273/codegeeky/internal/api/dto"
    "github.com/omkar273/codegeeky/internal/service"
    "github.com/omkar273/codegeeky/internal/types"
)

type EnrollmentHandler struct {
    enrollmentSvc service.EnrollmentService
}

func NewEnrollmentHandler(enrollmentSvc service.EnrollmentService) *EnrollmentHandler {
    return &EnrollmentHandler{enrollmentSvc: enrollmentSvc}
}

// @Summary Initialize course enrollment
// @Description Start the enrollment process for a course
// @Tags enrollments
// @Accept json
// @Produce json
// @Param request body dto.InitializeEnrollmentRequest true "Enrollment request"
// @Success 200 {object} dto.EnrollmentResponse
// @Router /enrollments/initialize [post]
func (h *EnrollmentHandler) InitializeEnrollment(c *gin.Context) {
    var req dto.InitializeEnrollmentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Set user ID from context
    req.UserID = types.GetUserID(c)

    response, err := h.enrollmentSvc.InitializeEnrollment(c, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, response)
}

// @Summary Get enrollment status
// @Description Get the current status of an enrollment
// @Tags enrollments
// @Produce json
// @Param enrollmentId path string true "Enrollment ID"
// @Success 200 {object} dto.EnrollmentStatusResponse
// @Router /enrollments/{enrollmentId}/status [get]
func (h *EnrollmentHandler) GetEnrollmentStatus(c *gin.Context) {
    enrollmentID := c.Param("enrollmentId")

    response, err := h.enrollmentSvc.GetEnrollmentStatus(c, enrollmentID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, response)
}
```

#### 4.2 Router Setup

```go
// File: internal/api/router.go
// Add to existing router setup

func (r *Router) setupEnrollmentRoutes() {
    enrollmentGroup := r.engine.Group("/api/v1/enrollments")
    enrollmentGroup.Use(r.middleware.AuthRequired())

    enrollmentHandler := v1.NewEnrollmentHandler(r.services.EnrollmentService())

    enrollmentGroup.POST("/initialize", enrollmentHandler.InitializeEnrollment)
    enrollmentGroup.GET("/:enrollmentId/status", enrollmentHandler.GetEnrollmentStatus)
}
```

### Phase 5: Frontend Integration

#### 5.1 JavaScript/TypeScript Integration

```typescript
// enrollment.service.ts
export interface EnrollmentRequest {
  course_id: string;
  coupon_code?: string;
  success_url: string;
  cancel_url: string;
  metadata?: Record<string, string>;
}

export interface EnrollmentResponse {
  enrollment_id: string;
  status: string;
  payment_required: boolean;
  pricing: PricingInfo;
  payment_session?: PaymentSessionInfo;
}

export class EnrollmentService {
  constructor(private apiClient: ApiClient) {}

  async initializeEnrollment(
    request: EnrollmentRequest
  ): Promise<EnrollmentResponse> {
    const response = await this.apiClient.post(
      "/api/v1/enrollments/initialize",
      request
    );
    return response.data;
  }

  async getEnrollmentStatus(
    enrollmentId: string
  ): Promise<EnrollmentStatusResponse> {
    const response = await this.apiClient.get(
      `/api/v1/enrollments/${enrollmentId}/status`
    );
    return response.data;
  }

  async processRazorpayPayment(
    enrollmentData: EnrollmentResponse
  ): Promise<void> {
    return new Promise((resolve, reject) => {
      const options = {
        key: enrollmentData.payment_session!.razorpay_key,
        order_id: enrollmentData.payment_session!.razorpay_order_id,
        amount: Math.round(enrollmentData.pricing.total_payable * 100),
        currency: enrollmentData.pricing.currency,
        name: "Your Platform Name",
        description: "Course Enrollment",
        handler: async (response: any) => {
          try {
            await this.verifyPayment({
              razorpay_payment_id: response.razorpay_payment_id,
              razorpay_order_id: response.razorpay_order_id,
              razorpay_signature: response.razorpay_signature,
              enrollment_id: enrollmentData.enrollment_id,
            });
            resolve();
          } catch (error) {
            reject(error);
          }
        },
        modal: {
          ondismiss: () => reject(new Error("Payment cancelled")),
        },
      };

      const rzp = new (window as any).Razorpay(options);
      rzp.open();
    });
  }

  private async verifyPayment(verificationData: any): Promise<void> {
    await this.apiClient.post("/api/v1/payments/verify", verificationData);
  }
}
```

#### 5.2 React Component Example

```tsx
// EnrollmentButton.tsx
import React, { useState } from "react";
import { EnrollmentService } from "./enrollment.service";

interface EnrollmentButtonProps {
  courseId: string;
  courseName: string;
  price: number;
}

export const EnrollmentButton: React.FC<EnrollmentButtonProps> = ({
  courseId,
  courseName,
  price,
}) => {
  const [loading, setLoading] = useState(false);
  const [enrollmentStatus, setEnrollmentStatus] = useState<string>("");
  const enrollmentService = new EnrollmentService();

  const handleEnrollment = async () => {
    setLoading(true);
    try {
      // Initialize enrollment
      const enrollmentResponse = await enrollmentService.initializeEnrollment({
        course_id: courseId,
        success_url: `${window.location.origin}/enrollment/success`,
        cancel_url: `${window.location.origin}/enrollment/cancel`,
      });

      if (!enrollmentResponse.payment_required) {
        // Free course - enrollment complete
        setEnrollmentStatus("enrolled");
        return;
      }

      // Process payment
      await enrollmentService.processRazorpayPayment(enrollmentResponse);

      // Check final status
      const statusResponse = await enrollmentService.getEnrollmentStatus(
        enrollmentResponse.enrollment_id
      );
      setEnrollmentStatus(statusResponse.enrollment_status);
    } catch (error) {
      console.error("Enrollment failed:", error);
      alert("Enrollment failed. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  if (enrollmentStatus === "enrolled") {
    return <button disabled>Enrolled ✓</button>;
  }

  return (
    <button
      onClick={handleEnrollment}
      disabled={loading}
      className="enrollment-button"
    >
      {loading ? "Processing..." : `Enroll for ₹${price}`}
    </button>
  );
};
```

## 🧪 Testing Strategy

### Integration Tests

```go
// File: tests/integration/enrollment_test.go
package integration

import (
    "testing"
    "context"

    "github.com/stretchr/testify/assert"
    "github.com/omkar273/codegeeky/internal/api/dto"
)

func TestEnrollmentWorkflow(t *testing.T) {
    // Setup test environment
    testEnv := setupTestEnvironment(t)
    defer testEnv.Cleanup()

    // Test data
    userID := "usr_test123"
    courseID := "int_test123"

    t.Run("Successful Paid Course Enrollment", func(t *testing.T) {
        // Initialize enrollment
        req := &dto.InitializeEnrollmentRequest{
            CourseID:   courseID,
            UserID:     userID,
            SuccessURL: "http://test.com/success",
            CancelURL:  "http://test.com/cancel",
        }

        response, err := testEnv.EnrollmentService.InitializeEnrollment(context.Background(), req)
        assert.NoError(t, err)
        assert.True(t, response.PaymentRequired)
        assert.NotEmpty(t, response.EnrollmentID)

        // Simulate successful payment webhook
        webhookEvent := createPaymentSuccessWebhook(response.PaymentSession.PaymentID, response.EnrollmentID)
        err = testEnv.WebhookService.ProcessWebhook(context.Background(), webhookEvent)
        assert.NoError(t, err)

        // Verify enrollment completion
        status, err := testEnv.EnrollmentService.GetEnrollmentStatus(context.Background(), response.EnrollmentID)
        assert.NoError(t, err)
        assert.Equal(t, "enrolled", status.EnrollmentStatus)
    })

    t.Run("Free Course Enrollment", func(t *testing.T) {
        // Create free course
        freeCourseID := createFreeCourse(t, testEnv)

        req := &dto.InitializeEnrollmentRequest{
            CourseID:   freeCourseID,
            UserID:     userID,
            SuccessURL: "http://test.com/success",
            CancelURL:  "http://test.com/cancel",
        }

        response, err := testEnv.EnrollmentService.InitializeEnrollment(context.Background(), req)
        assert.NoError(t, err)
        assert.False(t, response.PaymentRequired)

        // Verify immediate enrollment
        status, err := testEnv.EnrollmentService.GetEnrollmentStatus(context.Background(), response.EnrollmentID)
        assert.NoError(t, err)
        assert.Equal(t, "enrolled", status.EnrollmentStatus)
    })
}
```

## 📊 Monitoring & Metrics

### Key Metrics to Track

```go
// File: internal/metrics/enrollment.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    EnrollmentInitiations = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "enrollment_initiations_total",
            Help: "Total number of enrollment initiations",
        },
        []string{"course_id", "payment_required"},
    )

    EnrollmentCompletions = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "enrollment_completions_total",
            Help: "Total number of successful enrollments",
        },
        []string{"course_id", "payment_method"},
    )

    PaymentProcessingDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "payment_processing_duration_seconds",
            Help: "Time taken to process payments",
        },
        []string{"gateway", "status"},
    )

    EnrollmentConversionRate = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "enrollment_conversion_rate",
            Help: "Enrollment conversion rate by course",
        },
        []string{"course_id"},
    )
)
```

## 🚀 Deployment Checklist

- [ ] Database migrations executed
- [ ] Environment variables configured
- [ ] Razorpay credentials set up
- [ ] Webhook endpoints configured
- [ ] Monitoring dashboards deployed
- [ ] Load testing completed
- [ ] Security audit passed
- [ ] Documentation updated

## 🔍 Troubleshooting Guide

### Common Issues

1. **Payment Webhook Not Received**

   - Check webhook URL accessibility
   - Verify webhook signature validation
   - Check firewall rules

2. **Enrollment Status Not Updating**

   - Check database transactions
   - Verify webhook processing logic
   - Check for race conditions

3. **Razorpay Order Creation Fails**
   - Verify API credentials
   - Check amount formatting (paise conversion)
   - Validate request parameters

This implementation guide provides a practical roadmap for building the enrollment workflow on top of your existing architecture. Each phase builds incrementally, allowing for testing and validation at each step.
