# RabbitMQ Go Messaging Project

A demonstration of message queuing with RabbitMQ using Go.
- Direct exchange type for message routing
- Topic exchange type for message filtering

## Prerequisites

- Docker and Docker Compose
- Go 1.16+
- RabbitMQ Go client: `github.com/rabbitmq/amqp091-go`

## Project Setup

1. Initialize the Go project:
```bash
mkdir rabbitmq-go
cd rabbitmq-go
go mod init rabbitmq-go
go get github.com/rabbitmq/amqp091-go
```

2. Create `docker-compose.yml`:
```yaml
version: '3.8'

services:
  rabbitmq:
    image: rabbitmq:4.0-management
    container_name: rabbitmq
    ports:
      - "5672:5672"     # RabbitMQ messaging port
      - "15672:15672"   # Management interface
    environment:
      RABBITMQ_DEFAULT_USER: guest    # Default username
      RABBITMQ_DEFAULT_PASS: guest    # Default password
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq

volumes:
  rabbitmq_data:

```

3. Start RabbitMQ:
```bash
docker-compose up -d
```

## Project Structure

```
rabbitmq-go/
│
├── docker-compose.yml
├── go.mod
├── go.sum
├── direct/
│   ├── send/
│   │   └── send.go
│   └── receive/
│       └── receive.go
├── topic/
│   ├── send/
│   │   └── send.go
│   └── receive/
│       └── receive.go
└── README.md
```

## Running the Application

### Start RabbitMQ:
```bash
docker-compose up -d
```

### Direct Exchange

#### Build and run the consumer:
```bash
cd direct/receive
go run receive.go
```

#### In another terminal, run the producer:
```bash
cd direct/send
go run send.go 
```

### Topic Exchange

#### Build and run the consumer:
```bash
cd topic/receive
# To create Q1 and bind it to "*.data.*"
go run receive.go Q1 "*.data.*"
```

#### In another terminal, run the consumer:
```bash
cd topic/receive
# To create Q2 and bind it to "*.*.temperature" and "sensor.#"
go run receive.go Q2 "*.*.temperature" "sensor.#"
```

#### In another terminal, run the producer:
```bash
cd topic/send
go run send.go "*.data.*" 1 100 25 75 1200 7.2
go run send.go "sensor.#" 2 200 30 80 1500 8.5
go run send.go "sensor.temperature" 3 300 35 85 1800 9.8
go run send.go "sensor.humidity" 4 400 40 90 2100 10.5 
go run send.go "error.data.tds" 5 500 45 95 2400 11.2
go run send.go "sensor.room1" 6 600 50 100 2700 12.5
```


## Monitoring

Access the RabbitMQ Management Interface:
- URL: http://localhost:15672
- Username: guest
- Password: guest

Here you can monitor:
- Message rates
- Queue status
- Exchange bindings
- Consumer status
- Connection details

## Error Handling

The application includes error handling for:
- Connection failures
- Channel operations
- Message publishing
- Message consumption
- JSON marshaling/unmarshaling
