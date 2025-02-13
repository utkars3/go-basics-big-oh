package main

import (
	"fmt"
	"log"
	"microservices/user-service/internal/handler"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

func main() {
	// Initialize Consul client
	client, err := api.NewClient(&api.Config{
		Address: "localhost:8500", // Replace with your Consul server address
	})
	if err != nil {
		log.Fatalf("Error creating Consul client: %v", err)
	}

	// Create a new Consul registration object
	registration := &api.AgentServiceRegistration{
		ID:      "user", // Unique ID for the service
		Name:    "user", // Service name
		Port:    8080,   // Port on which the service is running
		Tags:    []string{"api", "user"},
		Address: "127.0.0.1", // Address of the service (localhost or other IP)
	}

	// Register the service with Consul
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatalf("Error registering service with Consul: %v", err)
	}
	fmt.Println("Service registered with Consul")

	// Setup Gin router
	r := gin.Default()
	r.GET("/user/:id", handler.GetUserHandler)

	// Start the server in a separate goroutine
	go func() {
		fmt.Println("User Service is running on port 8080...")
		if err := r.Run(":8080"); err != nil {
			log.Fatal("Error starting User Service: ", err)
		}
	}()

	// Gracefully handle shutdown signals
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	// Block until an interrupt signal is received
	<-signalChannel

	// Deregister the service from Consul on shutdown
	err = client.Agent().ServiceDeregister(registration.ID)
	if err != nil {
		log.Printf("Error deregistering service from Consul: %v", err)
	} else {
		fmt.Println("Service deregistered from Consul")
	}

	// Perform any necessary cleanup before exiting
	fmt.Println("Shutting down User Service...")
}
