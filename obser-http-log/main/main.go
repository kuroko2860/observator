package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/qiniu/qmgo"
	"kuroko.com/processor/internal/config"
	"kuroko.com/processor/internal/service"
)

func main() {
	client, err := qmgo.NewClient(context.Background(), &qmgo.Config{Uri: config.MONGO_URI})
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB")
	db := client.Database(config.MONGO_DATABASE)

	// Connect to a server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to NATS")

	s := service.NewService(db)
	ticker := s.StartTickerUpdateData(config.INTERVAL)

	// Simple Async Subscriber
	go nc.Subscribe("logs", func(m *nats.Msg) {
		s.ReceiveNATSMsg(m)
	})
	// Create a channel to receive OS signals
	signalChan := make(chan os.Signal, 1)
	// Notify the channel of specific signals
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a channel to indicate when to stop the program
	stopChan := make(chan bool)

	// Goroutine to listen for signals
	go func() {
		for sig := range signalChan {
			fmt.Printf("Received signal: %s\n", sig)
			// Implement custom logic on signal reception
			if sig == syscall.SIGTERM {
				fmt.Println("SIGTERM received, cleaning up...")
				// Perform cleanup tasks here if needed
			} else if sig == syscall.SIGINT {
				fmt.Println("SIGINT received, gracefully shutting down...")
				stopChan <- true
			}
		}
	}()

	// Main process logic
	fmt.Println("Application is running. Press Ctrl+C to exit.")

	if <-stopChan {
		fmt.Println("Exiting the application...")
		client.Close(context.Background())
		ticker.Stop()
		nc.Close()
		return
	}

}
