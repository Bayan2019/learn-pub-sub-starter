package main

import (
	"fmt"

	"github.com/Bayan2019/learn-pub-sub-starter/internal/gamelogic"
	"github.com/Bayan2019/learn-pub-sub-starter/internal/pubsub"
	"github.com/Bayan2019/learn-pub-sub-starter/internal/routing"
)

// Ch 6. Serialization Lv 3. Consume Logs
// The handler should:
func handlerLogs() func(gamelog routing.GameLog) pubsub.Acktype {
	return func(gamelog routing.GameLog) pubsub.Acktype {
		// Defer printing a new prompt to the console.
		defer fmt.Print("> ")
		// Use the gamelogic.WriteLog function to write the log to disk.
		err := gamelogic.WriteLog(gamelog)
		if err != nil {
			fmt.Printf("error writing log: %v\n", err)
			return pubsub.NackRequeue
		}
		return pubsub.Ack
	}
}
