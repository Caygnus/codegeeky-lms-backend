# API Integration Examples

This document provides comprehensive examples for integrating with the event-driven pubsub system.

## Table of Contents

1. [Razorpay Integration](#razorpay-integration)
2. [WebSocket Integration](#websocket-integration)
3. [Microservice Communication](#microservice-communication)
4. [Event Publishing Patterns](#event-publishing-patterns)
5. [Error Handling Examples](#error-handling-examples)

## Razorpay Integration

### Setting up Razorpay Webhook Endpoint

```go
// internal/api/handlers/webhook.go
package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/omkar273/codegeeky/internal/webhook/receivers"
)

type WebhookHandler struct {
    razorpayReceiver *receivers.RazorpayReceiver
}

func NewWebhookHandler(razorpayReceiver *receivers.RazorpayReceiver) *WebhookHandler {
    return &WebhookHandler{
        razorpayReceiver: razorpayReceiver,
    }
}

func (h *WebhookHandler) SetupRoutes(router *gin.Engine) {
    webhooks := router.Group("/api/webhooks")
    {
        webhooks.POST("/razorpay", h.handleRazorpayWebhook)
    }
}

func (h *WebhookHandler) handleRazorpayWebhook(c *gin.Context) {
    h.razorpayReceiver.HandleWebhook(c.Writer, c.Request)
}
```

### Razorpay Event Processing

```go
// internal/webhook/handlers/razorpay_handler.go
package handlers

import (
    "context"
    "encoding/json"

    "github.com/ThreeDotsLabs/watermill/message"
    "github.com/omkar273/codegeeky/internal/logger"
    "github.com/omkar273/codegeeky/internal/service"
    "github.com/omkar273/codegeeky/internal/types"
)

type RazorpayEventHandler struct {
    paymentService service.PaymentService
    userService    service.UserService
    logger         *logger.Logger
}

func NewRazorpayEventHandler(
    paymentService service.PaymentService,
    userService service.UserService,
    logger *logger.Logger,
) *RazorpayEventHandler {
    return &RazorpayEventHandler{
        paymentService: paymentService,
        userService:    userService,
        logger:         logger,
    }
}

func (h *RazorpayEventHandler) HandleEvent(msg *message.Message) error {
    ctx := msg.Context()

    var event types.WebhookEvent
    if err := json.Unmarshal(msg.Payload, &event); err != nil {
        return err
    }

    // Only process Razorpay events
    if event.Source != types.EventSourceRazorpay {
        return nil
    }

    switch event.EventName {
    case "razorpay.payment.captured":
        return h.handlePaymentCaptured(ctx, event)
    case "razorpay.payment.failed":
        return h.handlePaymentFailed(ctx, event)
    case "razorpay.order.paid":
        return h.handleOrderPaid(ctx, event)
    case "razorpay.refund.processed":
        return h.handleRefundProcessed(ctx, event)
    default:
        h.logger.Debugw("unhandled razorpay event", "event_name", event.EventName)
        return nil
    }
}

func (h *RazorpayEventHandler) handlePaymentCaptured(ctx context.Context, event types.WebhookEvent) error {
    var razorpayEvent types.RazorpayPaymentEvent
    payloadBytes, _ := json.Marshal(event.Payload)
    if err := json.Unmarshal(payloadBytes, &razorpayEvent); err != nil {
        return err
    }

    // Update payment status in database
    updateReq := &service.UpdatePaymentStatusRequest{
        PaymentID:     razorpayEvent.Payment.ID,
        Status:        "completed",
        GatewayID:     razorpayEvent.Payment.ID,
        Amount:        float64(razorpayEvent.Payment.Amount) / 100, // Convert from paise
        Currency:      razorpayEvent.Payment.Currency,
        ProcessedAt:   event.CreatedAt,
    }

    if err := h.paymentService.UpdatePaymentStatus(ctx, updateReq); err != nil {
        return err
    }

    h.logger.Infow("payment captured successfully",
        "payment_id", razorpayEvent.Payment.ID,
        "amount", updateReq.Amount,
        "user_id", event.UserID,
    )

    return nil
}

func (h *RazorpayEventHandler) handlePaymentFailed(ctx context.Context, event types.WebhookEvent) error {
    var razorpayEvent types.RazorpayPaymentEvent
    payloadBytes, _ := json.Marshal(event.Payload)
    if err := json.Unmarshal(payloadBytes, &razorpayEvent); err != nil {
        return err
    }

    // Update payment status and handle failure logic
    updateReq := &service.UpdatePaymentStatusRequest{
        PaymentID:   razorpayEvent.Payment.ID,
        Status:      "failed",
        GatewayID:   razorpayEvent.Payment.ID,
        ProcessedAt: event.CreatedAt,
    }

    if err := h.paymentService.UpdatePaymentStatus(ctx, updateReq); err != nil {
        return err
    }

    // Trigger failure notifications
    if event.UserID != nil {
        notificationEvent := &types.WebhookEvent{
            ID:        generateEventID(),
            EventName: "notification.payment.failed",
            Source:    types.EventSourceInternal,
            UserID:    event.UserID,
            Payload: map[string]interface{}{
                "payment_id": razorpayEvent.Payment.ID,
                "amount":     float64(razorpayEvent.Payment.Amount) / 100,
                "reason":     "Payment processing failed",
            },
            CreatedAt: event.CreatedAt,
        }

        // Publish notification event (will be handled by notification service)
        // This would be done through your event publisher
    }

    return nil
}
```

## WebSocket Integration

### Frontend WebSocket Client

```javascript
// frontend/src/services/websocket.js
class EventWebSocket {
  constructor(token, userId) {
    this.token = token;
    this.userId = userId;
    this.ws = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.eventHandlers = new Map();
  }

  connect() {
    const wsUrl = `ws://localhost:8080/ws?token=${this.token}`;
    this.ws = new WebSocket(wsUrl);

    this.ws.onopen = () => {
      console.log("WebSocket connected");
      this.reconnectAttempts = 0;

      // Send subscription message
      this.subscribe(["payment.*", "user.*", "notification.*"]);
    };

    this.ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);
        this.handleMessage(message);
      } catch (error) {
        console.error("Failed to parse WebSocket message:", error);
      }
    };

    this.ws.onclose = () => {
      console.log("WebSocket connection closed");
      this.reconnect();
    };

    this.ws.onerror = (error) => {
      console.error("WebSocket error:", error);
    };
  }

  handleMessage(message) {
    switch (message.type) {
      case "event":
        this.handleEvent(message.payload);
        break;
      case "pong":
        // Handle ping-pong for keep-alive
        break;
      default:
        console.log("Unknown message type:", message.type);
    }
  }

  handleEvent(event) {
    const eventName = event.event_name;

    // Call registered handlers
    this.eventHandlers.forEach((handler, pattern) => {
      if (this.matchesPattern(eventName, pattern)) {
        handler(event);
      }
    });

    // Handle specific events
    switch (eventName) {
      case "razorpay.payment.captured":
        this.handlePaymentSuccess(event);
        break;
      case "razorpay.payment.failed":
        this.handlePaymentFailure(event);
        break;
      case "notification.payment.failed":
        this.showNotification(event.payload);
        break;
      default:
        console.log("Received event:", event);
    }
  }

  handlePaymentSuccess(event) {
    const payment = event.payload.payment;

    // Show success notification
    this.showNotification({
      type: "success",
      title: "Payment Successful!",
      message: `Payment of ₹${
        payment.amount / 100
      } has been processed successfully.`,
      duration: 5000,
    });

    // Update UI state
    this.updatePaymentStatus(payment.order_id, "completed");

    // Trigger page-specific handlers
    window.dispatchEvent(
      new CustomEvent("paymentSuccess", {
        detail: { payment, event },
      })
    );
  }

  handlePaymentFailure(event) {
    const payment = event.payload.payment;

    // Show error notification
    this.showNotification({
      type: "error",
      title: "Payment Failed",
      message: `Payment processing failed. Please try again.`,
      duration: 8000,
    });

    // Update UI state
    this.updatePaymentStatus(payment.order_id, "failed");

    // Trigger page-specific handlers
    window.dispatchEvent(
      new CustomEvent("paymentFailure", {
        detail: { payment, event },
      })
    );
  }

  subscribe(eventTypes) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(
        JSON.stringify({
          type: "subscribe",
          payload: {
            event_types: eventTypes,
            user_id: this.userId,
          },
        })
      );
    }
  }

  on(eventPattern, handler) {
    this.eventHandlers.set(eventPattern, handler);
  }

  off(eventPattern) {
    this.eventHandlers.delete(eventPattern);
  }

  matchesPattern(eventName, pattern) {
    // Simple pattern matching (supports wildcards)
    const regex = new RegExp(pattern.replace("*", ".*"));
    return regex.test(eventName);
  }

  showNotification(notification) {
    // Integrate with your notification system
    if (window.showToast) {
      window.showToast(notification);
    } else {
      console.log("Notification:", notification);
    }
  }

  updatePaymentStatus(orderId, status) {
    // Update UI components that show payment status
    const statusElements = document.querySelectorAll(
      `[data-order-id="${orderId}"]`
    );
    statusElements.forEach((element) => {
      element.classList.remove("pending", "completed", "failed");
      element.classList.add(status);
      element.textContent = status.charAt(0).toUpperCase() + status.slice(1);
    });
  }

  reconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      const delay = Math.pow(2, this.reconnectAttempts) * 1000; // Exponential backoff

      console.log(
        `Attempting to reconnect in ${delay}ms (attempt ${this.reconnectAttempts})`
      );

      setTimeout(() => {
        this.connect();
      }, delay);
    } else {
      console.error("Max reconnection attempts reached");
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }
}

