package main

import (
	"context"
	"encoding/json"
	"log"
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

	q, err := ch.QueueDeclare(
		"ktmt", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Create the sensor data
	data := SensorData{
		ID:          11,
		PacketNo:    126,
		Temperature: 30,
		Humidity:    60,
		TDS:         1100,
		PH:          5.0,
	}

	// Marshal the data into JSON format
	body, err := json.Marshal(data)
	failOnError(err, "Failed to encode JSON")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}
