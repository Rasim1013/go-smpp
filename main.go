package main

import (
	"log"
	"net/http"
	"smpp-sender/handlers"
)

func main() {
	// http.HandleFunc("/", handlers.HomeHandlers)
	http.HandleFunc("/smpp-api/v1/sendone", handlers.SendOneHandler)
	http.HandleFunc("/smpp-api/v1/sendbulk", handlers.SendBulkHandler)

	log.Println("Server is running on port 1013")
	if err := http.ListenAndServe(":1013", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