// Usage example
const wsClient = new EventWebSocket(userToken, userId);
wsClient.connect();

// Register custom event handlers
wsClient.on("user.profile.updated", (event) => {
  console.log("Profile updated:", event);
  refreshUserProfile();
});

wsClient.on("payment.*", (event) => {
  console.log("Payment event received:", event);
  updatePaymentDashboard();
});
```

### React Integration Example

```javascript
// frontend/src/hooks/useWebSocket.js
import { useEffect, useRef, useState } from "react";

export const useWebSocket = (token, userId) => {
  const [isConnected, setIsConnected] = useState(false);
  const [lastEvent, setLastEvent] = useState(null);
  const wsRef = useRef(null);
  const eventHandlersRef = useRef(new Map());

  useEffect(() => {
    if (!token || !userId) return;

    const ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);
    wsRef.current = ws;

    ws.onopen = () => {
      setIsConnected(true);
      // Subscribe to user-specific events
      ws.send(
        JSON.stringify({
          type: "subscribe",
          payload: {
            event_types: ["payment.*", "user.*", "notification.*"],
            user_id: userId,
          },
        })
      );
    };

    ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      if (message.type === "event") {
        setLastEvent(message.payload);

        // Call registered handlers
        eventHandlersRef.current.forEach((handler, pattern) => {
          if (matchesPattern(message.payload.event_name, pattern)) {
            handler(message.payload);
          }
        });
      }
    };

    ws.onclose = () => {
      setIsConnected(false);
    };

    return () => {
      ws.close();
    };
  }, [token, userId]);

  const subscribe = (pattern, handler) => {
    eventHandlersRef.current.set(pattern, handler);
  };

  const unsubscribe = (pattern) => {
    eventHandlersRef.current.delete(pattern);
  };

  const sendMessage = (message) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message));
    }
  };

  return {
    isConnected,
    lastEvent,
    subscribe,
    unsubscribe,
    sendMessage,
  };
};

