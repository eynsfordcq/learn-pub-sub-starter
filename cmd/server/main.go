package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	ampq "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connStr := "amqp://guest:guest@localhost:5672/"
	conn, err := ampq.Dial(connStr)
	if err != nil {
		log.Fatalf("fail to connect to rabbit mq: %v", err)
	}
	defer conn.Close()

	fmt.Println("Successfully connected to RabbitMQ.")
	fmt.Println("  [*] Waiting for messages. To exit, press CTRL+C")

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("fail to create channel: %v", err)
	}

	gamelogic.PrintServerHelp()

	for {
		inputSlice := gamelogic.GetInput()
		if len(inputSlice) == 0 {
			continue
		}

		input := inputSlice[0]
		switch input {
		case "pause":
			fmt.Println("sending pause message...")
			err = pubsub.PublishJson(
				channel,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: true,
				},
			)
			if err != nil {
				log.Printf("could not publish data: %v", err)
			}

		case "resume":
			fmt.Println("sending resume message...")
			err = pubsub.PublishJson(
				channel,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: false,
				},
			)
			if err != nil {
				log.Printf("could not publish data: %v", err)
			}

		case "quit":
			fmt.Println("exiting...")
			return

		default:
			fmt.Printf("not a valid command: %s\n", input)
		}
	}
}
