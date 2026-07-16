package pubsub

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Ch 3. Publishers & Queues Lv 4. Transient Queues
type Acktype int

// Ch 3. Publishers & Queues Lv 4. Transient Queues
type SimpleQueueType int

// Ch 3. Publishers & Queues Lv 4. Transient Queues
const (
	SimpleQueueDurable SimpleQueueType = iota
	SimpleQueueTransient
)

// Ch 3. Publishers & Queues Lv 4. Transient Queues
func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	// Ch 3. Publishers & Queues Lv 4. Transient Queues
	// Create a new .Channel() on the connection.
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not create channel: %v", err)
	}

	// Ch 3. Publishers & Queues Lv 5. Transient Queues
	// Declare a new queue using .QueueDeclare() method:
	queue, err := ch.QueueDeclare(
		queueName, // name
		// The durable parameter should only be true if queueType is durable.
		queueType == SimpleQueueDurable, // durable
		// The autoDelete parameter should be true if queueType is transient.
		queueType != SimpleQueueDurable, // delete when unused
		// The exclusive parameter should be true if queueType is transient.
		queueType != SimpleQueueDurable, // exclusive
		// The noWait parameter should be false.
		false, // no-wait
		// The args parameter should be nil.
		nil, // args
	)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not declare queue: %v", err)
	}

	// Ch 3. Publishers & Queues Lv 5. Transient Queues
	// Bind the queue to the exchange using .QueueBind() method.
	err = ch.QueueBind(
		queue.Name, // queue name
		key,        // routing key
		exchange,   // exchange
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not bind queue: %v", err)
	}

	// Ch 3. Publishers & Queues Lv 5. Transient Queues
	// Return the channel and queue.
	return ch, queue, nil
}

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T),
) error {
	// Ch 4. Subscribers & Routings Lv 1. Consumers
	// Call DeclareAndBind
	// to make sure that the given queue exists
	// and is bound to the exchange
	ch, queue, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		queueType,
	)
	if err != nil {
		return fmt.Errorf("could not declare and bind queue: %v", err)
	}

	// Ch 4. Subscribers & Routings Lv 1. Consumers
	// Get a new chan of amqp.Delivery structs
	// by using the channel.Consume method.
	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return fmt.Errorf("could not consume messages: %v", err)
	}

	// Ch 4. Subscribers & Routings Lv 1. Consumers
	// Unmarshal the body (raw bytes) of each message delivery
	// into the (generic) T type.
	unmarshaller := func(data []byte) (T, error) {
		var target T
		err := json.Unmarshal(data, &target)
		return target, err
	}

	// Ch 4. Subscribers & Routings Lv 1. Consumers
	// Start a goroutine that ranges over the channel of deliveries,
	// and for each message:
	go func() {
		defer ch.Close()
		for msg := range msgs {
			target, err := unmarshaller(msg.Body)
			if err != nil {
				fmt.Printf("could not unmarshal message: %v\n", err)
				continue
			}
			// Ch 4. Subscribers & Routings Lv 1. Consumers
			// Call the given handler function with the unmarshaled message
			handler(target)
			// Ch 4. Subscribers & Routings Lv 1. Consumers
			// Acknowledge the message with delivery.Ack(false)
			// to remove it from the queue
			msg.Ack(false)
		}
	}()
	return nil
}
