package db

import (
	"fmt"
	"forum/internal/entity"
	"time"
)

func (d *Database) SaveMessage(msg entity.Message) error {
	query := `INSERT INTO messages (sender, receiver, content, timestamp) 
	          VALUES (?, ?, ?, ?)`

	_, err := d.db.Exec(query, msg.SenderID, msg.ReceiverID, msg.Content, time.Now())
	if err != nil {
		// fmt.Println(err)
		return fmt.Errorf("could not save message: %v", err)
	}

	return nil
}