function matchesPattern(eventName, pattern) {
  const regex = new RegExp(pattern.replace("*", ".*"));
  return regex.test(eventName);
}
```

```javascript
// frontend/src/components/PaymentStatus.jsx
import React, { useEffect, useState } from "react";
import { useWebSocket } from "../hooks/useWebSocket";

const PaymentStatus = ({ orderId, token, userId }) => {
  const [paymentStatus, setPaymentStatus] = useState("pending");
  const [paymentDetails, setPaymentDetails] = useState(null);
  const { isConnected, subscribe, unsubscribe } = useWebSocket(token, userId);

  useEffect(() => {
    const handlePaymentEvent = (event) => {
      if (event.payload.payment?.order_id === orderId) {
        switch (event.event_name) {
          case "razorpay.payment.captured":
            setPaymentStatus("completed");
            setPaymentDetails(event.payload.payment);
            break;
          case "razorpay.payment.failed":
            setPaymentStatus("failed");
            setPaymentDetails(event.payload.payment);
            break;
          default:
            break;
        }
      }
    };

    subscribe("razorpay.payment.*", handlePaymentEvent);

    return () => {
      unsubscribe("razorpay.payment.*");
    };
  }, [orderId, subscribe, unsubscribe]);

  const getStatusColor = () => {
    switch (paymentStatus) {
      case "completed":
        return "text-green-600";
      case "failed":
        return "text-red-600";
      default:
        return "text-yellow-600";
    }
  };

  const getStatusIcon = () => {
    switch (paymentStatus) {
      case "completed":
        return "✅";
      case "failed":
        return "❌";
      default:
        return "⏳";
    }
  };

  return (
    <div className="payment-status-card p-4 border rounded-lg">
      <div className="flex items-center space-x-2">
        <span className="text-2xl">{getStatusIcon()}</span>
        <div>
          <h3 className="font-semibold">Payment Status</h3>
          <p className={`text-sm ${getStatusColor()}`}>
            {paymentStatus.charAt(0).toUpperCase() + paymentStatus.slice(1)}
          </p>
        </div>
      </div>

      {paymentDetails && (
        <div className="mt-4 text-sm text-gray-600">
          <p>Amount: ₹{paymentDetails.amount / 100}</p>
          <p>Payment ID: {paymentDetails.id}</p>
          <p>Order ID: {paymentDetails.order_id}</p>
        </div>
      )}

      <div className="mt-2 text-xs text-gray-500">
        WebSocket: {isConnected ? "🟢 Connected" : "🔴 Disconnected"}
      </div>
    </div>
  );
};

