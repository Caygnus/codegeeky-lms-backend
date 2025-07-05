# 🚀 LMS Enrollment API Examples & Test Scenarios

## 📋 Overview

This document provides comprehensive API examples, request/response samples, and test scenarios for the LMS enrollment workflow with Razorpay integration.

## 🔄 Complete Enrollment Flow Examples

### Scenario 1: Paid Course Enrollment

#### Step 1: Get Course Information

```http
GET /api/v1/internships/int_01HKZM7X9P0QJ2K3L4M5N6O7P8
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**Response:**

```json
{
  "internship": {
    "id": "int_01HKZM7X9P0QJ2K3L4M5N6O7P8",
    "title": "Full Stack Web Development Bootcamp",
    "description": "Comprehensive 12-week program covering React, Node.js, and MongoDB",
    "price": "4999.00",
    "currency": "INR",
    "flat_discount": "2000.00",
    "percentage_discount": "0.00",
    "final_price": "2999.00",
    "duration_in_weeks": 12,
    "level": "intermediate",
    "prerequisites": ["basic_programming", "html_css"],
    "skills": ["react", "nodejs", "mongodb", "javascript"],
    "categories": [
      {
        "id": "cat_web_development",
        "name": "Web Development"
      }
    ],
    "enrollment_info": {
      "can_enroll": true,
      "is_enrolled": false,
      "enrollment_deadline": "2024-03-30T23:59:59Z",
      "available_spots": 25,
      "enrolled_count": 45
    }
  }
}
```

#### Step 2: Initialize Enrollment

```http
POST /api/v1/enrollments/initialize
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

{
  "course_id": "int_01HKZM7X9P0QJ2K3L4M5N6O7P8",
  "coupon_code": "EARLY50",
  "payment_method_preference": "razorpay",
  "success_url": "https://myapp.com/enrollment/success",
  "cancel_url": "https://myapp.com/enrollment/cancel",
  "metadata": {
    "source": "web_app",
    "campaign": "spring_2024",
    "referrer": "google_ads"
  }
}
```

**Response:**

```json
{
  "enrollment_id": "enr_01HKZM8Y1A2B3C4D5E6F7G8H9I",
  "status": "pending",
  "payment_required": true,
  "course": {
    "id": "int_01HKZM7X9P0QJ2K3L4M5N6O7P8",
    "title": "Full Stack Web Development Bootcamp"
  },
  "pricing": {
    "original_amount": "4999.00",
    "discount_amount": "2500.00",
    "coupon_discount": "500.00",
    "final_amount": "2499.00",
    "currency": "INR",
    "tax_amount": "449.82",
    "total_payable": "2948.82"
  },
  "payment_session": {
    "payment_id": "pay_01HKZM8Z2B3C4D5E6F7G8H9I0J",
    "razorpay_order_id": "order_NfJZ5mUg7MUlEe",
    "razorpay_key": "rzp_test_1DP5mmOlF5G5ag",
    "checkout_url": "https://checkout.razorpay.com/v1/checkout.js",
    "expires_at": "2024-01-15T11:30:00Z"
  },
  "idempotency_key": "enr_01HKZM8Y1A2B3C4D5E6F7G8H9I_1705312200"
}
```

#### Step 3: Frontend Payment Processing

```javascript
// Initialize Razorpay payment
const options = {
  key: "rzp_test_1DP5mmOlF5G5ag",
  order_id: "order_NfJZ5mUg7MUlEe",
  amount: 294882, // in paise
  currency: "INR",
  name: "Your LMS Platform",
  description: "Full Stack Web Development Bootcamp",
  image: "https://yourapp.com/logo.png",
  handler: function (response) {
    // Payment successful
    console.log("Payment ID:", response.razorpay_payment_id);
    console.log("Order ID:", response.razorpay_order_id);
    console.log("Signature:", response.razorpay_signature);

    // Verify payment on backend
    verifyPayment(response);
  },
  prefill: {
    name: "John Doe",
    email: "john@example.com",
    contact: "9876543210",
  },
  theme: {
    color: "#3399cc",
  },
  modal: {
    ondismiss: function () {
      console.log("Payment cancelled");
      handlePaymentCancellation();
    },
  },
};

