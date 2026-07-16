package main

import (
	"fmt"
	"log"

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
	// _, queue, err := pubsub.DeclareAndBind(
	// 	conn,
	// 	// exchange: peril_direct (this is a constant in the internal/routing package)
	// 	routing.ExchangePerilDirect,
	// 	// queueName: pause.username where username is the user's input.
	// 	// The pause section of the name is the routing key constant
	// 	// in the internal/routing package
	// 	// and is joined by a ..
	// 	routing.PauseKey+"."+username,
	// 	// routingKey: pause (this is a constant in the internal/routing package)
	// 	routing.PauseKey,
	// 	// queueType: transient
	// 	pubsub.SimpleQueueTransient,
	// )
	// if err != nil {
	// 	log.Fatalf("could not subscribe to pause: %v", err)
	// }
	// fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	// Ch 3. Publishers & Queues Lv 4. Transient Queues
	// wait for ctrl+c
	// After declaring and binding the queue,
	// the client should wait for a Ctrl+C signal to exit.
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// <-signalChan
	// fmt.Println("RabbitMQ connection closed.")

	// Ch 3. Publishers & Queues Lv 6. Client REPL
	// use the NewGameState function in internal/gamelogic
	// to create a new game state.
	gs := gamelogic.NewGameState(username)

	// Ch 4. Subscribers & Routings Lv 1. Consumers
	// after creating the game state,
	// replace your previous DeclareAndBind call
	// with pubsub.SubscribeJSON
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+gs.GetUsername(),
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
		handlerPause(gs),
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}

	// Ch 3. Publishers & Queues Lv 6. Client REPL
	// Add a REPL loop similar to what you did
	// in the cmd/server application.
	for {
		// Ch 3. Publishers & Queues Lv 6. Client REPL
		// the "words" from the GetInput command.
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "move":
			// Ch 3. Publishers & Queues Lv 6. Client REPL
			// The move command allows a player
			// to move their units to a new location.
			// It accepts two arguments:
			// the destination, and the ID of the unit.
			// Call the gamestate.CommandMove method
			_, err := gs.CommandMove(words)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// TODO: publish the move
		case "spawn":
			// Ch 3. Publishers & Queues Lv 6. Client REPL
			// The spawn command allows a player
			// to add a new unit to the map under their control.
			// Use the gamestate.CommandSpawn method
			err = gs.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "status":
			// Ch 3. Publishers & Queues Lv 6. Client REPL
			// The status command uses the gamestate.CommandStatus method
			gs.CommandStatus()
		case "help":
			// Ch 3. Publishers & Queues Lv 6. Client REPL
			// The help command uses the gamelogic.PrintClientHelp function
			gamelogic.PrintClientHelp()
		case "spam":
			// TODO: publish n malicious logs
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			// Ch 3. Publishers & Queues Lv 6. Client REPL
			// The quit command uses the gamelogic.PrintQuit function
			// to print a message, then exit the REPL.
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("unknown command")
		}
	}
}
