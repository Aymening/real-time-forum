package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Initialize Database
	InitDB()

	// WebSocket route
	http.HandleFunc("/ws", handleConnections)

	// Serve frontend files
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	// Start server
	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
