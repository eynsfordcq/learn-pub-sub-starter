package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("fail to get username: %v", err)
	}

	queueName := fmt.Sprintf("%s.%s", routing.PauseKey, username)

	_, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		queueName,
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
	)
	if err != nil {
		log.Fatalf("fail to bind queue: %v", err)
	}

	fmt.Printf("Queue '%s' declared and bound to exchange '%s' with routing key '%s'\n",
		queue.Name,
		routing.ExchangePerilDirect,
		routing.PauseKey,
	)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	fmt.Println("\nShutting down...")

}