export default PaymentStatus;
```

## Microservice Communication

### Publishing Events from Services

```go
// internal/service/payment_service.go
package service

import (
    "context"
    "time"

    "github.com/omkar273/codegeeky/internal/types"
    "github.com/omkar273/codegeeky/internal/webhook/publisher"
)

type PaymentService struct {
    eventPublisher publisher.WebhookPublisher
    // other dependencies
}

func (s *PaymentService) CreatePayment(ctx context.Context, req *CreatePaymentRequest) (*Payment, error) {
    // Create payment logic
    payment, err := s.createPaymentInDB(ctx, req)
    if err != nil {
        return nil, err
    }

    // Publish payment created event
    event := &types.WebhookEvent{
        ID:        generateEventID(),
        EventName: "payment.created",
        Source:    types.EventSourceInternal,
        UserID:    &req.UserID,
        Payload: types.PaymentCreatedEvent{
            PaymentID:   payment.ID,
            UserID:      req.UserID,
            Amount:      req.Amount,
            Currency:    req.Currency,
            Status:      "pending",
            Gateway:     req.Gateway,
            CreatedAt:   time.Now(),
        },
        Metadata: types.EventMetadata{
            Version:       "1.0",
            CorrelationID: getCorrelationID(ctx),
            TraceID:       getTraceID(ctx),
        },
        CreatedAt: time.Now(),
    }

    if err := s.eventPublisher.PublishWebhook(ctx, event); err != nil {
        s.logger.Errorw("failed to publish payment created event", "error", err)
        // Don't fail the main operation for event publishing errors
    }

    return payment, nil
}

func (s *PaymentService) UpdatePaymentStatus(ctx context.Context, req *UpdatePaymentStatusRequest) error {
    // Update payment in database
    oldPayment, err := s.getPaymentByID(ctx, req.PaymentID)
    if err != nil {
        return err
    }

    err = s.updatePaymentInDB(ctx, req)
    if err != nil {
        return err
    }

    // Publish status change event if status actually changed
    if oldPayment.Status != req.Status {
        event := &types.WebhookEvent{
            ID:        generateEventID(),
            EventName: "payment.status.changed",
            Source:    types.EventSourceInternal,
            UserID:    &oldPayment.UserID,
            Payload: types.PaymentStatusChangedEvent{
                PaymentID:   req.PaymentID,
                UserID:      oldPayment.UserID,
                OldStatus:   oldPayment.Status,
                NewStatus:   req.Status,
                Amount:      oldPayment.Amount,
                Currency:    oldPayment.Currency,
                ChangedAt:   time.Now(),
            },
            CreatedAt: time.Now(),
        }

        if err := s.eventPublisher.PublishWebhook(ctx, event); err != nil {
            s.logger.Errorw("failed to publish payment status changed event", "error", err)
        }
    }

    return nil
}
```

### Cross-Service Event Handling

```go
// internal/service/notification_service.go
package service

