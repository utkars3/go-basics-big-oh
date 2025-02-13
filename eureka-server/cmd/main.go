package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	// Store registered services in memory
	services = make(map[string]Service)
	mu       sync.Mutex
)

type Service struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Host     string   `json:"host"`
	Port     int      `json:"port"`
	Health   string   `json:"health"`
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	URL string `json:"url"`
}

func main() {
	// Register routes
	http.HandleFunc("/eureka/apps", registerService)
	http.HandleFunc("/eureka/apps/", discoverService)
	http.HandleFunc("/eureka/apps/deregister", deregisterService)
	http.HandleFunc("/eureka/apps/log", logRegisteredServices) //

	// Start the server
	log.Println("Mock Eureka Server started on port 8761...")
	log.Fatal(http.ListenAndServe(":8761", nil))
}

// Register service in the mock Eureka server
func registerService(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var newService Service
	err := json.NewDecoder(r.Body).Decode(&newService)
	if err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	mu.Lock()
	fmt.Println(newService.ID)
	services[newService.ID] = newService
	mu.Unlock()

	log.Printf("Service %s registered\n", newService.Name)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newService)
}

// Discover service by ID
func discoverService(w http.ResponseWriter, r *http.Request) {
	serviceID := r.URL.Path[len("/eureka/apps/"):]

	mu.Lock()
	service, exists := services[serviceID]
	// fmt.Println(services["user"])
	// service, exists := services["user"]
	// service, exists := services["user"]
	mu.Unlock()
	fmt.Println(service)

	if !exists {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(service)
}

// Deregister a service from the mock Eureka server
func deregisterService(w http.ResponseWriter, r *http.Request) {
	fmt.Println("calling me")
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var serviceToRemove Service
	err := json.NewDecoder(r.Body).Decode(&serviceToRemove)
	if err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	mu.Lock()
	delete(services, serviceToRemove.ID)
	mu.Unlock()

	log.Printf("Service %s deregistered\n", serviceToRemove.Name)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serviceToRemove)
}

func logRegisteredServices(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if len(services) == 0 {
		log.Println("No registered services.")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("No registered services.\n"))
		return
	}

	log.Println("Currently Registered Services:")
	for _, service := range services {
		log.Printf("ID=%s, Name=%s, Host=%s, Port=%d, Health=%s, URL=%s\n",
			service.ID, service.Name, service.Host, service.Port, service.Health, service.Metadata.URL)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(services)
}
