package entity

type Contact struct {
	Id          int `json:"id"`
	UserName    string `json:"name"`
	LastMessage string
}
