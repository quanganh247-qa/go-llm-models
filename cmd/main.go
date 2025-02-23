package main

import (
	"log"
	"vet-tails/ai/internal/router"
)

func main() {
	// Setup router
	router := router.SetupRouter()

	// Start server
	log.Fatal(router.Run(":8080"))
}
