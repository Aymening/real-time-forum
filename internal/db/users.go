package db

import (
	"fmt"
	"forum/internal/entity"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int64             `json:"id"`
	UserName  string            `json:"name"`
	Email     string            `json:"email,omitempty"`
	Password  string            `json:"password,omitempty"`
	Token     string            `json:"-"`
	Posts     []*Post           `json:"posts,omitempty"`
	Errors    map[string]string `json:"-,omitempty"`
	Reactions map[string][]int  `json:"reactions,omitempty"`
}

func (d *Database) GetAllUsersExceptCurrent(r *http.Request) ([]entity.Contact, error) {
	// Get the session token from the request
	token, err := GetSessionToken(r)
	if err != nil {
		return nil, err
	}

	// Retrieve the current user's ID based on the session token
	var currentUserID int
	err = d.db.QueryRow("SELECT id FROM users WHERE token = ?", token).Scan(&currentUserID)
	if err != nil {
		return nil, err
	}

	// Query to get all users except the current user
	rows, err := d.db.Query(`SELECT 
    u.id AS user_id,
    u.username,
    MAX(m.timestamp) AS last_message_time
FROM 
    users u
LEFT JOIN 
    messages m 
ON 
    (m.sender = u.id OR m.receiver = u.id)
    AND (m.sender = ? OR m.receiver = ?)
WHERE 
    u.id != ?
GROUP BY 
    u.id, u.username
ORDER BY 
    last_message_time DESC;`, currentUserID, currentUserID, currentUserID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var users []entity.Contact
	for rows.Next() {
		var user entity.Contact
		if err := rows.Scan(&user.Id, &user.UserName, &user.LastMessage); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// func (d *Database) GetAllUsers() ([]User, error) {
// 	rows, err := d.db.Query("SELECT id, username FROM users")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var users []User
// 	for rows.Next() {
// 		var user User
// 		if err := rows.Scan(&user.Id, &user.UserName); err != nil {
// 			return nil, err
// 		}
// 		users = append(users, user)
// 	}

// 	return users, nil
// }

func (d *Database) InsertUser(users map[string]*User, name, email, Password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(Password), 12) // 2¹² times
	if err != nil {
		return err
	}

	expression := `INSERT INTO users (username, email, password)
	VALUES (?, ?, ?)`

	stmnt, err := d.db.Prepare(expression)
	if err != nil {
		return err
	}

	defer stmnt.Close()

	row, err := stmnt.Exec(name, email, passwordHash)
	if err != nil {
		return err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return err
	}

	users[name] = &User{
		Id:       id,
		UserName: name,
	}

	return nil
}

func (d *Database) Authenticate(email, Password string) (int, error) {
	var id int
	var passwordHash []byte

	expression := `SELECT id, password From users WHERE email = ?`
	row := d.db.QueryRow(expression, email)
	err := row.Scan(&id, &passwordHash)
	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(passwordHash, []byte(Password))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (d *Database) InsertToken(id int, token string) error {

	expression := `UPDATE users SET Token = ? WHERE id = ?;`
	_, err := d.db.Exec(expression, token, id)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) DatabaseVerification(name, email string) (bool, bool, error) {

	// Username Verification
	usernameExists := false
	usernameExpression := `SELECT EXISTS (SELECT * FROM users WHERE username LIKE ?);`
	err := d.db.QueryRow(usernameExpression, name).Scan(&usernameExists)
	if err != nil {
		return false, false, err
	}

	// email verification
	emailExists := false
	emailExpression := `SELECT EXISTS (SELECT * FROM users WHERE email LIKE ?);`
	err = d.db.QueryRow(emailExpression, email).Scan(&emailExists)
	if err != nil {
		return false, false, err
	}

	return usernameExists, emailExists, nil
}

func (d *Database) TokenVerification(token string) (string, error) {
	var user User

	expression := `SELECT username FROM users WHERE token = ? `

	row := d.db.QueryRow(expression, token)

	err := row.Scan(&user.UserName)
	if err != nil {
		return "", err
	}

	return user.UserName, nil
}

// retrieve username of sender from its id
func (db *Database) GetSenderUsernameByID(senderID int) (string, error) {
	var username string
	query := `SELECT username FROM users WHERE id = ?`
	err := db.db.QueryRow(query, senderID).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}
