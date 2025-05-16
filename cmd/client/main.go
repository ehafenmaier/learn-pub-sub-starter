package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"strconv"
	"time"
)

func main() {
	fmt.Println("Starting Peril client...")

	// Connect to RabbitMQ
	conn, ch, err := pubsub.ConnectRabbitMQ()
	if err != nil {
		fmt.Printf("Failed to connect to RabbitMQ: %s\n", err)
		return
	}
	defer conn.Close()
	defer ch.Close()
	fmt.Println("Successfully connected to RabbitMQ")

	// Prompt user for username
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Printf("Failed to get username: %s\n", err)
		return
	}

	// Create new game state
	gameState := gamelogic.NewGameState(username)

	// Subscribe to the pause queue
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.TransientQueue,
		handlerPause(gameState))
	if err != nil {
		fmt.Printf("Failed to subscribe to pause queue: %s\n", err)
		return
	}

	// Subscribe to the move queue
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+username,
		routing.ArmyMovesPrefix+".*",
		pubsub.TransientQueue,
		handlerMove(gameState, ch))
	if err != nil {
		fmt.Printf("Failed to subscribe to move queue: %s\n", err)
		return
	}

	// Subscribe to the war queue
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.WarRecognitionsPrefix,
		routing.WarRecognitionsPrefix+".*",
		pubsub.DurableQueue,
		handlerWar(gameState, ch))
	if err != nil {
		fmt.Printf("Failed to subscribe to war queue: %s\n", err)
		return
	}

	// Start the game client REPL loop
	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "spawn":
			err := gameState.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "move":
			mv, err := gameState.CommandMove(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = pubsub.PublishJSON(ch, routing.ExchangePerilTopic, routing.ArmyMovesPrefix+"."+username, mv)
			if err != nil {
				fmt.Printf("Failed to publish move message: %s\n", err)
				continue
			}
			fmt.Printf("Move to %s successful!\n", mv.ToLocation)
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			if len(words) < 2 {
				fmt.Println("Usage: spam <message>")
				continue
			}
			spamCount, err := strconv.Atoi(words[1])
			if err != nil {
				fmt.Println("Invalid spam message. Must be a number.")
				continue
			}
			for i := 0; i < spamCount; i++ {
				err = pubsub.PublishGob(
					ch,
					routing.ExchangePerilTopic,
					routing.GameLogSlug+"."+username,
					routing.GameLog{
						Username:    username,
						Message:     gamelogic.GetMaliciousLog(),
						CurrentTime: time.Now(),
					})
				if err != nil {
					fmt.Printf("Failed to publish spam message: %s\n", err)
					continue
				}
			}
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}
}
