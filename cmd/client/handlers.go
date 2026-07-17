package main

import (
	"fmt"

	"github.com/Bayan2019/learn-pub-sub-starter/internal/gamelogic"
	"github.com/Bayan2019/learn-pub-sub-starter/internal/pubsub"
	"github.com/Bayan2019/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Ch 4. Subscribers & Routings Lv 1. Consumers
func handlerPause(
	gs *gamelogic.GameState,
) func(routing.PlayingState) pubsub.Acktype {
	// Ch 5. Delivery Lv 2. Ack and Nack
	// Update your internal/pubsub.SubscribeJSON function's handler parameter
	// to return an "acktype" instead of nothing.
	return func(playingState routing.PlayingState) pubsub.Acktype {
		// Ch 4. Subscribers & Routings Lv 1. Consumers
		// Use defer fmt.Print("> ")
		// to display a new prompt (> )
		// when the function exits.
		defer fmt.Print("> ")
		// Ch 4. Subscribers & Routings Lv 1. Consumers
		// Use the game state's HandlePause method
		// to pause the game for the client.
		gs.HandlePause(playingState)
		// Ch 5. Delivery Lv 2. Ack and Nack
		return pubsub.Ack
	}
}

// Ch 5. Delivery Lv 2. Ack and Nack
// Update your internal/pubsub.SubscribeJSON function's handler parameter
// to return an "acktype" instead of nothing.
func handlerMove(
	gs *gamelogic.GameState,
	// 5. Delivery Lv 5. Nack Requeue
	// Update the "move" handler
	// (and its registration in main.go)
	// to accept an AMQP channel.
	publishCh *amqp.Channel,
) func(gamelogic.ArmyMove) pubsub.Acktype {
	defer fmt.Print("> ")

	return func(move gamelogic.ArmyMove) pubsub.Acktype {
		moveOutcome := gs.HandleMove(move)

		switch moveOutcome {
		case gamelogic.MoveOutcomeSamePlayer:
			return pubsub.NackDiscard
		case gamelogic.MoveOutcomeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			// 5. Delivery Lv 5. Nack Requeue
			// Publish a message to the "topic" exchange with
			// the routing key $WARPREFIX.$USERNAME.
			err := pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilTopic,
				// routing.WarRecognitionsPrefix contains
				// the $WARPREFIX constant
				// The $USERNAME should be
				// the name of the player consuming the move.
				routing.WarRecognitionsPrefix+"."+gs.GetUsername(),
				gamelogic.RecognitionOfWar{
					Attacker: move.Player,
					Defender: gs.GetPlayerSnap(),
				},
			)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				// If publishing the war declaration fails, "NackRequeue" the message.
				return pubsub.NackRequeue
			}
			// 5. Delivery Lv 5. Nack Requeue
			// NackRequeue the message...
			// return pubsub.NackRequeue
			// 5. Delivery Lv 6. Nack Requeue Fix
			// Otherwise, "Ack" the message.
			return pubsub.Ack

		}
		fmt.Println("error: unknown move outcome")
		return pubsub.NackDiscard
	}
}

// 5. Delivery Lv 5. Nack Requeue
// Create a new handler that consumes
// all the war messages that the "move" handler publishes,
// no matter the username in the routing key.
// Ch 6. Serialization Lv 2. Game Logs
// Update the war handler function in the client to publish game logs.
func handlerWar(
	gs *gamelogic.GameState,
	// Ch 6. Serialization Lv 2. Game Logs
	publishCh *amqp.Channel,
) func(dw gamelogic.RecognitionOfWar) pubsub.Acktype {
	return func(dw gamelogic.RecognitionOfWar) pubsub.Acktype {
		// 5. Delivery Lv 5. Nack Requeue
		// to ensure a new prompt is printed after the handler is done.
		defer fmt.Print("> ")
		// 5. Delivery Lv 5. Nack Requeue
		// Call the gamestate's HandleWar method with the message's body.
		// Ch 6. Serialization Lv 2. Game Logs
		// Capture the winner and loser return values
		// from the GameState's HandleWar method
		warOutcome, winner, loser := gs.HandleWar(dw)
		switch warOutcome {
		case gamelogic.WarOutcomeNotInvolved:
			// NackRequeue the message so
			// another client can try to consume it.
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			// NackDiscard the message.
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			// Ack the message.
			// Ch 6. Serialization Lv 2. Game Logs
			err := publishGameLog(
				publishCh,
				gs.GetUsername(),
				fmt.Sprintf(
					"%s won a war against %s",
					winner,
					loser,
				),
			)
			if err != nil {
				// Ch 6. Serialization Lv 2. Game Logs
				// If a publishing fails, NackRequeue,
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeYouWon:
			// Ack the message.
			// Ch 6. Serialization Lv 2. Game Logs
			err := publishGameLog(
				publishCh,
				gs.GetUsername(),
				fmt.Sprintf(
					"%s won a war against %s",
					winner,
					loser,
				),
			)
			if err != nil {
				// Ch 6. Serialization Lv 2. Game Logs
				// If a publishing fails, NackRequeue,
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			// Ack the message.
			// Ch 6. Serialization Lv 2. Game Logs
			err := publishGameLog(
				publishCh,
				gs.GetUsername(),
				fmt.Sprintf(
					"A war between %s and %s resulted in a draw",
					winner,
					loser,
				),
			)
			if err != nil {
				// Ch 6. Serialization Lv 2. Game Logs
				// If a publishing fails, NackRequeue,
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		}
		// if it's anything else,
		// print an error
		fmt.Println("error: unknown war outcome")
		// and NackDiscard the message.
		return pubsub.NackDiscard
	}
}
