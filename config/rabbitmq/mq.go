package rabbitmq

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func ConnectRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
		return nil, nil, err
	}

	return conn, ch, nil
}

func PublishMessage(ch *amqp.Channel, queueName string, message string) error {
	fmt.Println(message)
	// Declaring a queue
	q, err := ch.QueueDeclare(
		queueName, // Nama queue
		false,     // Durable
		false,     // Auto-delete
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		return err
	}

	// Publish a message to the queue
	err = ch.Publish(
		"",     // Exchange
		q.Name, // Routing key (nama queue)
		false,  // Mandatory
		false,  // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
		return err
	}

	log.Printf(" [x] Sent %s", message)
	return nil
}

func ConsumeMessages(ch *amqp.Channel, queueName string) (<-chan amqp.Delivery, error) {
	// Declare a queue (the queue must be the same as when publishing the message)
	q, err := ch.QueueDeclare(
		queueName, // Nama queue
		false,     // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		return nil, err
	}

	// Retrieving messages from queue
	msgs, err := ch.Consume(
		q.Name, // Nama queue
		"",     // Consumer
		true,   // Auto-ack
		false,  // Exclusive
		false,  // No-local
		false,  // No-wait
		nil,    // Arguments
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
