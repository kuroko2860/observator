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

	echoSwagger "github.com/swaggo/echo-swagger"
	_ "kuroko.com/analystics/docs"
)

// @title			Todo Application
// @description	This is a todo list management application
// @version		1.0
// @host			localhost:8585
// @BasePath		/api
func main() {
	client, err := qmgo.NewClient(context.Background(), &qmgo.Config{Uri: config.MONGO_URI})
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB")
	db := client.Database(config.MONGO_DATABASE)

	s := service.NewService(db)
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

	go func() {
		if <-stopChan {
			fmt.Println("Exiting the application...")
			os.Exit(0)
		}
	}()
	r := router.New()
	v1 := r.Group("/api")
	apiHandler := handler.NewHandler(s)
	apiHandler.RegisterRoutes(v1)
	r.GET("/swagger/*", echoSwagger.WrapHandler)
	r.Logger.Fatal(r.Start("127.0.0.1:8585"))

	// Main process logic
	fmt.Println("Application is running. Press Ctrl+C to exit.")

}