import (
    "context"
    "encoding/json"

    "github.com/ThreeDotsLabs/watermill/message"
    "github.com/omkar273/codegeeky/internal/types"
)

type NotificationService struct {
    emailService EmailService
    smsService   SMSService
    pushService  PushNotificationService
    logger       *logger.Logger
}

func (s *NotificationService) HandlePaymentEvents(msg *message.Message) error {
    ctx := msg.Context()

    var event types.WebhookEvent
    if err := json.Unmarshal(msg.Payload, &event); err != nil {
        return err
    }

    switch event.EventName {
    case "payment.created":
        return s.handlePaymentCreated(ctx, event)
    case "payment.status.changed":
        return s.handlePaymentStatusChanged(ctx, event)
    case "razorpay.payment.captured":
        return s.handlePaymentCaptured(ctx, event)
    case "razorpay.payment.failed":
        return s.handlePaymentFailed(ctx, event)
    default:
        return nil // Not interested in this event
    }
}

func (s *NotificationService) handlePaymentCaptured(ctx context.Context, event types.WebhookEvent) error {
    if event.UserID == nil {
        return nil
    }

    var razorpayEvent types.RazorpayPaymentEvent
    eventBytes, _ := json.Marshal(event.Payload)
    if err := json.Unmarshal(eventBytes, &razorpayEvent); err != nil {
        return err
    }

    // Get user details
    user, err := s.userService.GetUserByID(ctx, *event.UserID)
    if err != nil {
        return err
    }

    // Send success email
    emailReq := &SendEmailRequest{
        To:       user.Email,
        Subject:  "Payment Successful",
        Template: "payment_success",
        Data: map[string]interface{}{
            "user_name":    user.Name,
            "amount":       float64(razorpayEvent.Payment.Amount) / 100,
            "currency":     razorpayEvent.Payment.Currency,
            "payment_id":   razorpayEvent.Payment.ID,
            "order_id":     razorpayEvent.Payment.OrderID,
            "processed_at": event.CreatedAt,
        },
    }

    if err := s.emailService.SendEmail(ctx, emailReq); err != nil {
        s.logger.Errorw("failed to send payment success email",
            "user_id", *event.UserID,
            "error", err,
        )
        // Don't return error as email failure shouldn't fail the event processing
    }

    // Send push notification
    pushReq := &SendPushNotificationRequest{
        UserID: *event.UserID,
        Title:  "Payment Successful",
        Body:   fmt.Sprintf("Your payment of ₹%.2f has been processed successfully", float64(razorpayEvent.Payment.Amount)/100),
        Data: map[string]string{
            "type":       "payment_success",
            "payment_id": razorpayEvent.Payment.ID,
            "order_id":   razorpayEvent.Payment.OrderID,
        },
    }

    if err := s.pushService.SendPush(ctx, pushReq); err != nil {
        s.logger.Errorw("failed to send payment success push notification",
            "user_id", *event.UserID,
            "error", err,
        )
    }

    return nil
}
```

## Event Publishing Patterns

### Transactional Outbox Pattern

```go
// internal/repository/outbox_repository.go
package repository

import (
    "context"
    "encoding/json"
    "time"

    "github.com/omkar273/codegeeky/ent"
    "github.com/omkar273/codegeeky/internal/postgres"
    "github.com/omkar273/codegeeky/internal/types"
)

type OutboxRepository struct {
    db postgres.IClient
}

type OutboxEvent struct {
    ID        string                 `json:"id"`
    EventName string                 `json:"event_name"`
    Payload   map[string]interface{} `json:"payload"`
    Status    string                 `json:"status"` // pending, published, failed
    CreatedAt time.Time              `json:"created_at"`
    UpdatedAt time.Time              `json:"updated_at"`
}

