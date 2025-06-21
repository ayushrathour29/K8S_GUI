package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"k8_gui/internal/k8s"
	"k8_gui/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	clientset, metricsClient, err := k8s.InitK8sClient()
	if err != nil {
		log.Printf("Warning: Error initializing Kubernetes client: %v", err)
		log.Println("Starting server with limited functionality (auth endpoints will still work)")
		clientset = nil
		metricsClient = nil
	} else {
		log.Println("Kubernetes client initialized successfully")
	}

	router := server.NewRouter(clientset, metricsClient)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
