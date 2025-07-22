package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	ampq "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connStr := "amqp://guest:guest@localhost:5672/"
	conn, err := ampq.Dial(connStr)
	if err != nil {
		fmt.Println("fail to connect to rabbit mq", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Successfully connected to RabbitMQ.")
	fmt.Println("  [*] Waiting for messages. To exit, press CTRL+C")

	// create a channel to receive signals
	signals := make(chan os.Signal, 1)
	// notify the channel for SIGINT (CTRL+C) and SIGTERM (kill)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	// block until a signal is received
	<-signals
	fmt.Println("\nShutting down...")

}
