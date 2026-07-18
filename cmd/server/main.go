package main

import (
	"fmt"
	"log"

	"github.com/Bayan2019/learn-pub-sub-starter/internal/gamelogic"
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

	// Ch 3. Publishers & Queues Lv 9. Durable
	// to declare and bind a queue to the new peril_topic exchange.
	// _, queue, err := pubsub.DeclareAndBind(
	// 	conn,
	// 	// new peril_topic exchange
	// 	routing.ExchangePerilTopic,
	// 	// named game_logs.
	// 	routing.GameLogSlug,
	// 	// The routing key should be game_logs.*.
	// 	routing.GameLogSlug+".*",
	// 	// It should be a durable queue
	// 	pubsub.SimpleQueueDurable,
	// )
	// if err != nil {
	// 	log.Fatalf("could not subscribe to pause: %v", err)
	// }
	// fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	// Ch 6. Serialization Lv 3. Consume Logs
	// Update the server to SubscribeGob
	// to the game_logs queue instead of just declaring it.
	err = pubsub.SubscribeGob(
		conn,
		// new peril_topic exchange
		routing.ExchangePerilTopic,
		// named game_logs.
		routing.GameLogSlug,
		// The routing key should be game_logs.*.
		routing.GameLogSlug+".*",
		// Use a durable queue with the war handler.
		pubsub.SimpleQueueDurable,
		handlerLogs(),
	)
	if err != nil {
		log.Fatalf("could not subscribe to war declarations: %v", err)
	}

	// Ch 3. Publishers & Queues Lv 5. Decoupling
	// Run the PrintServerHelp function in internal/gamelogic
	// as the server starts up so that you can see the commands
	// the user of the REPL can use.
	gamelogic.PrintServerHelp()

	// Ch 3. Publishers & Queues Lv 5. Decoupling
	// Start an infinite loop
	for {
		// Ch 3. Publishers & Queues Lv 5. Decoupling
		// use the GetInput function in internal/gamelogic
		// to wait for a slice of input "words" from the user.
		inputs := gamelogic.GetInput()

		// Ch 3. Publishers & Queues Lv 5. Decoupling
		// If the slice is empty, continue to the next iteration of the loop.
		if len(inputs) == 0 {
			continue
		} else {
			// Ch 3. Publishers & Queues Lv 5. Decoupling
			// Check the first word:
			switch inputs[0] {
			case "pause":
				// Ch 3. Publishers & Queues Lv 5. Decoupling
				// If it's "pause", log to the console
				// that you're sending a pause message,
				// and publish the pause message as you were doing before.
				log.Println("Publishing paused game state")
				// Ch 3. Publishers & Queues Lv 1. Exchanges and Queues
				// use the PublishJSON function
				// to publish a message to the exchange!
				err = pubsub.PublishJSON(
					// Use the channel you created.
					publishCh,
					// Use the internal/routing package's
					// ExchangePerilDirect string for the exchange.
					routing.ExchangePerilDirect,
					// Use the internal/routing package's
					// PauseKey string for the routing key.
					routing.PauseKey,
					// The data to send is the JSON-marshaled internal/routing's
					// PlayingState struct,
					// with the IsPaused field set to true.
					routing.PlayingState{
						IsPaused: true,
					},
				)
				if err != nil {
					log.Printf("could not publish time: %v", err)
				}
			case "resume":
				// Ch 3. Publishers & Queues Lv 5. Decoupling
				// If it's "resume", log to the console that you're sending a resume message,
				// and publish the resume message as you were doing before.
				// The only difference is that the IsPaused field should be set to false.
				log.Println("Publishing resumes game state")
				err = pubsub.PublishJSON(
					publishCh,
					routing.ExchangePerilDirect,
					routing.PauseKey,
					routing.PlayingState{
						IsPaused: false,
					},
				)
				if err != nil {
					log.Printf("could not publish time: %v", err)
				}
			case "quit":
				// Ch 3. Publishers & Queues Lv 5. Decoupling
				// If it's "quit",
				// log to the console that you're exiting, and break out of the loop.
				log.Println("Quiting the server")
				// break OuterLoop
				return
			default:
				log.Printf("We don't understand command: %s\n", inputs[0])
			}
		}
	}

}
