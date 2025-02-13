package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"forum/internal/entity"
	"forum/internal/web/server"

	_ "github.com/mattn/go-sqlite3"
)

// type Message struct {
// 	SenderID   int    `json:"sender"`
// 	ReceiverID int    `json:"receiver"`
// 	Content    string `json:"content"`
// }

var db *sql.DB

// func initDB() {
// 	var err error
// 	db, err = sql.Open("sqlite3", "./internal/db/data.db")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// _, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
//     //     id INTEGER PRIMARY KEY AUTOINCREMENT,
//     //     sender TEXT NOT NULL,
//     //     receiver TEXT NOT NULL,
//     //     content TEXT NOT NULL,
//     //     timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
//     // )`)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func getConversation(w http.ResponseWriter, r *http.Request) {
	sender := r.URL.Query().Get("sender")
	receiver := (r.URL.Query)().Get("receiver")

	rows, err := db.Query("SELECT sender, receiver, content FROM messages WHERE (sender = ? AND receiver = ?) OR (sender = ? AND receiver = ?)", sender, receiver, receiver, sender)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []entity.Message
	for rows.Next() {
		var msg entity.Message
		err := rows.Scan(&msg.SenderID, &msg.ReceiverID, &msg.Content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		messages = append(messages, msg)
	}

	json.NewEncoder(w).Encode(messages)
}

func main() {

	// initDB()
	defer db.Close()

	app, err := server.InitApp()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/get-conversation", getConversation)

	mux.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/chat.html")
	})

	serveur := http.Server{
		Addr:    ":8228",
		Handler: app.Routes(),
	}

	log.Println("\u001b[38;2;255;165;0mListing on http://localhost:8228...\033[0m")
	log.Fatal(serveur.ListenAndServe())
}
