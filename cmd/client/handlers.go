package main

import (
	"fmt"

	"github.com/Bayan2019/learn-pub-sub-starter/internal/gamelogic"
	"github.com/Bayan2019/learn-pub-sub-starter/internal/routing"
)

// Ch 4. Subscribers & Routings Lv 1. Consumers
func handlerPause(
	gs *gamelogic.GameState,
) func(routing.PlayingState) {

	return func(playingState routing.PlayingState) {
		// Ch 4. Subscribers & Routings Lv 1. Consumers
		// Use defer fmt.Print("> ")
		// to display a new prompt (> )
		// when the function exits.
		defer fmt.Print("> ")
		// Ch 4. Subscribers & Routings Lv 1. Consumers
		// Use the game state's HandlePause method
		// to pause the game for the client.
		gs.HandlePause(playingState)
	}
}

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) {
	return func(move gamelogic.ArmyMove) {
		defer fmt.Print("> ")
		gs.HandleMove(move)
	}
}