const rzp = new Razorpay(options);
rzp.open();
```

#### Step 4: Payment Verification (Backend Webhook)

```json
// Razorpay webhook payload
{
  "entity": "event",
  "account_id": "acc_BFQ7uQEaa30GJJ",
  "event": "payment.captured",
  "contains": ["payment"],
  "payload": {
    "payment": {
      "entity": "payment",
      "id": "pay_NfJZ6mUg7MUlEf",
      "amount": 294882,
      "currency": "INR",
      "status": "captured",
      "order_id": "order_NfJZ5mUg7MUlEe",
      "method": "card",
      "description": "Full Stack Web Development Bootcamp",
      "notes": {
        "enrollment_id": "enr_01HKZM8Y1A2B3C4D5E6F7G8H9I",
        "course_id": "int_01HKZM7X9P0QJ2K3L4M5N6O7P8",
        "user_id": "usr_01HKZM8X0Z1Y2X3W4V5U6T7S8R"
      },
      "created_at": 1705312800
    }
  },
  "created_at": 1705312800
}
```

#### Step 5: Check Enrollment Status

```http
GET /api/v1/enrollments/enr_01HKZM8Y1A2B3C4D5E6F7G8H9I/status
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**Response:**

```json
{
  "enrollment_id": "enr_01HKZM8Y1A2B3C4D5E6F7G8H9I",
  "enrollment_status": "enrolled",
  "payment_status": "completed",
  "payment_id": "pay_01HKZM8Z2B3C4D5E6F7G8H9I0J",
  "razorpay_payment_id": "pay_NfJZ6mUg7MUlEf",
  "completed_at": "2024-01-15T10:40:00Z",
  "course_access": {
    "has_access": true,
    "access_granted_at": "2024-01-15T10:40:00Z",
    "next_lesson_url": "/api/v1/courses/int_01HKZM7X9P0QJ2K3L4M5N6O7P8/lessons/1"
  }
}
```

### Scenario 2: Free Course Enrollment

#### Initialize Free Course Enrollment

```http
POST /api/v1/enrollments/initialize
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

{
  "course_id": "int_01HKZM7X9P0QJ2K3L4M5N6O7P9",
  "success_url": "https://myapp.com/enrollment/success",
  "cancel_url": "https://myapp.com/enrollment/cancel"
}
```

**Response:**

```json
{
  "enrollment_id": "enr_01HKZM8Y1A2B3C4D5E6F7G8H9J",
  "status": "enrolled",
  "payment_required": false,
  "course": {
    "id": "int_01HKZM7X9P0QJ2K3L4M5N6O7P9",
    "title": "Introduction to Programming"
  },
  "pricing": {
    "original_amount": "0.00",
    "discount_amount": "0.00",
    "final_amount": "0.00",
    "currency": "INR",
    "tax_amount": "0.00",
    "total_payable": "0.00"
  },
  "course_access": {
    "has_access": true,
    "access_granted_at": "2024-01-15T10:30:00Z",
    "next_lesson_url": "/api/v1/courses/int_01HKZM7X9P0QJ2K3L4M5N6O7P9/lessons/1"
  }
}
```

## 🚨 Error Handling Examples

### Error 1: Already Enrolled

```http
POST /api/v1/enrollments/initialize
```

**Response (409 Conflict):**

```json
{
  "error": {
    "code": "already_enrolled",
    "message": "User is already enrolled in this course",
    "details": {
      "enrollment_id": "enr_01HKZM8Y1A2B3C4D5E6F7G8H9K",
      "enrollment_status": "enrolled",
      "enrolled_at": "2024-01-10T14:30:00Z"
    },
    "actions": [
      {
        "type": "redirect",
        "url": "/api/v1/courses/int_01HKZM7X9P0QJ2K3L4M5N6O7P8/access",
        "label": "Continue Learning"
      }
    ]
  }
}
```

### Error 2: Course Not Available

