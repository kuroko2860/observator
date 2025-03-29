package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Connect to NATS
	nc, err := nats.Connect(config.NATS_URL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	s := service.NewService(db)
	// ---------------- http logs ----------------
	// Simple Async Subscriber
	go nc.Subscribe("logs", func(m *nats.Msg) {
		s.ReceiveNATSMsg(m)
	})
	ticker := s.StartTickerUpdateData(config.INTERVAL)
	// ---------------- http logs ----------------

	// ---------------- trace data ----------------
	s.StartProcessTrace(nc)
	// ---------------- trace data ----------------

	fmt.Println("Application is running. Press Ctrl+C to exit.")

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

	if <-stopChan {
		fmt.Println("Exiting the application...")
		ticker.Stop()
		client.Close(context.Background())
		time.Sleep(1 * time.Second)
		return
	}

}
