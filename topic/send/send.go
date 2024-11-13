package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SensorData struct {
	ID          int     `json:"id"`
	PacketNo    int     `json:"packet_no"`
	Temperature int     `json:"temperature"`
	Humidity    int     `json:"humidity"`
	TDS         int     `json:"tds"`
	PH          float64 `json:"pH"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"sensors_topic", // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Generate JSON body
	body := bodyFrom(os.Args)

	err = ch.PublishWithContext(ctx,
		"sensors_topic",       // exchange
		severityFrom(os.Args), // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent sensor data: %s", body)
}

func bodyFrom(args []string) []byte {
	data := SensorData{
		ID:          11,
		PacketNo:    126,
		Temperature: 30,
		Humidity:    60,
		TDS:         1100,
		PH:          5.0,
	}

	// Override defaults if args are provided
	if len(args) >= 3 {
		if id, err := strconv.Atoi(args[2]); err == nil {
			data.ID = id
		} else {
			log.Printf("Invalid ID, using default: %v", err)
		}
	}
	if len(args) >= 4 {
		if packetNo, err := strconv.Atoi(args[3]); err == nil {
			data.PacketNo = packetNo
		} else {
			log.Printf("Invalid PacketNo, using default: %v", err)
		}
	}
	if len(args) >= 5 {
		if temp, err := strconv.Atoi(args[4]); err == nil {
			data.Temperature = temp
		} else {
			log.Printf("Invalid Temperature, using default: %v", err)
		}
	}
	if len(args) >= 6 {
		if hum, err := strconv.Atoi(args[5]); err == nil {
			data.Humidity = hum
		} else {
			log.Printf("Invalid Humidity, using default: %v", err)
		}
	}
	if len(args) >= 7 {
		if tds, err := strconv.Atoi(args[6]); err == nil {
			data.TDS = tds
		} else {
			log.Printf("Invalid TDS, using default: %v", err)
		}
	}
	if len(args) >= 8 {
		if ph, err := strconv.ParseFloat(args[7], 64); err == nil {
			data.PH = ph
		} else {
			log.Printf("Invalid pH, using default: %v", err)
		}
	}

	body, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %s", err)
	}

	return body
}

func severityFrom(args []string) string {
	if len(args) < 2 || args[1] == "" {
		return "anonymous.info"
	}
	return args[1]
}