```http
POST /api/v1/enrollments/initialize
```

**Response (422 Unprocessable Entity):**

```json
{
  "error": {
    "code": "course_not_available",
    "message": "Course is not available for enrollment",
    "details": {
      "course_id": "int_01HKZM7X9P0QJ2K3L4M5N6O7P8",
      "status": "draft",
      "available_from": "2024-02-01T00:00:00Z"
    },
    "actions": [
      {
        "type": "notify",
        "label": "Notify me when available"
      }
    ]
  }
}
```

### Error 3: Payment Failed

```http
GET /api/v1/enrollments/enr_01HKZM8Y1A2B3C4D5E6F7G8H9I/status
```

**Response:**

```json
{
  "enrollment_id": "enr_01HKZM8Y1A2B3C4D5E6F7G8H9I",
  "enrollment_status": "failed",
  "payment_status": "failed",
  "payment_id": "pay_01HKZM8Z2B3C4D5E6F7G8H9I0J",
  "razorpay_payment_id": "pay_NfJZ6mUg7MUlEf",
  "failed_at": "2024-01-15T10:35:00Z",
  "error": {
    "code": "payment_failed",
    "message": "Payment was declined by the bank",
    "razorpay_error": "BAD_REQUEST_ERROR",
    "retry_options": [
      {
        "type": "retry_payment",
        "url": "/api/v1/enrollments/enr_01HKZM8Y1A2B3C4D5E6F7G8H9I/retry",
        "expires_at": "2024-01-16T10:35:00Z"
      },
      {
        "type": "change_payment_method",
        "url": "/api/v1/enrollments/enr_01HKZM8Y1A2B3C4D5E6F7G8H9I/update-payment"
      }
    ]
  }
}
```

## 🔄 Advanced Scenarios

### Scenario 3: Coupon Application

#### Apply Coupon Code

```http
POST /api/v1/coupons/validate
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

{
  "coupon_code": "EARLY50",
  "course_id": "int_01HKZM7X9P0QJ2K3L4M5N6O7P8",
  "user_id": "usr_01HKZM8X0Z1Y2X3W4V5U6T7S8R"
}
```

**Response:**

```json
{
  "valid": true,
  "coupon": {
    "code": "EARLY50",
    "type": "percentage",
    "value": "20.00",
    "description": "Early Bird 20% Discount",
    "applicable_courses": ["int_01HKZM7X9P0QJ2K3L4M5N6O7P8"],
    "expires_at": "2024-02-29T23:59:59Z",
    "usage_left": 89
  },
  "discount_calculation": {
    "original_amount": "4999.00",
    "coupon_discount": "999.80",
    "other_discounts": "2000.00",
    "total_discount": "2999.80",
    "final_amount": "1999.20"
  }
}
```

### Scenario 4: Bulk Enrollment (Enterprise)

#### Initialize Bulk Enrollment

```http
POST /api/v1/enrollments/bulk/initialize
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

{
  "course_id": "int_01HKZM7X9P0QJ2K3L4M5N6O7P8",
  "user_emails": [
    "employee1@company.com",
    "employee2@company.com",
    "employee3@company.com"
  ],
  "billing_details": {
    "company_name": "Acme Corp",
    "billing_email": "billing@acmecorp.com",
    "purchase_order": "PO-2024-001",
    "payment_terms": "net_30"
  },
  "success_url": "https://acmecorp.com/enrollment/success",
  "cancel_url": "https://acmecorp.com/enrollment/cancel"
}
```

**Response:**

