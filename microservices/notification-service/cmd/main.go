// package main

// import (
// 	"fmt"
// 	"log"
// 	"microservices/notification-service/internal/handler"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"

// 	"github.com/gin-gonic/gin"
// 	"github.com/hashicorp/consul/api"
// )

// type ServiceRegistry struct {
// 	ConsulClient *api.Client
// }

// func NewServiceRegistry() (*ServiceRegistry, error) {
// 	// Initialize Consul client
// 	client, err := api.NewClient(&api.Config{
// 		Address: "localhost:8500", // Consul address (adjust as needed)
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("error creating Consul client: %v", err)
// 	}
// 	return &ServiceRegistry{ConsulClient: client}, nil
// }

// // registerService registers the notification service with Consul
// func (sr *ServiceRegistry) registerService() error {
// 	registration := &api.AgentServiceRegistration{
// 		ID:      "notification", // Unique ID for the service
// 		Name:    "notification", // Service name
// 		Port:    8081,           // Port on which the notification service is running
// 		Tags:    []string{"notification"},
// 		Address: "127.0.0.1", // Address of the service (localhost or other IP)
// 	}

// 	// Register the service with Consul
// 	err := sr.ConsulClient.Agent().ServiceRegister(registration)
// 	if err != nil {
// 		return fmt.Errorf("error registering service with Consul: %v", err)
// 	}
// 	fmt.Println("Notification service registered with Consul")
// 	return nil
// }

// // deregisterService deregisters the notification service from Consul
// func (sr *ServiceRegistry) deregisterService() error {
// 	err := sr.ConsulClient.Agent().ServiceDeregister("notification")
// 	if err != nil {
// 		return fmt.Errorf("error deregistering service from Consul: %v", err)
// 	}
// 	fmt.Println("Notification service deregistered from Consul")
// 	return nil
// }

// // handler for the notification service (just as an example)
// func GetNotificationHandler(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"message": "Notification service is running!"})
// }

// func main() {
// 	// Initialize service registry with Consul client
// 	registry, err := NewServiceRegistry()
// 	if err != nil {
// 		log.Fatalf("Error initializing service registry: %v", err)
// 	}

// 	// Register the notification service with Consul on startup
// 	err = registry.registerService()
// 	if err != nil {
// 		log.Fatalf("Error registering service with Consul: %v", err)
// 	}

// 	// Set up Gin router
// 	r := gin.Default()

// 	// Define routes
// 	r.POST("/notification", handler.NotifyHandler)

// 	// Gracefully handle shutdown signals
// 	signalChannel := make(chan os.Signal, 1)
// 	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

// 	// Start the server in a goroutine
// 	go func() {
// 		log.Println("Notification service running on port 8081")
// 		if err := r.Run(":8081"); err != nil {
// 			log.Fatal("Error starting Notification Service: ", err)
// 		}
// 	}()

// 	// Block until an interrupt signal is received
// 	<-signalChannel

// 	// Deregister the service from Consul on shutdown
// 	err = registry.deregisterService()
// 	if err != nil {
// 		log.Printf("Error deregistering service from Consul: %v", err)
// 	}

// 	// Perform any necessary cleanup before exiting
// 	fmt.Println("Shutting down Notification Service...")
// }

package main

import (
	"fmt"
	"log"
	"microservices/notification-service/internal/handler"

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
		ID:      "notification", // Unique ID for the service
		Name:    "notification", // Service name
		Port:    8081,           // Port on which the service is running
		Tags:    []string{"notification"},
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
	r.POST("/notification", handler.NotifyHandler)

	// Start the server in a separate goroutine
	go func() {
		fmt.Println("Notification Service is running on port 8081...")
		if err := r.Run(":8081"); err != nil {
			log.Fatal("Error starting Notification Service: ", err)
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
