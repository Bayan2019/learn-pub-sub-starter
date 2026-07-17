package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
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

// Ch 6. Serialization Lv 2. Game Logs
// Add a PublishGob function to the internal/pubsub package.
func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	// Ch 6. Serialization Lv 2. Game Logs
	// It should be similar to the PublishJSON function, but encode to gob
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(val)
	if err != nil {
		return err
	}
	// Ch 6. Serialization Lv 2. Game Logs
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
			// Ch 6. Serialization Lv 2. Game Logs
			// Set the ContentType option to application/gob
			ContentType: "application/json",
			Body:        buffer.Bytes(),
		},
	)
}
