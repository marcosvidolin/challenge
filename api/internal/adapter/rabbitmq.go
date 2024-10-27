package adapter

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// rabbitMQAdapter defines an adapter for RabbitMQ
type rabbitmqAdapter struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

// NewRabbitMQAdapter creates a RabbitMQ adapter
// returns an error if any issue connectin to the
// server occours
func NewRabbitMQAdapter(url, queeu string) (*rabbitmqAdapter, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create RabbitMQ channel: %w", err)
	}

	queue, err := channel.QueueDeclare(queeu, true, false, false, false, nil)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare RabbitMQ queue: %w", err)
	}

	return &rabbitmqAdapter{
		conn:    conn,
		channel: channel,
		queue:   queue,
	}, nil
}

// Publish sends a message to the queue
func (r *rabbitmqAdapter) Publish(message, contentType string) error {
	err := r.channel.Publish(
		"",           // Exchange
		r.queue.Name, // Routing key (queue name)
		false,        // Mandatory
		false,        // Immediate
		amqp.Publishing{
			ContentType: contentType,
			Body:        []byte(message),
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}

// Consume starts consuming messages from queue
func (r *rabbitmqAdapter) Consume(ch chan<- string) error {
	msgs, err := r.channel.Consume(
		r.queue.Name, // Queue name
		"",           // Consumer
		true,         // Auto-ack
		false,        // Exclusive
		false,        // No-local
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming messages: %w", err)
	}

	go func() {
		for m := range msgs {
			ch <- string(m.Body)
		}
	}()

	return nil
}

// Close closes the connection and channel to RabbitMQ
func (r *rabbitmqAdapter) Close() {
	r.channel.Close()
	r.conn.Close()
}
