# System Architecture Documentation

## Overview

This document provides a comprehensive overview of the event-driven pubsub system architecture, showing how different components interact and the flow of data through the system.

## Current vs Proposed Architecture

### Current Architecture

```
┌─────────────────┐
│  Application    │
│     Events      │
└─────────┬───────┘
          │
          ▼
┌─────────────────┐    ┌─────────────────┐
│   Webhook       │───▶│   Memory        │
│   Publisher     │    │   PubSub        │
└─────────────────┘    └─────────┬───────┘
                                 │
                                 ▼
                    ┌─────────────────┐
                    │   Webhook       │
                    │   Handler       │
                    └─────────┬───────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │   HTTP Client   │
                    │   (External     │
                    │   Webhooks)     │
                    └─────────────────┘
```

**Limitations:**

- In-memory only (no persistence)
- Single node deployment
- No external event ingestion
- No real-time frontend connectivity

### Proposed Architecture

```
                    ┌─────────────────────────────────────────────────────────────┐
                    │                    EXTERNAL SYSTEMS                         │
                    │                                                             │
    ┌───────────────┼─────────────┐    ┌─────────────┐    ┌─────────────────────┼─┐
    │               │             │    │             │    │                     │ │
    │   Razorpay    │  Other      │    │ Micro-      │    │    Frontend         │ │
    │   Webhooks    │  External   │    │ services    │    │   Applications      │ │
    │               │  APIs       │    │             │    │                     │ │
    └───────────────┼─────────────┘    └─────────────┘    └─────────────────────┼─┘
                    │                                                             │
                    └─────────────────────┬───────────────────────────────────────┘
                                          │
                    ┌─────────────────────▼───────────────────────────────────────┐
                    │                 API GATEWAY LAYER                           │
                    │                                                             │
    ┌───────────────┼─────────────┐    ┌─────────────┐    ┌─────────────────────┼─┐
    │               │             │    │             │    │                     │ │
    │   Webhook     │  REST API   │    │ GraphQL     │    │    WebSocket        │ │
    │   Receivers   │  Endpoints  │    │ Endpoints   │    │    Gateway          │ │
    │               │             │    │             │    │                     │ │
    └───────────────┼─────────────┘    └─────────────┘    └─────────────────────┼─┘
                    │                                                             │
                    └─────────────────────┬───────────────────────────────────────┘
                                          │
                    ┌─────────────────────▼───────────────────────────────────────┐
                    │              EVENT PROCESSING CORE                          │
                    │                                                             │
    ┌───────────────┼─────────────┐    ┌─────────────┐    ┌─────────────────────┼─┐
    │               │             │    │             │    │                     │ │
    │   Event       │   Event     │    │   Event     │    │    Event Schema     │ │
    │   Router      │ Validator   │    │Transformer  │    │    Registry         │ │
    │               │             │    │             │    │                     │ │
    └───────────────┼─────────────┘    └─────────────┘    └─────────────────────┼─┘
                    │                                                             │
                    └─────────────────────┬───────────────────────────────────────┘
                                          │
                    ┌─────────────────────▼───────────────────────────────────────┐
                    │              PUBSUB INFRASTRUCTURE                          │
                    │                                                             │
    ┌───────────────┼─────────────┐    ┌─────────────┐    ┌─────────────────────┼─┐
    │               │             │    │             │    │                     │ │
    │   Kafka       │   Redis     │    │   NATS      │    │     Memory          │ │
    │ (Production)  │(Mid-scale)  │    │(Cloud)      │    │  (Development)      │ │
    │               │             │    │             │    │                     │ │
    └───────────────┼─────────────┘    └─────────────┘    └─────────────────────┼─┘
                    │                                                             │
                    └─────────────────────┬───────────────────────────────────────┘
                                          │
                    ┌─────────────────────▼───────────────────────────────────────┐
                    │               EVENT HANDLERS                                │
                    │                                                             │
    ┌───────────────┼─────────────┐    ┌─────────────┐    ┌─────────────────────┼─┐
    │               │             │    │             │    │                     │ │
    │   Webhook     │ Microservice│    │  Real-time  │    │     Audit           │ │
    │   Handler     │   Handler   │    │   Handler   │    │     Handler         │ │
    │               │             │    │             │    │                     │ │
    └───────────────┼─────────────┘    └─────────────┘    └─────────────────────┼─┘
                    │                                                             │
                    └─────────────────────┬───────────────────────────────────────┘
                                          │
                    ┌─────────────────────▼───────────────────────────────────────┐
                    │                OUTPUT CHANNELS                              │
                    │                                                             │
    ┌───────────────┼─────────────┐    ┌─────────────┐    ┌─────────────────────┼─┐
    │               │             │    │             │    │                     │ │
    │   HTTP        │ WebSocket   │    │ Server-Sent │    │    Event Store      │ │
    │   Webhooks    │Connections  │    │   Events    │    │   (Database)        │ │
    │               │             │    │             │    │                     │ │
    └───────────────┼─────────────┘    └─────────────┘    └─────────────────────┼─┘
                    │                                                             │
                    └─────────────────────────────────────────────────────────────┘
```

