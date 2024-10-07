package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

// 2. Message Brokers 4. Connect
func main() {
	// Declare a connection string with the value: amqp://guest:guest@localhost:5672/.
	// This is how your application will know where to connect to the RabbitMQ server.
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"
	// Call amqp.Dial with the connection string to create a new connection to RabbitMQ.
	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	// Defer a .Close() of the connection to ensure it's closed when the program exits.
	defer conn.Close()
	// Print a message to the console that the connection was successful.
	fmt.Println("Peril game server connected to RabbitMQ!")

	// If a signal is received,
	// print a message to the console that the program is shutting down and close the connection.
	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
