package server

import (
	"encoding/json"
	"fmt"
	"forum/internal/api"
	"forum/internal/db"
	"forum/internal/entity"
	"forum/internal/web/handlers"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type App struct {
	Handlers handlers.Handler
	Api      api.Api
	DB       *db.Database
}

var mutex sync.Mutex
var clients = make(map[*websocket.Conn]int)
var broadcast = make(chan entity.Message)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewApp() (*App, error) {
	database, err := db.NewDatabase()
	if err != nil {
		return nil, err
	}

	if err := database.ExecuteAllTableInDataBase(); err != nil {
		return nil, err
	}

	return &App{
		Handlers: handlers.Handler{
			DB: database,
		},
		Api: api.Api{
			Users:    make(map[string]*db.User),
			Comments: make(map[int][]*db.Comment),
		},
		DB: database,
	}, nil
}

func InitApp() (*App, error) {

	// mux := http.NewServeMux()
	// mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	app, err := NewApp()
	if err != nil {
		return nil, err
	}

	if err := app.DB.GetPostsFromDB(app.Api.Users, &app.Api.Posts); err != nil {
		return nil, err
	}

	if err := app.DB.GetAllCommentsFromDataBase(app.Api.Comments); err != nil {
		return nil, err
	}

	app.Handlers.Api = &app.Api

	return app, nil
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./assets/templates/chat.html")
}

func (app *App) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := app.DB.GetAllUsersExceptCurrent(r)
	if err != nil {
		fmt.Print(err)
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// getchat **********************************>
func (app *App) GetChat(w http.ResponseWriter, r *http.Request) {

	// Get session token
	token, err := db.GetSessionToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get sender ID
	senderID, err := app.DB.GetSenderID(token)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	// Get receiver ID
	receiverID, err := db.GetReceiverID(r)
	if err != nil {
		http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
		return
	}

	// Fetch messages
	messages, err := app.DB.FetchMessages(senderID, receiverID)
	if err != nil {
		// fmt.Println(err)
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// =========================================================================>

// func (app *App) HandleConnections(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("connected")
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	defer func() {
// 		fmt.Println("closed")
// 		ws.Close()
// 	}()
// 	token, err := db.GetSessionToken(r)
// 	if err != nil {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	// Get sender ID
// 	senderID, err := app.DB.GetSenderID(token)
// 	if err != nil {
// 		http.Error(w, "Invalid session", http.StatusUnauthorized)
// 		return
// 	}

// 	clients[ws] = senderID
// 	for {
// 		var msg entity.Message
// 		err := ws.ReadJSON(&msg)

// 		if err != nil {
// 			fmt.Println(err)
// 			delete(clients, ws)
// 			break
// 		}
// 		msg.SenderID = senderID
// 		fmt.Println(senderID)

// 		// // here save messages in database function
// 		app.DB.SaveMessage(msg)
// 		ws.WriteJSON(msg)
// 		for conn, id := range clients {
// 			if id == msg.ReceiverID {
// 				conn.WriteJSON(msg)
// 				break
// 			}

// 		}

// 	}
// }

func (app *App) HandleConnections(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New WebSocket connection established")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading WebSocket:", err)
		return
	}
	defer func() {
		fmt.Println("WebSocket closed")
		ws.Close()
	}()

	// Get session token
	token, err := db.GetSessionToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get sender ID
	senderID, err := app.DB.GetSenderID(token)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	senderUsername, _ := app.DB.GetSenderUsernameByID(senderID)

	// Store user connection
	mutex.Lock()
	clients[ws] = senderID // Keep the original map structure
	mutex.Unlock()

	// fmt.Println("User", senderID, "connected")

	for {
		var msg entity.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			// fmt.Println("Error reading message:", err)
			mutex.Lock()
			delete(clients, ws) // Remove client on disconnect
			mutex.Unlock()
			break
		}

		// Assign the sender ID
		msg.SenderID = senderID
		// fmt.Println("Message from", senderID, "to", msg.ReceiverID, ":", msg.Content)

		// Save message in database
		app.DB.SaveMessage(msg)

		// Find the recipient's WebSocket connection
		var recipientConn *websocket.Conn
		mutex.Lock()
		for conn, id := range clients {
			if id == msg.ReceiverID {
				recipientConn = conn
				break
			}
		}
		mutex.Unlock()

		// If recipient is online, send message and notification
		if recipientConn != nil {
			notification := map[string]string{
				"type":    "notification",
				"from":    senderUsername,
				"message": "You have a new message!",
			}
			recipientConn.WriteJSON(notification) // Notify receiver
			recipientConn.WriteJSON(msg)
			// Send message to receiver
			fmt.Println("A message was saved")
		}
		ws.WriteJSON(msg)
		// else {
		// 	app.DB.SaveMessage(msg)
		// }
	}
}
