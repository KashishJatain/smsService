# Polyglot Distributed SMS Service

A two-microservice system for sending and storing SMS notifications.

| Service | Language | Port | Role |
|---|---|---|---|
| **sms-sender** | Java 17 / Spring Boot 3 | `8080` | Gateway – validates requests, checks Redis block list, calls mock vendor, publishes to Kafka |
| **sms-store** | Go 1.22 / net/http | `8081` | Persistence – consumes Kafka events, stores to MongoDB, serves history API |

---

## Architecture

```
Client
  │
  ▼
SMS Sender (Java :8080)
  ├── Redis  ← block-list check
  ├── Mock 3P Vendor  ← simulated send
  └── Kafka (sms-events topic) ──► SMS Store (Go :8081)
                                        ├── MongoDB  ← persist record
                                        └── GET /v1/user/:id/messages
```

---

## Prerequisites

- Docker ≥ 24 and Docker Compose v2
- Java 17+, Maven 3.9+, Go 1.22+

---

## Quick Start

```bash
# 1. Clone / unzip the project
cd sms-service

# 2. Start everything
docker compose up --build

# 3. Run the end-to-end demo (in a separate terminal once services are healthy)
bash demo.sh
```

---

## API Reference

### SMS Sender  (`http://localhost:8080`)

#### `POST /v1/sms/send`
Send an SMS message.

**Request body**
```json
{
  "userId":      "user-123",
  "phoneNumber": "+919876543210",
  "message":     "Hello!"
}
```

**Response – 200 SUCCESS**
```json
{
  "messageId":   "3f2504e0-...",
  "status":      "SUCCESS",
  "message":     "SMS sent and queued for storage successfully.",
  "phoneNumber": "+919876543210"
}
```

**Response – 403 BLOCKED**
```json
{
  "status":      "BLOCKED",
  "message":     "User is blocked from receiving SMS.",
  "phoneNumber": "+919876543210"
}
```

**Response – 502 FAILED**
```json
{
  "status":  "FAILED",
  "message": "Failed to send SMS: VENDOR_TIMEOUT: simulated failure"
}
```

---

#### `POST /v1/sms/block/{userId}`
Add a user to the Redis block list.

**Response – 200**
```
userId=user-123 has been blocked.
```

---

#### `DELETE /v1/sms/block/{userId}`
Remove a user from the block list.

**Response – 200**
```
userId=user-123 has been unblocked.
```

---

### SMS Store  (`http://localhost:8081`)

#### `GET /v1/user/{userId}/messages`
Fetch all SMS records for a user (newest first).

**Response – 200**
```json
{
  "userId":   "user-123",
  "count":    2,
  "messages": [
    {
      "id":             "665f...",
      "messageId":      "3f2504e0-...",
      "userId":         "user-123",
      "phoneNumber":    "+919876543210",
      "message":        "Hello!",
      "status":         "SUCCESS",
      "vendorResponse": "MOCK-REF-001",
      "timestamp":      "2024-06-01T10:00:00Z",
      "createdAt":      "2024-06-01T10:00:01Z"
    }
  ]
}
```

#### `GET /health`
Health probe.

```json
{ "status": "ok", "service": "sms-store" }
```

---

## Project Structure

```
sms-service/
├── docker-compose.yml
├── demo.sh
├── README.md
│
├── sms-sender/                         # Java Spring Boot
│   ├── pom.xml
│   ├── Dockerfile
│   └── src/
│       ├── main/java/com/sms/sender/
│       │   ├── SmsSenderApplication.java
│       │   ├── controller/
│       │   │   ├── SmsController.java
│       │   │   └── GlobalExceptionHandler.java
│       │   ├── service/
│       │   │   ├── SmsService.java
│       │   │   ├── BlockListService.java
│       │   │   └── VendorService.java
│       │   ├── kafka/
│       │   │   └── SmsEventProducer.java
│       │   ├── model/
│       │   │   ├── SmsRequest.java
│       │   │   ├── SmsResponse.java
│       │   │   └── SmsEvent.java
│       │   └── config/
│       │       ├── KafkaConfig.java
│       │       └── RedisConfig.java
│       └── test/java/com/sms/sender/
│           └── service/
│               ├── SmsServiceTest.java
│               └── BlockListServiceTest.java
│
└── sms-store/                          # Go net/http
    ├── go.mod
    ├── Dockerfile
    ├── cmd/server/main.go
    ├── config/
    │   ├── config.go
    │   └── kafka_consumer.go
    └── internal/
        ├── model/sms.go
        ├── repository/sms_repository.go
        ├── service/
        │   ├── sms_service.go
        │   └── sms_service_test.go
        └── handler/
            ├── sms_handler.go
            └── sms_handler_test.go
```

---

## System Flow

### Send SMS Flow

```text
You → POST /v1/sms/send
       │
       ▼
SmsController.java
(Receives HTTP request)
       │
       ▼
SmsService.java
(Coordinates the process)
       │
       ├──► BlockListService → Redis
       │      Checks if user is blocked
       │
       ├──► VendorService
       │      Sends SMS (mock implementation)
       │
       └──► SmsEventProducer → Kafka
              Publishes SMS event
                              │
                              ▼
                     kafka_consumer.go
                     Consumes Kafka event
                              │
                              ▼
                     sms_service.go
                     Converts event to DB record
                              │
                              ▼
                     sms_repository.go → MongoDB
                     Stores SMS record

You → GET /v1/user/user-123/messages
       │
       ▼
sms_handler.go
       │
       ▼
sms_service.go
       │
       ▼
sms_repository.go → MongoDB
       │
       ▼
Returns JSON list of all SMS records                     

---

## Environment Variables

### sms-sender
| Variable | Default | Description |
|---|---|---|
| `REDIS_HOST` | `localhost` | Redis hostname |
| `REDIS_PORT` | `6379` | Redis port |
| `KAFKA_BOOTSTRAP_SERVERS` | `localhost:9092` | Kafka brokers |
| `vendor.mock-failure-rate` | `0.1` | Fraction of vendor calls that simulate failure |

### sms-store
| Variable | Default | Description |
|---|---|---|
| `SERVER_PORT` | `8081` | HTTP port |
| `MONGO_URI` | `mongodb://localhost:27017` | MongoDB connection string |
| `MONGO_DB_NAME` | `sms_store` | MongoDB database name |
| `KAFKA_BOOTSTRAP_SERVERS` | `localhost:9092` | Kafka brokers |
| `KAFKA_TOPIC` | `sms-events` | Topic to consume |
| `KAFKA_GROUP_ID` | `sms-store-consumer-group` | Consumer group |

---

## Design Decisions

- **Fail open on Redis**: if Redis is unreachable the block-list check is skipped (logged as an error) so SMS delivery continues rather than silently dropping messages.
- **Audit all events**: even BLOCKED and FAILED sends are published to Kafka so the SMS Store builds a complete audit trail.
- **Partition by userId**: Kafka messages are keyed by `userId`, so all events for a user land on the same partition and are consumed in order.
- **Idempotent storage**: MongoDB has a unique index on `message_id` so duplicate Kafka deliveries are safely rejected.
- **Graceful shutdown**: both services listen for SIGTERM, drain in-flight work, and close connections cleanly.
