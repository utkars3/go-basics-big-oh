package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

type ServiceRegistry struct {
	ConsulClient *api.Client
}

func NewServiceRegistry() (*ServiceRegistry, error) {
	// Initialize Consul client
	client, err := api.NewClient(&api.Config{
		Address: "localhost:8500", // Consul address (adjust as needed)
	})
	if err != nil {
		return nil, fmt.Errorf("error creating Consul client: %v", err)
	}
	return &ServiceRegistry{ConsulClient: client}, nil
}

// proxyRequest handles proxying the requests to the appropriate service using Consul for service discovery
func (sr *ServiceRegistry) proxyRequest(c *gin.Context) {
	serviceName := c.Param("service")

	// Query Consul for service information
	services, _, err := sr.ConsulClient.Catalog().Service(serviceName, "", nil)
	if err != nil || len(services) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found in Consul"})
		return
	}

	// Assuming we're using the first available service (you could implement load balancing if necessary)
	service := services[0]
	url := fmt.Sprintf("http://%s:%d", service.ServiceAddress, service.ServicePort)

	// Proxy the request to the discovered service
	resp, err := http.Get(url + c.Request.URL.Path)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service unavailable"})
		return
	}
	defer resp.Body.Close()

	c.DataFromReader(resp.StatusCode, resp.ContentLength, "application/json", resp.Body, nil)
}

func main() {
	// Initialize service registry with Consul client
	registry, err := NewServiceRegistry()
	if err != nil {
		log.Fatalf("Error initializing service registry: %v", err)
	}

	// Set up Gin router
	r := gin.Default()

	// Proxy requests using the service name (from URL path)
	r.GET("/:service/*path", registry.proxyRequest)

	log.Println("API Gateway running on port 5000")
	r.Run(":5000")
}
