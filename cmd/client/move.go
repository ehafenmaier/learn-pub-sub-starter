package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.AckType {
	return func(move gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Print("> ")

		// Connect to RabbitMQ
		conn, ch, err := pubsub.ConnectRabbitMQ()
		if err != nil {
			fmt.Printf("Failed to connect to RabbitMQ: %s\n", err)
			return pubsub.NackDiscard
		}
		defer conn.Close()
		defer ch.Close()

		outcome := gs.HandleMove(move)
		var ackType pubsub.AckType

		switch outcome {
		case gamelogic.MoveOutComeSafe:
			ackType = pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			key := routing.WarRecognitionsPrefix + "." + gs.GetUsername()
			err = pubsub.PublishJSON(ch, routing.ExchangePerilTopic, key, move)
			if err != nil {
				ackType = pubsub.NackRequeue
			} else {
				ackType = pubsub.Ack
			}
		case gamelogic.MoveOutcomeSamePlayer:
			ackType = pubsub.NackDiscard
		default:
			ackType = pubsub.NackDiscard
		}

		return ackType
	}
}