```json
{
  "bulk_enrollment_id": "bulk_01HKZM8Y1A2B3C4D5E6F7G8H9K",
  "status": "pending",
  "enrollments": [
    {
      "enrollment_id": "enr_01HKZM8Y1A2B3C4D5E6F7G8H9L",
      "user_email": "employee1@company.com",
      "status": "pending"
    },
    {
      "enrollment_id": "enr_01HKZM8Y1A2B3C4D5E6F7G8H9M",
      "user_email": "employee2@company.com",
      "status": "pending"
    },
    {
      "enrollment_id": "enr_01HKZM8Y1A2B3C4D5E6F7G8H9N",
      "user_email": "employee3@company.com",
      "status": "pending"
    }
  ],
  "pricing": {
    "per_user_amount": "2999.00",
    "user_count": 3,
    "subtotal": "8997.00",
    "bulk_discount": "899.70",
    "tax_amount": "1457.55",
    "total_amount": "9554.85",
    "currency": "INR"
  },
  "invoice": {
    "invoice_id": "inv_01HKZM8Y1A2B3C4D5E6F7G8H9O",
    "invoice_url": "/api/v1/invoices/inv_01HKZM8Y1A2B3C4D5E6F7G8H9O/download",
    "due_date": "2024-02-14T23:59:59Z"
  }
}
```

### Scenario 5: Payment Retry

#### Retry Failed Payment

```http
POST /api/v1/enrollments/enr_01HKZM8Y1A2B3C4D5E6F7G8H9I/retry
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

{
  "payment_method_preference": "razorpay",
  "retry_reason": "user_requested"
}
```

**Response:**

```json
{
  "enrollment_id": "enr_01HKZM8Y1A2B3C4D5E6F7G8H9I",
  "status": "pending_payment",
  "retry_attempt": 2,
  "payment_session": {
    "payment_id": "pay_01HKZM8Z2B3C4D5E6F7G8H9I0K",
    "razorpay_order_id": "order_NfJZ7mUg7MUlEg",
    "razorpay_key": "rzp_test_1DP5mmOlF5G5ag",
    "expires_at": "2024-01-15T12:30:00Z"
  },
  "pricing": {
    "final_amount": "2948.82",
    "currency": "INR"
  }
}
```

## 🧪 Test Scenarios

### Test Case 1: Complete Happy Path

```javascript
// Test: Successful paid course enrollment
async function testPaidCourseEnrollment() {
  // 1. Get course info
  const courseResponse = await api.get("/api/v1/internships/int_test123");
  expect(courseResponse.status).toBe(200);
  expect(courseResponse.data.internship.price).toBe("4999.00");

  // 2. Initialize enrollment
  const enrollmentResponse = await api.post("/api/v1/enrollments/initialize", {
    course_id: "int_test123",
    success_url: "http://test.com/success",
    cancel_url: "http://test.com/cancel",
  });

  expect(enrollmentResponse.status).toBe(200);
  expect(enrollmentResponse.data.payment_required).toBe(true);
  expect(
    enrollmentResponse.data.payment_session.razorpay_order_id
  ).toBeDefined();

  // 3. Simulate successful payment webhook
  const webhookPayload = createSuccessfulPaymentWebhook(
    enrollmentResponse.data.payment_session.razorpay_order_id,
    enrollmentResponse.data.enrollment_id
  );

  await simulateWebhook("/webhooks/razorpay", webhookPayload);

  // 4. Verify enrollment completion
  const statusResponse = await api.get(
    `/api/v1/enrollments/${enrollmentResponse.data.enrollment_id}/status`
  );

  expect(statusResponse.data.enrollment_status).toBe("enrolled");
  expect(statusResponse.data.course_access.has_access).toBe(true);
}
```

### Test Case 2: Free Course Enrollment

```javascript
async function testFreeCourseEnrollment() {
  const response = await api.post("/api/v1/enrollments/initialize", {
    course_id: "int_free123",
    success_url: "http://test.com/success",
    cancel_url: "http://test.com/cancel",
  });

  expect(response.status).toBe(200);
  expect(response.data.payment_required).toBe(false);
  expect(response.data.status).toBe("enrolled");
}
```

### Test Case 3: Duplicate Enrollment

```javascript
async function testDuplicateEnrollment() {
  // First enrollment
  await api.post("/api/v1/enrollments/initialize", {
    course_id: "int_test123",
  });

  // Second enrollment attempt
  const response = await api.post("/api/v1/enrollments/initialize", {
    course_id: "int_test123",
  });

  expect(response.status).toBe(409);
  expect(response.data.error.code).toBe("already_enrolled");
}
```

