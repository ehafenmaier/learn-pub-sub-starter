package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
)

func handlerWar(gs *gamelogic.GameState) func(rw gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(rw gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")

		// Connect to RabbitMQ
		conn, ch, err := pubsub.ConnectRabbitMQ()
		if err != nil {
			fmt.Printf("Failed to connect to RabbitMQ: %s\n", err)
			return pubsub.NackDiscard
		}
		defer conn.Close()
		defer ch.Close()

		outcome, _, _ := gs.HandleWar(rw)
		var ackType pubsub.AckType

		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			ackType = pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			ackType = pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			ackType = pubsub.Ack
		case gamelogic.WarOutcomeYouWon:
			ackType = pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			ackType = pubsub.Ack
		default:
			fmt.Println("Error: Unknown war outcome")
			ackType = pubsub.NackDiscard
		}

		return ackType
	}
}
