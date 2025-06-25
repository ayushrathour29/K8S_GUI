package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"k8_gui/internal/k8s"
	"k8_gui/internal/server"
	"k8_gui/internal/utils"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(utils.LogNoEnvFile)
	}

	clientset, metricsClient, err := k8s.InitK8sClient()
	if err != nil {
		log.Printf(utils.LogWarnInitK8sClient, err)
		log.Println(utils.LogLimitedServer)
		clientset = nil
		metricsClient = nil
	} else {
		log.Println(utils.LogK8sClientInitSuccess)
	}

	router := server.NewRouter(clientset, metricsClient)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
