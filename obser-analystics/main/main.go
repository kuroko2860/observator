package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/qiniu/qmgo"
	"kuroko.com/analystics/internal/api/handler"
	"kuroko.com/analystics/internal/api/router"
	"kuroko.com/analystics/internal/config"
	"kuroko.com/analystics/internal/service"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	client, err := qmgo.NewClient(context.Background(), &qmgo.Config{Uri: config.MONGO_URI})
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB")
	db := client.Database(config.MONGO_DATABASE)

	ctx := context.Background()
	driver, err := neo4j.NewDriverWithContext(
		config.Neo4jURI,
		neo4j.BasicAuth(config.Neo4jUsername, config.Neo4jPassword, ""))
	if err != nil {
		panic(err)
	}
	defer driver.Close(ctx)
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to Neo4j")

	s := service.NewService(db, driver)
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

	r := router.New()
	v1 := r.Group("/api")
	apiHandler := handler.NewHandler(s)
	apiHandler.RegisterRoutes(v1)
	r.Logger.Fatal(r.Start("127.0.0.1:8585"))

	// Main process logic
	fmt.Println("Application is running. Press Ctrl+C to exit.")

	if <-stopChan {
		fmt.Println("Exiting the application...")
		return
	}

}
