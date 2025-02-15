package entity

import "database/sql"

type Contact struct {
	Id          int    `json:"id"`
	UserName    string `json:"name"`
	LastMessage sql.NullString
}
