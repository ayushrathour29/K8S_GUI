package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"k8_gui/internal/k8s"
	"k8_gui/internal/server"
)

func main() {
	// Initialize Kubernetes client
	clientset, err := k8s.InitK8sClient()
	if err != nil {
		log.Fatalf("Error initializing Kubernetes client: %v", err)
	}

	// Creating router with clientset
	router := server.NewRouter(clientset)

	// Starting  HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
