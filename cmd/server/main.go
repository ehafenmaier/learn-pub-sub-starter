package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func main() {
	fmt.Println("Starting Peril server...")

	// Connect to RabbitMQ
	conn, ch, err := pubsub.ConnectRabbitMQ()
	if err != nil {
		fmt.Printf("Failed to connect to RabbitMQ: %s\n", err)
		return
	}
	defer conn.Close()
	defer ch.Close()
	fmt.Println("Successfully connected to RabbitMQ")

	// Subscribe to the game logs queue
	err = pubsub.SubscribeGob(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.DurableQueue,
		handlerGameLog())
	if err != nil {
		fmt.Printf("Failed to subscribe to game logs queue: %s\n", err)
		return
	}

	// Print server help
	gamelogic.PrintServerHelp()

	// Start the game server loop
	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "pause":
			fmt.Println("Pausing game...")
			msg := routing.PlayingState{IsPaused: true}
			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, msg)
			if err != nil {
				fmt.Printf("Failed to publish message: %s\n", err)
			}
		case "resume":
			fmt.Println("Resuming game...")
			msg := routing.PlayingState{IsPaused: false}
			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, msg)
			if err != nil {
				fmt.Printf("Failed to publish message: %s\n", err)
			}
		case "quit":
			fmt.Println("Shutting down Peril server...")
			return
		case "help":
			gamelogic.PrintServerHelp()
		default:
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}
}
