package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/websocket"
    _ "github.com/mattn/go-sqlite3"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type Message struct {
    Sender   string `json:"sender"`
    Receiver string `json:"receiver"`
    Content  string `json:"content"`
}

var db *sql.DB

func initDB() {
    var err error
    db, err = sql.Open("sqlite3", "./database.db")
    if err != nil {
        log.Fatal(err)
    }
    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        sender TEXT NOT NULL,
        receiver TEXT NOT NULL,
        content TEXT NOT NULL,
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
    )`)

    if err != nil {
        log.Fatal(err)
    }
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    defer ws.Close()

    clients[ws] = true

    for {
        var msg Message
        err := ws.ReadJSON(&msg)
        if err != nil {
            delete(clients, ws)
            break
        }

        broadcast <- msg
    }
}

func handleMessages() {
    for {
        msg := <-broadcast
        for client := range clients {
            err := client.WriteJSON(msg)
            if err != nil {
                log.Printf("Error: %v", err)
                client.Close()
                delete(clients, client)
            }
        }

        // Save message to SQLite
        stmt, _ := db.Prepare("INSERT INTO messages(sender, receiver, content) VALUES(?, ?, ?)")
        stmt.Exec(msg.Sender, msg.Receiver, msg.Content)
    }
}

func getConversation(w http.ResponseWriter, r *http.Request) {
    sender := r.URL.Query().Get("sender")
    receiver := r.URL.Query().Get("receiver")

    rows, err := db.Query("SELECT sender, receiver, content FROM messages WHERE (sender = ? AND receiver = ?) OR (sender = ? AND receiver = ?)", sender, receiver, receiver, sender)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var messages []Message
    for rows.Next() {
        var msg Message
        err := rows.Scan(&msg.Sender, &msg.Receiver, &msg.Content)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        messages = append(messages, msg)
    }

    json.NewEncoder(w).Encode(messages)
}

func main() {
    initDB()
    defer db.Close()

    go handleMessages()

    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/", fs)
    http.HandleFunc("/ws", handleConnections)
    http.HandleFunc("/get-conversation", getConversation)

    fmt.Println("Server started on :8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}