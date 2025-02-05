package main

import (
	"fmt"
	"forum/internal/web/server"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	app, err := server.InitApp()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("app: ", app)

	// Start a goroutine to handle WebSocket connections separately
	go func() {
		// Set up a WebSocket server (use the correct handler from your app)
		http.HandleFunc("/ws/chat", app.Handlers.WebSocketHandler)
		log.Fatal(http.ListenAndServe(":8229", nil)) // Use a different port for WebSocket if needed
	}()

	// Now set up the regular HTTP server
	serveur := http.Server{
		Addr:    ":8228", // This can remain the same, for your normal HTTP requests
		Handler: app.Routes(),
	}

	log.Println("\u001b[38;2;255;165;0mListing on http://localhost:8228...\033[0m")
	log.Fatal(serveur.ListenAndServe())
}

// func main() {
// 	app, err := server.InitApp()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("app: ", app)

// 	serveur := http.Server{
// 		Addr:    ":8228",
// 		Handler: app.Routes(),
// 	}

// 	log.Println("\u001b[38;2;255;165;0mListing on http://localhost:8228...\033[0m")
// 	log.Fatal(serveur.ListenAndServe())
// }
