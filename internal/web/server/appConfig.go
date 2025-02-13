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

	"github.com/gorilla/websocket"
)

type App struct {
	Handlers handlers.Handler
	Api      api.Api
	DB       *db.Database
}

var clients = make(map[*websocket.Conn]bool)
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
	users, err := app.DB.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// getchat **********************************>
func (app *App) GetChat(w http.ResponseWriter, r *http.Request) {
	fmt.Println("this ")
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
	fmt.Println("ddd")
	messages, err := app.DB.FetchMessages(senderID, receiverID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// =========================================================================>

func (app *App) HandleConnections(w http.ResponseWriter, r *http.Request) {
	fmt.Println("connected")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		fmt.Println("closed")
		ws.Close()
	}()
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

	clients[ws] = true
	for {
		var msg entity.Message
		err := ws.ReadJSON(&msg)

		if err != nil {
			fmt.Println(err)
			delete(clients, ws)
			break
		}
		msg.SenderID = senderID
		fmt.Println(senderID)
		// // here save messages in database function
		app.DB.SaveMessage(msg)

	}
}
