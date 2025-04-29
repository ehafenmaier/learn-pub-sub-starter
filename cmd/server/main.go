package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	// Connect to RabbitMQ
	connStr := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connStr)
	if err != nil {
		fmt.Printf("Failed to connect to RabbitMQ: %s\n", err)
		return
	}
	defer conn.Close()
	fmt.Println("Successfully connected to RabbitMQ")

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("Failed to open a channel: %s\n", err)
		return
	}
	defer ch.Close()

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
