package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SensorData struct {
	ID          int     `json:"id"`
	PacketNo    int     `json:"packet_no"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	TDS         int     `json:"tds"`
	PH          float64 `json:"pH"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	// Check for at least two arguments: queue name and at least one binding key
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s [queue_name] [binding_key]...", os.Args[0])
	}

	// Parse command-line arguments
	queueName := os.Args[1]
	bindingKeys := os.Args[2:]

	// Connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare an exchange
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

	// Declare a queue based on the queue name argument
	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Bind the queue with each specified binding key
	for _, bindingKey := range bindingKeys {
		log.Printf("Binding queue %s to exchange %s with routing key %s", queueName, "sensors_topic", bindingKey)
		err = ch.QueueBind(
			q.Name,          // queue name
			bindingKey,      // routing key
			"sensors_topic", // exchange
			false,
			nil,
		)
		failOnError(err, "Failed to bind a queue")
	}

	// Consume messages from the queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	// Set up a channel to block and receive messages
	var forever chan struct{}

	go func() {
		for d := range msgs {
			var data SensorData
			err := json.Unmarshal(d.Body, &data)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
				continue
			}
			fmt.Printf("[%s] Received Sensor Data:\n", queueName)
			fmt.Printf("ID: %d, PacketNo: %d, Temperature: %.2f, Humidity: %.2f, TDS: %d, pH: %.2f\n",
				data.ID, data.PacketNo, data.Temperature, data.Humidity, data.TDS, data.PH)
		}
	}()

	log.Printf(" [*] Waiting for messages on queue %s. To exit press CTRL+C", queueName)
	<-forever
}
