# Payment Service

## Overview

The Payment Service has been sanitized to focus only on CRUD operations for the payments table. All third-party payment gateway logic has been removed to make it reusable across the system.

## Features

### Core Payment Operations

- **Create**: Create new payment records with idempotency support
- **GetByID**: Retrieve payments by ID
- **GetByIdempotencyKey**: Retrieve payments by idempotency key
- **Update**: Update existing payment records
- **Delete**: Delete payment records
- **List**: Retrieve paginated list of payments with filtering

### Payment Attempt Operations

- **CreateAttempt**: Create payment attempt records
- **GetAttempt**: Retrieve payment attempts by ID
- **ListAttempts**: Get all attempts for a payment
- **GetLatestAttempt**: Get the latest attempt for a payment

### Status Management

- **UpdateStatus**: Update payment status with metadata
- **MarkAsSuccess**: Mark payment as successful
- **MarkAsFailed**: Mark payment as failed with error message
- **MarkAsRefunded**: Mark payment as refunded

## Usage

### Basic Payment Creation

```go
// Create a payment service
paymentService := NewPaymentService(serviceParams)

// Create a payment request
req := &dto.CreatePaymentRequest{
    PaymentRequest: dto.PaymentRequest{
        ReferenceID:            "enrollment-123",
        ReferenceType:          types.PaymentDestinationTypeEnrollment,
        DestinationID:          "internship-456",
        DestinationType:        types.PaymentDestinationTypeInternship,
        Amount:                 10000, // 100.00 in paisa
        Currency:               "INR",
        PaymentGatewayProvider: types.PaymentGatewayProviderRazorpay,
        PaymentMethodType:      &types.PaymentMethodTypeCard,
        IdempotencyKey:         "unique-key-123",
        Metadata:               map[string]string{"source": "web"},
        TrackAttempts:          true,
    },
}

// Create the payment
payment, err := paymentService.Create(ctx, req)
if err != nil {
    // Handle error
}
```

### Status Updates

```go
// Mark payment as successful
payment, err := paymentService.MarkAsSuccess(ctx, paymentID, &gatewayPaymentID, metadata)

// Mark payment as failed
payment, err := paymentService.MarkAsFailed(ctx, paymentID, "Payment declined", metadata)

// Mark payment as refunded
payment, err := paymentService.MarkAsRefunded(ctx, paymentID, metadata)
```

### Payment Attempts

```go
// Create a payment attempt
attemptReq := &dto.PaymentAttemptRequest{
    PaymentID:     paymentID,
    PaymentStatus: types.PaymentStatusPending,
    Metadata:      types.MetadataFromEnt(payment.Metadata),
}

attempt, err := paymentService.CreateAttempt(ctx, attemptReq)

// Get latest attempt
latestAttempt, err := paymentService.GetLatestAttempt(ctx, paymentID)
```

## Integration with External Payment Processing

The sanitized payment service is designed to work with external payment processing systems:

1. **Create Payment Record**: Use the service to create a payment record in your database
2. **External Processing**: Handle actual payment processing in a separate service/module
3. **Update Status**: Use the service's status update methods to reflect the payment result

### Example Workflow

```go
// 1. Create payment record
payment, err := paymentService.Create(ctx, createReq)

// 2. Send to external payment processor
// (This would be handled by your external payment service)
externalPaymentID, err := externalPaymentProcessor.Process(payment.ID, payment.Amount)

// 3. Update payment status based on external processor response
if err != nil {
    paymentService.MarkAsFailed(ctx, payment.ID, err.Error(), nil)
} else {
    paymentService.MarkAsSuccess(ctx, payment.ID, &externalPaymentID, nil)
}
```

## Dependencies

The payment service depends on:

- `ServiceParams` (database, logger, repositories)
- Payment repository for data persistence
- No external payment gateway dependencies

## Benefits

1. **Reusable**: Can be used across different parts of your system
2. **Testable**: Easy to unit test without external dependencies
3. **Flexible**: Can integrate with any payment processing system
4. **Idempotent**: Supports safe retries with idempotency keys
5. **Auditable**: Tracks payment attempts and status changes

## Future Enhancements

When you're ready to implement real payment processing:

1. Create a separate payment processing service
2. Use the gateway registry (`gateway.GatewayRegistryService`) for third-party integrations
3. Integrate with the payment service for status updates
4. Add webhook handling for payment status updates
