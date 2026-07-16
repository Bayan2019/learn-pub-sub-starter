package pubsub

import (
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