### Test Case 4: Payment Failure Recovery

```javascript
async function testPaymentFailureRecovery() {
  // Initialize enrollment
  const enrollmentResponse = await api.post("/api/v1/enrollments/initialize", {
    course_id: "int_test123",
  });

  // Simulate failed payment webhook
  const failureWebhook = createFailedPaymentWebhook(
    enrollmentResponse.data.payment_session.razorpay_order_id,
    enrollmentResponse.data.enrollment_id
  );

  await simulateWebhook("/webhooks/razorpay", failureWebhook);

  // Check failure status
  const statusResponse = await api.get(
    `/api/v1/enrollments/${enrollmentResponse.data.enrollment_id}/status`
  );

  expect(statusResponse.data.enrollment_status).toBe("failed");
  expect(statusResponse.data.error.retry_options).toBeDefined();

  // Retry payment
  const retryResponse = await api.post(
    `/api/v1/enrollments/${enrollmentResponse.data.enrollment_id}/retry`
  );

  expect(retryResponse.status).toBe(200);
  expect(retryResponse.data.payment_session.razorpay_order_id).toBeDefined();
}
```

## 📊 Performance Testing

### Load Test Scenarios

#### 1. Concurrent Enrollments

```javascript
// Test: 100 concurrent enrollment requests
const concurrentEnrollments = Array.from({ length: 100 }, (_, i) =>
  api.post("/api/v1/enrollments/initialize", {
    course_id: "int_popular123",
    user_id: `usr_test${i}`,
    success_url: "http://test.com/success",
    cancel_url: "http://test.com/cancel",
  })
);

const results = await Promise.allSettled(concurrentEnrollments);
const successCount = results.filter((r) => r.status === "fulfilled").length;
const averageResponseTime = calculateAverageResponseTime(results);

expect(successCount).toBeGreaterThan(95); // 95% success rate
expect(averageResponseTime).toBeLessThan(500); // Under 500ms
```

#### 2. Webhook Processing Load

```javascript
// Test: Process 1000 webhook events in parallel
const webhookEvents = Array.from({ length: 1000 }, (_, i) =>
  createPaymentWebhook(`order_${i}`, `enr_${i}`)
);

const startTime = Date.now();
const results = await Promise.allSettled(
  webhookEvents.map((event) => simulateWebhook("/webhooks/razorpay", event))
);
const processingTime = Date.now() - startTime;

expect(results.filter((r) => r.status === "fulfilled").length).toBe(1000);
expect(processingTime).toBeLessThan(10000); // Under 10 seconds
```

## 🔒 Security Testing

### Security Test Cases

#### 1. Webhook Signature Validation

```javascript
async function testWebhookSignatureValidation() {
  const validWebhook = createValidWebhook();
  const invalidWebhook = { ...validWebhook };
  delete invalidWebhook.headers["x-razorpay-signature"];

  // Valid webhook should be processed
  const validResponse = await simulateWebhook(
    "/webhooks/razorpay",
    validWebhook
  );
  expect(validResponse.status).toBe(200);

  // Invalid webhook should be rejected
  const invalidResponse = await simulateWebhook(
    "/webhooks/razorpay",
    invalidWebhook
  );
  expect(invalidResponse.status).toBe(401);
}
```

#### 2. Idempotency Testing

```javascript
async function testIdempotency() {
  const enrollmentRequest = {
    course_id: "int_test123",
    idempotency_key: "test_idempotency_123",
  };

  // First request
  const response1 = await api.post(
    "/api/v1/enrollments/initialize",
    enrollmentRequest
  );
  expect(response1.status).toBe(200);

  // Duplicate request with same idempotency key
  const response2 = await api.post(
    "/api/v1/enrollments/initialize",
    enrollmentRequest
  );
  expect(response2.status).toBe(200);
  expect(response2.data.enrollment_id).toBe(response1.data.enrollment_id);
}
```

This comprehensive API examples document provides practical request/response samples, error handling patterns, and test scenarios to help implement and validate the enrollment workflow system.
