# Webhook and PubSub System PRD

## Overview

This document outlines the design and implementation of a robust webhook and pubsub system that handles both incoming webhooks from payment gateways and outgoing webhooks to clients. The system uses Watermill for message routing and processing.

## System Architecture

### 1. Core Components

#### 1.1 PubSub Router

- Uses Watermill for message routing
- Implements retry mechanisms with configurable parameters
- Handles dead letter queues (DLQ) for failed messages
- Provides correlation IDs for message tracking
- Implements panic recovery

#### 1.2 Webhook Service

- Orchestrates webhook operations
- Manages webhook publishers and handlers
- Handles service lifecycle (start/stop)
- Configurable through application settings

#### 1.3 Webhook Handler

- Processes incoming webhook messages
- Validates webhook payloads
- Routes messages to appropriate handlers
- Implements retry logic for failed deliveries

### 2. Message Flow

#### 2.1 Incoming Webhooks (Payment Gateways)

1. Webhook endpoint receives request
2. Validates webhook signature
3. Processes webhook payload
4. Stores data in database
5. Publishes event to internal pubsub system
6. Sends acknowledgment response

#### 2.2 Outgoing Webhooks

1. Internal event triggers webhook
2. Message published to pubsub system
3. Webhook handler processes message
4. Validates user configuration
5. Builds webhook payload
6. Sends webhook to client endpoint
7. Handles retries and failures

### 3. Payment Gateway Integration

#### 3.1 Supported Gateways

- Razorpay
- Stripe
- (Other payment gateways as needed)

#### 3.2 Webhook Processing

1. Signature Verification

   - Validate webhook signatures
   - Prevent unauthorized requests
   - Implement gateway-specific verification

2. Payload Processing

   - Parse gateway-specific payloads
   - Map to internal event types
   - Handle different event types (payment.success, payment.failed, etc.)

3. Data Storage
   - Store webhook events
   - Update payment status
   - Maintain audit trail

### 4. Database Schema

#### 4.1 Webhook Events

```sql
CREATE TABLE webhook_events (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    event_name VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    retry_count INT DEFAULT 0,
    last_error TEXT,
    metadata JSONB
);
```

#### 4.2 Webhook Deliveries

```sql
CREATE TABLE webhook_deliveries (
    id UUID PRIMARY KEY,
    event_id UUID NOT NULL,
    endpoint VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    response_code INT,
    response_body TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (event_id) REFERENCES webhook_events(id)
);
```

### 5. Configuration

#### 5.1 Webhook Settings

```yaml
webhook:
  enabled: true
  pubsub: "memory" # or "kafka"
  max_retries: 3
  initial_interval: "1s"
  max_interval: "1m"
  multiplier: 2
  max_elapsed_time: "5m"
  users:
    - id: "user_id"
      enabled: true
      endpoint: "https://client-webhook-url"
      headers:
        Authorization: "Bearer token"
      excluded_events:
        - "event.name"
```

### 6. Error Handling

#### 6.1 Retry Strategy

- Exponential backoff
- Maximum retry attempts
- Dead letter queue for failed messages
- Error logging and monitoring

#### 6.2 Failure Scenarios

- Invalid webhook signatures
- Network failures
- Client endpoint unavailability
- Payload validation errors
- Database errors

### 7. Monitoring and Logging

#### 7.1 Metrics

- Webhook delivery success rate
- Average delivery time
- Retry counts
- Error rates by type
- Queue sizes

#### 7.2 Logging

- Webhook receipt and processing
- Delivery attempts and results
- Error details
- Performance metrics

### 8. Security

#### 8.1 Authentication

- Webhook signature verification
- API key validation
- IP whitelisting

#### 8.2 Data Protection

- Payload encryption
- Secure storage
- Access control

### 9. Implementation Guidelines

#### 9.1 Code Organization

```
internal/
  webhook/
    handler/     # Webhook handlers
    publisher/   # Webhook publishers
    payload/     # Payload builders
    dto/         # Data transfer objects
  pubsub/
    router/      # Message router
    memory/      # In-memory pubsub
    kafka/       # Kafka implementation
```

#### 9.2 Best Practices

- Use dependency injection
- Implement proper error handling
- Follow logging standards
- Write comprehensive tests
- Document API endpoints

### 10. Testing Strategy

#### 10.1 Unit Tests

- Handler logic
- Payload building
- Signature verification
- Database operations

#### 10.2 Integration Tests

- End-to-end webhook flow
- PubSub message routing
- Database interactions
- External service calls

#### 10.3 Load Tests

- Concurrent webhook processing
- Message queue performance
- Database performance

### 11. Deployment

#### 11.1 Requirements

- Go 1.21+
- PostgreSQL 13+
- Kafka (optional)
- Redis (optional)

#### 11.2 Configuration

- Environment variables
- Configuration files
- Secret management

### 12. Future Enhancements

#### 12.1 Planned Features

- Webhook management UI
- Real-time monitoring dashboard
- Advanced retry strategies
- Rate limiting
- Webhook replay functionality

#### 12.2 Scalability

- Horizontal scaling
- Load balancing
- Message partitioning
- Caching strategies