func (r *OutboxRepository) SaveEvent(ctx context.Context, event *types.WebhookEvent) error {
    return r.db.WithTx(ctx, func(ctx context.Context) error {
        client := r.db.Querier(ctx)

        payloadBytes, _ := json.Marshal(event.Payload)
        var payload map[string]interface{}
        json.Unmarshal(payloadBytes, &payload)

        _, err := client.OutboxEvent.Create().
            SetID(event.ID).
            SetEventName(event.EventName).
            SetPayload(payload).
            SetStatus("pending").
            SetCreatedAt(event.CreatedAt).
            SetUpdatedAt(time.Now()).
            Save(ctx)

        return err
    })
}

func (r *OutboxRepository) GetPendingEvents(ctx context.Context, limit int) ([]*OutboxEvent, error) {
    client := r.db.Querier(ctx)

    events, err := client.OutboxEvent.Query().
        Where(outboxevent.StatusEQ("pending")).
        Order(ent.Asc(outboxevent.FieldCreatedAt)).
        Limit(limit).
        All(ctx)

    if err != nil {
        return nil, err
    }

    result := make([]*OutboxEvent, len(events))
    for i, event := range events {
        result[i] = &OutboxEvent{
            ID:        event.ID,
            EventName: event.EventName,
            Payload:   event.Payload,
            Status:    event.Status,
            CreatedAt: event.CreatedAt,
            UpdatedAt: event.UpdatedAt,
        }
    }

    return result, nil
}

func (r *OutboxRepository) MarkAsPublished(ctx context.Context, eventID string) error {
    client := r.db.Querier(ctx)

    return client.OutboxEvent.UpdateOneID(eventID).
        SetStatus("published").
        SetUpdatedAt(time.Now()).
        Exec(ctx)
}
```

```go
// internal/service/outbox_publisher.go
package service

import (
    "context"
    "time"

    "github.com/omkar273/codegeeky/internal/repository"
    "github.com/omkar273/codegeeky/internal/types"
    "github.com/omkar273/codegeeky/internal/webhook/publisher"
)

type OutboxPublisher struct {
    outboxRepo repository.OutboxRepository
    publisher  publisher.WebhookPublisher
    logger     *logger.Logger
}

func NewOutboxPublisher(
    outboxRepo repository.OutboxRepository,
    publisher publisher.WebhookPublisher,
    logger *logger.Logger,
) *OutboxPublisher {
    return &OutboxPublisher{
        outboxRepo: outboxRepo,
        publisher:  publisher,
        logger:     logger,
    }
}

func (p *OutboxPublisher) Start(ctx context.Context) {
    ticker := time.NewTicker(5 * time.Second) // Poll every 5 seconds
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            p.processOutboxEvents(ctx)
        }
    }
}

func (p *OutboxPublisher) processOutboxEvents(ctx context.Context) {
    events, err := p.outboxRepo.GetPendingEvents(ctx, 100)
    if err != nil {
        p.logger.Errorw("failed to get pending outbox events", "error", err)
        return
    }

    for _, outboxEvent := range events {
        event := &types.WebhookEvent{
            ID:        outboxEvent.ID,
            EventName: outboxEvent.EventName,
            Payload:   outboxEvent.Payload,
            CreatedAt: outboxEvent.CreatedAt,
        }

        if err := p.publisher.PublishWebhook(ctx, event); err != nil {
            p.logger.Errorw("failed to publish outbox event",
                "event_id", outboxEvent.ID,
                "error", err,
            )
            continue
        }

        if err := p.outboxRepo.MarkAsPublished(ctx, outboxEvent.ID); err != nil {
            p.logger.Errorw("failed to mark outbox event as published",
                "event_id", outboxEvent.ID,
                "error", err,
            )
        }

        p.logger.Debugw("published outbox event",
            "event_id", outboxEvent.ID,
            "event_name", outboxEvent.EventName,
        )
    }
}
```

## Error Handling Examples

### Circuit Breaker Implementation

```go
// internal/circuitbreaker/circuit_breaker.go
package circuitbreaker