## Component Interactions

### Event Flow Sequence

```
┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│   External  │  │  Webhook    │  │   Event     │  │   PubSub    │  │   Event     │
│   Service   │  │  Receiver   │  │ Processor   │  │   System    │  │  Handlers   │
│ (Razorpay)  │  │             │  │             │  │             │  │             │
└──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘
       │                │                │                │                │
       │                │                │                │                │
   1.  │─── Webhook ────▶│                │                │                │
       │                │                │                │                │
   2.  │                │─ Validate ─────▶│                │                │
       │                │  Signature      │                │                │
       │                │                │                │                │
   3.  │                │                │─ Transform ────▶│                │
       │                │                │  to Internal   │                │
       │                │                │  Event Format  │                │
       │                │                │                │                │
   4.  │                │                │                │─ Route Event ─▶│
       │                │                │                │  to Handlers   │
       │                │                │                │                │
   5.  │◀── Response ───│◀─── Success ───│◀─── Success ───│◀─── Success ───│
       │    (200 OK)    │      (ACK)     │      (ACK)     │      (ACK)     │
       │                │                │                │                │
```

### Real-time Data Flow

```
┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│   Payment   │  │   Event     │  │  Real-time  │  │  WebSocket  │  │   Frontend  │
│   Service   │  │   Bus       │  │   Handler   │  │   Gateway   │  │ Application │
│             │  │             │  │             │  │             │  │             │
└──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘
       │                │                │                │                │
       │                │                │                │                │
   1.  │─── Publish ────▶│                │                │                │
       │    Event        │                │                │                │
       │                │                │                │                │
   2.  │                │─── Route ──────▶│                │                │
       │                │    to RT        │                │                │
       │                │    Handler      │                │                │
       │                │                │                │                │
   3.  │                │                │─── Broadcast ──▶│                │
       │                │                │    to User      │                │
       │                │                │    Connection   │                │
       │                │                │                │                │
   4.  │                │                │                │─── Send Event ─▶│
       │                │                │                │    via WS       │
       │                │                │                │                │
   5.  │                │                │                │                │─ Update UI
       │                │                │                │                │  Show Notification
       │                │                │                │                │
```

## Deployment Architecture

### Development Environment

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            DEVELOPMENT ENVIRONMENT                              │
│                                                                                 │
│   ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────────────────┐ │
│   │                 │    │                 │    │                             │ │
│   │   Application   │    │   Memory        │    │        PostgreSQL          │ │
│   │    Server       │────│   PubSub        │    │        Database             │ │
│   │  (Port: 8080)   │    │                 │    │      (Port: 5432)           │ │
│   │                 │    │                 │    │                             │ │
│   └─────────────────┘    └─────────────────┘    └─────────────────────────────┘ │
│                                                                                 │
│   ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────────────────┐ │
│   │                 │    │                 │    │                             │ │
│   │    Frontend     │    │     Redis       │    │         Ngrok               │ │
│   │   Development   │    │     Cache       │    │    (Webhook Testing)        │ │
│   │  (Port: 3000)   │    │  (Port: 6379)   │    │                             │ │
│   │                 │    │                 │    │                             │ │
│   └─────────────────┘    └─────────────────┘    └─────────────────────────────┘ │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### Production Environment

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                             PRODUCTION ENVIRONMENT                              │
│                                                                                 │
│ ┌─────────────────────────────────────────────────────────────────────────────┐ │
│ │                              LOAD BALANCER                                   │ │
│ │                           (NGINX/Cloudflare)                                 │ │
│ └─────────────────────────┬───────────────────────────────────────────────────┘ │
│                           │                                                     │
│ ┌─────────────────────────▼───────────────────────────────────────────────────┐ │
│ │                         APPLICATION TIER                                    │ │
│ │                                                                             │ │
│ │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │ │
│ │  │   App       │  │   App       │  │   App       │  │      WebSocket      │ │ │
│ │  │ Instance 1  │  │ Instance 2  │  │ Instance 3  │  │      Gateway        │ │ │
│ │  │             │  │             │  │             │  │                     │ │ │
│ │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────────────┘ │ │
│ └─────────────────────────┬───────────────────────────────────────────────────┘ │
│                           │                                                     │
│ ┌─────────────────────────▼───────────────────────────────────────────────────┐ │
│ │                        MESSAGE BROKER TIER                                  │ │
│ │                                                                             │ │
│ │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │ │
│ │  │   Kafka     │  │   Kafka     │  │   Kafka     │  │        Redis        │ │ │
│ │  │  Broker 1   │  │  Broker 2   │  │  Broker 3   │  │       Cluster       │ │ │
│ │  │             │  │             │  │             │  │                     │ │ │
│ │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────────────┘ │ │
│ └─────────────────────────┬───────────────────────────────────────────────────┘ │
│                           │                                                     │
│ ┌─────────────────────────▼───────────────────────────────────────────────────┐ │
│ │                          DATABASE TIER                                      │ │
│ │                                                                             │ │
│ │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │ │
│ │  │ PostgreSQL  │  │ PostgreSQL  │  │ PostgreSQL  │  │     Monitoring      │ │ │
│ │  │   Master    │  │  Replica 1  │  │  Replica 2  │  │    (Prometheus)     │ │ │
│ │  │             │  │             │  │             │  │                     │ │ │
│ │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────────────┘ │ │
│ └─────────────────────────────────────────────────────────────────────────────┘ │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
```

## Data Flow Patterns

### Event Publishing Pattern

```
Application Service
        │
        ▼
