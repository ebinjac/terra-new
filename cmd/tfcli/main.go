package main

import (
	"fmt"
	"log"
	"os"

	"github.com/terraform-cli/pkg/server"
)

func main() {
	port := 8080
	if envPort := os.Getenv("PORT"); envPort != "" {
		fmt.Sscanf(envPort, "%d", &port)
	}

	srv := server.NewServer(port)
	log.Printf("Starting server on port %d...\n", port)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}