import (
    "context"
    "errors"
    "sync"
    "time"
)

type State int

const (
    StateClosed State = iota
    StateHalfOpen
    StateOpen
)

type CircuitBreaker struct {
    name           string
    maxFailures    int
    resetTimeout   time.Duration
    halfOpenCount  int

    mu             sync.RWMutex
    state          State
    failures       int
    lastFailTime   time.Time
    halfOpenTries  int
}

func NewCircuitBreaker(name string, maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        name:          name,
        maxFailures:   maxFailures,
        resetTimeout:  resetTimeout,
        halfOpenCount: 3,
        state:         StateClosed,
    }
}

func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
    if !cb.canExecute() {
        return errors.New("circuit breaker is open")
    }

    err := fn()
    cb.recordResult(err)
    return err
}

func (cb *CircuitBreaker) canExecute() bool {
    cb.mu.RLock()
    defer cb.mu.RUnlock()

    switch cb.state {
    case StateClosed:
        return true
    case StateOpen:
        return time.Since(cb.lastFailTime) > cb.resetTimeout
    case StateHalfOpen:
        return cb.halfOpenTries < cb.halfOpenCount
    default:
        return false
    }
}

func (cb *CircuitBreaker) recordResult(err error) {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    if err != nil {
        cb.failures++
        cb.lastFailTime = time.Now()

        if cb.state == StateHalfOpen {
            cb.state = StateOpen
            cb.halfOpenTries = 0
        } else if cb.failures >= cb.maxFailures {
            cb.state = StateOpen
        }
    } else {
        cb.failures = 0

        if cb.state == StateHalfOpen {
            cb.halfOpenTries++
            if cb.halfOpenTries >= cb.halfOpenCount {
                cb.state = StateClosed
                cb.halfOpenTries = 0
            }
        } else if cb.state == StateOpen && time.Since(cb.lastFailTime) > cb.resetTimeout {
            cb.state = StateHalfOpen
            cb.halfOpenTries = 0
        }
    }
}

func (cb *CircuitBreaker) GetState() State {
    cb.mu.RLock()
    defer cb.mu.RUnlock()
    return cb.state
}
```

### Webhook Delivery with Circuit Breaker

```go
// internal/webhook/delivery/delivery.go
package delivery

import (
    "context"
    "fmt"

    "github.com/omkar273/codegeeky/internal/circuitbreaker"
    "github.com/omkar273/codegeeky/internal/httpclient"
)

type WebhookDelivery struct {
    client          httpclient.Client
    circuitBreakers map[string]*circuitbreaker.CircuitBreaker
    logger          *logger.Logger
}

func NewWebhookDelivery(client httpclient.Client, logger *logger.Logger) *WebhookDelivery {
    return &WebhookDelivery{
        client:          client,
        circuitBreakers: make(map[string]*circuitbreaker.CircuitBreaker),
        logger:          logger,
    }
}

func (w *WebhookDelivery) DeliverWebhook(ctx context.Context, endpoint string, payload []byte, headers map[string]string) error {
    cb := w.getCircuitBreaker(endpoint)

    return cb.Execute(ctx, func() error {
        req := &httpclient.Request{
            Method:  "POST",
            URL:     endpoint,
            Headers: headers,
            Body:    payload,
        }

        resp, err := w.client.Send(ctx, req)
        if err != nil {
            return fmt.Errorf("webhook delivery failed: %w", err)
        }

        if resp.StatusCode >= 400 {
            return fmt.Errorf("webhook returned status %d", resp.StatusCode)
        }

        return nil
    })
}

func (w *WebhookDelivery) getCircuitBreaker(endpoint string) *circuitbreaker.CircuitBreaker {
    if cb, exists := w.circuitBreakers[endpoint]; exists {
        return cb
    }

    cb := circuitbreaker.NewCircuitBreaker(
        fmt.Sprintf("webhook_%s", endpoint),
        5,  // Max failures
        30*time.Second, // Reset timeout
    )

    w.circuitBreakers[endpoint] = cb
    return cb
}
```

This comprehensive API integration guide provides practical examples for implementing all the key components of your event-driven system, from Razorpay webhook handling to real-time frontend updates.