┌───────────────┐
│   Business    │
│    Logic      │
│  Execution    │
└───────┬───────┘
        │
        ▼
┌───────────────┐     ┌───────────────┐
│   Database    │────▶│   Outbox      │
│  Transaction  │     │    Table      │
└───────────────┘     └───────┬───────┘
                              │
                              ▼
                      ┌───────────────┐
                      │   Outbox      │
                      │  Publisher    │
                      │  (Background) │
                      └───────┬───────┘
                              │
                              ▼
                      ┌───────────────┐
                      │   Message     │
                      │    Broker     │
                      └───────────────┘
```

### Event Processing Pattern

```
Message Broker
        │
        ▼
┌───────────────┐
│   Consumer    │
│   Group       │
└───────┬───────┘
        │
        ▼
┌───────────────┐     ┌───────────────┐
│   Message     │────▶│   Dead Letter │
│  Processing   │     │     Queue     │
│               │     │  (On Failure) │
└───────┬───────┘     └───────────────┘
        │
        ▼
┌───────────────┐
│   Business    │
│    Logic      │
│  Processing   │
└───────┬───────┘
        │
        ▼
┌───────────────┐
│   Acknowledge │
│    Message    │
└───────────────┘
```

## Scalability Considerations

### Horizontal Scaling Strategy

1. **Stateless Application Design**

   - All application instances are identical
   - No shared state between instances
   - Session data stored in external systems

2. **Message Partitioning**

   - Events partitioned by user ID or tenant ID
   - Ensures ordered processing per partition
   - Enables parallel processing across partitions

3. **Consumer Groups**

   - Multiple consumer instances in same group
   - Automatic load balancing
   - Fault tolerance through rebalancing

4. **Database Scaling**
   - Read replicas for query scaling
   - Connection pooling
   - Database sharding for extreme scale

### Performance Optimization

1. **Batch Processing**

   - Process multiple events in single operation
   - Reduce network overhead
   - Improve throughput

2. **Connection Pooling**

   - Reuse database connections
   - Reduce connection establishment overhead
   - Configure appropriate pool sizes

3. **Caching Strategy**

   - Cache frequently accessed data
   - Use Redis for session storage
   - Implement cache invalidation patterns

4. **Async Processing**
   - Non-blocking I/O operations
   - Background job processing
   - Event-driven architecture benefits

## Security Architecture

### Authentication & Authorization Flow

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Client    │    │   API       │    │    Auth     │    │   Event     │
│ Application │    │  Gateway    │    │  Service    │    │   System    │
└──────┬──────┘    └──────┬──────┘    └──────┬──────┘    └──────┬──────┘
       │                  │                  │                  │
   1.  │─── Request ──────▶│                  │                  │
       │    + JWT Token    │                  │                  │
       │                  │                  │                  │
   2.  │                  │─── Validate ────▶│                  │
       │                  │    Token         │                  │
       │                  │                  │                  │
   3.  │                  │◀─── User Info ───│                  │
       │                  │                  │                  │
   4.  │                  │─── Process ──────────────────────────▶│
       │                  │    Event         │                  │
       │                  │    + User Context│                  │
       │                  │                  │                  │
   5.  │◀─── Response ────│◀─── Success ─────────────────────────│
       │                  │                  │                  │
```

### Webhook Security

1. **Signature Validation**

   - HMAC-SHA256 signature verification
   - Prevent request tampering
   - Ensure payload integrity

2. **IP Allowlisting**

   - Restrict webhook sources
   - Block unauthorized requests
   - Firewall-level protection

3. **Rate Limiting**

   - Prevent abuse and DoS attacks
   - Per-IP and per-endpoint limits
   - Configurable thresholds

4. **Request Size Limits**
   - Prevent memory exhaustion
   - Limit payload sizes
   - Early request rejection

This architecture provides a robust, scalable foundation for your event-driven system while maintaining security and performance standards.
