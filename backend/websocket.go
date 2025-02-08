package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]string) // WebSocket connections mapped to usernames
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Handles new WebSocket connections
func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade error:", err)
		return
	}
	defer conn.Close()

	var username string
	err = conn.ReadJSON(&username) // First message is the username
	if err != nil {
		log.Println("Username Read Error:", err)
		return
	}

	clients[conn] = username // Store the username

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			delete(clients, conn)
			break
		}

		// Save to database
		_, err = db.Exec("INSERT INTO messages (sender, receiver, message) VALUES (?, ?, ?)", msg.Sender, msg.Receiver, msg.Message)
		if err != nil {
			log.Println("DB error:", err)
			continue
		}

		// Send message to the receiver in real-time
		for client, name := range clients {
			if name == msg.Receiver { // Find the right receiver
				err = client.WriteJSON(msg)
				if err != nil {
					log.Println("Send error:", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}

// Message structure
type Message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}
