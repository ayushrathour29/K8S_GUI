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
	clientset, metricsClient, err := k8s.InitK8sClient()
	if err != nil {
		log.Fatalf("Error initializing Kubernetes client: %v", err)
	}

	router := server.NewRouter(clientset, metricsClient)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
