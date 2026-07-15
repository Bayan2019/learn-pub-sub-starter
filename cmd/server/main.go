package main

import (
	"fmt"
	"log"

	"github.com/Bayan2019/learn-pub-sub-starter/internal/pubsub"
	"github.com/Bayan2019/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	fmt.Println("Peril game server connected to RabbitMQ!")

	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// <-signalChan
	// fmt.Println("RabbitMQ connection closed.")

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}
	// Ch 3. Publishers & Queues Lv 1. Exchanges and Queues
	// use the PublishJSON function to publish a message to the exchange!
	err = pubsub.PublishJSON(
		// Use the channel you created.
		publishCh,
		// Use the internal/routing package's ExchangePerilDirect string for the exchange.
		routing.ExchangePerilDirect,
		// Use the internal/routing package's PauseKey string for the routing key.
		routing.PauseKey,
		// The data to send is the JSON-marshaled internal/routing's PlayingState struct,
		// with the IsPaused field set to true.
		routing.PlayingState{
			IsPaused: true,
		},
	)
	if err != nil {
		log.Printf("could not publish time: %v", err)
	}
	fmt.Println("Pause message sent!")
}
