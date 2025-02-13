package db

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/entity"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("sqlite3", "./internal/db/data.db")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Database{
		db: db,
	}, nil
}

// ***************************************************************************************

// **Step 1: Get the token from cookies**
func GetSessionToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// **Step 2: Get the sender ID using the token**
func (db *Database) GetSenderID(token string) (int, error) {
	var senderID int
	query := `SELECT id FROM users WHERE token = ?`
	err := db.db.QueryRow(query, token).Scan(&senderID)
	if err != nil {
		return 0, err
	}
	return senderID, nil
}

// **Step 3: Get the receiver ID from query parameters**
func GetReceiverID(r *http.Request) (int, error) {
	receiverIDStr := r.URL.Query().Get("receiver")
	if receiverIDStr == "" {
		return 0, errors.New("missing receiver ID")
	}

	receiverID, err := strconv.Atoi(receiverIDStr)
	if err != nil {
		return 0, err
	}

	return receiverID, nil
}

// **Step 4: Fetch messages from the database**
func (db *Database) FetchMessages(senderID, receiverID int) ([]entity.Message, error) {
	fmt.Println(senderID, receiverID)
	var messages []entity.Message
	query := `SELECT sender, receiver, content FROM messages 
              WHERE (sender = ? AND receiver = ?) OR (sender = ? AND receiver = ?) 
              ORDER BY timestamp ASC`

	rows, err := db.db.Query(query, senderID, receiverID, receiverID, senderID)
	if err != nil {
		// fmt.Println(err)
		return nil, err
	}
	// fmt.Println(messages)
	defer rows.Close()

	for rows.Next() {
		var msg entity.Message
		err := rows.Scan(&msg.SenderID, &msg.ReceiverID, &msg.Content)
		if err != nil {
			return nil, err
		}
		fmt.Println(msg.Content)
		messages = append(messages, msg)
	}
	return messages, nil
}
