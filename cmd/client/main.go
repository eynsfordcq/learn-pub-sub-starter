package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

	connStr := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Fatalf("fail to connect to rabbit mq: %v", err)
	}
	defer conn.Close()
	fmt.Println("Successfully connected to RabbitMQ.")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel: %v", err)
	}
	defer ch.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("fail to get username: %v", err)
	}

	gameState := gamelogic.NewGameState(username)

	// subscribe to pause messages
	pauseQueueName := fmt.Sprintf("%s.%s", routing.PauseKey, username)
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		pauseQueueName,
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
		handlerPause(gameState),
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause queue: %v", err)
	}
	fmt.Printf("subscribed to pause messages on queue: %s\n", pauseQueueName)

	// subscribe to army moves
	armyMovesQueueName := fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, username)
	armyMovesBindingKey := fmt.Sprintf("%s.*", routing.ArmyMovesPrefix)
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		armyMovesQueueName,
		armyMovesBindingKey,
		pubsub.SimpleQueueTransient,
		handlerMove(gameState),
	)
	if err != nil {
		log.Fatalf("could not subscribe to army moves: %v", err)
	}
	fmt.Printf("subscribed to army moves with binding key %s on queue: %s\n",
		armyMovesBindingKey,
		armyMovesQueueName,
	)

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		input := words[0]
		switch input {
		case "spawn":
			err = gameState.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
			}

		case "move":
			move, err := gameState.CommandMove(words)
			if err != nil {
				fmt.Println(err)
			}

			err = pubsub.PublishJson(
				ch,
				routing.ExchangePerilTopic,
				fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, username),
				move,
			)
			if err != nil {
				log.Printf("fail to publish move: %v", err)
				continue
			}
			log.Printf("moved %d units to %s",
				len(move.Units),
				move.ToLocation,
			)

		case "status":
			gameState.CommandStatus()

		case "help":
			gamelogic.PrintClientHelp()

		case "spam":
			fmt.Println("Spamming not allowed yet!")

		case "quit":
			gamelogic.PrintQuit()
			return

		default:
			fmt.Println("not a valid command")
		}
	}

}
