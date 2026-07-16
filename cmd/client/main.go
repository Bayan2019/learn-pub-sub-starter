package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Bayan2019/learn-pub-sub-starter/internal/gamelogic"
	"github.com/Bayan2019/learn-pub-sub-starter/internal/pubsub"
	"github.com/Bayan2019/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Ch 3. Publishers & Queues Lv 4. Transient Queues
func main() {
	fmt.Println("Starting Peril client...")

	// Ch 3. Publishers & Queues Lv 4. Transient Queues
	// to connect to Rabbit, similar to the cmd/server package.
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game client connected to RabbitMQ!")

	// Ch 3. Publishers & Queues Lv 4. Transient Queues
	// Use the ClientWelcome() function in internal/gamelogic
	// to prompt the user for a username.
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could not get username: %v", err)
	}

	// Ch 3. Publishers & Queues Lv 4. Transient Queues
	// Declare and bind a transient queue
	_, queue, err := pubsub.DeclareAndBind(
		conn,
		// exchange: peril_direct (this is a constant in the internal/routing package)
		routing.ExchangePerilDirect,
		// queueName: pause.username where username is the user's input.
		// The pause section of the name is the routing key constant
		// in the internal/routing package
		// and is joined by a ..
		routing.PauseKey+"."+username,
		// routingKey: pause (this is a constant in the internal/routing package)
		routing.PauseKey,
		// queueType: transient
		pubsub.SimpleQueueTransient,
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	// Ch 3. Publishers & Queues Lv 4. Transient Queues
	// wait for ctrl+c
	// After declaring and binding the queue,
	// the client should wait for a Ctrl+C signal to exit.
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("RabbitMQ connection closed.")
}
