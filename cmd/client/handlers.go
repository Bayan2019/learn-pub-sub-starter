package main

import (
	"fmt"

	"github.com/Bayan2019/learn-pub-sub-starter/internal/gamelogic"
	"github.com/Bayan2019/learn-pub-sub-starter/internal/pubsub"
	"github.com/Bayan2019/learn-pub-sub-starter/internal/routing"
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
) func(gamelogic.ArmyMove) pubsub.Acktype {
	return func(move gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print("> ")
		moveOutcome := gs.HandleMove(move)
		switch moveOutcome {
		case gamelogic.MoveOutcomeSamePlayer:
			return pubsub.NackDiscard
		case gamelogic.MoveOutcomeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			return pubsub.Ack
		}
		fmt.Println("error: unknown move outcome")
		return pubsub.NackDiscard
	}
}
