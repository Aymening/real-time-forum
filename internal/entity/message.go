package entity

type Message struct {
	SenderID   int    `json:"sender"`
	ReceiverID int    `json:"receiver"`
	Content    string `json:"content"`
}
