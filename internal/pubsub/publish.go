package pubsub

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Ch 3. Publishers & Queues Lv 1. Exchanges and Queues
// Create a new package: internal/pubsub.
// Create an exported PublishJSON function in the internal/pubsub package.
func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	// Ch 3. Publishers & Queues Lv 1. Exchanges and Queues
	// Marshal the val to JSON bytes
	dat, err := json.Marshal(val)
	if err != nil {
		return err
	}
	// Ch 3. Publishers & Queues Lv 1. Exchanges and Queues
	// Use the channel's .PublishWithContext method to publish the message
	// to the exchange with the routing key.
	return ch.PublishWithContext(
		// Set ctx to context.Background()
		context.Background(),
		exchange,
		key,
		// Set mandatory to false.
		false,
		// Set immediate to false
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        dat,
		},
	)
}
