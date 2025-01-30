package helpers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"forum/internal/db"
)

func IsValidInput(user db.User) string {
	/*********** Name *************/
	nameRegex := `^[A-Za-z0-9._-]{2,15}$`
	NameR := regexp.MustCompile(nameRegex)
	if !NameR.MatchString(user.UserName) {
		return "invalid name !"
	} else if len(user.UserName) > 15 {
		return "too long name 🥲"
	}

	/*********** Email *************/
	emailRegex := `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	EmailR := regexp.MustCompile(emailRegex)
	if !EmailR.MatchString(user.Email) {
		return "invalid email !!!!"
	} else if len(user.Email) > 30 {
		return "too long email 🥲"
	}

	/*********** Password *************/

	if !isValidPassword(user.Password) || len(user.Password) < 8 {
		return "Your password needs to have a minimum of 8 characters, including an uppercase letter, a lowercase letter, a number, and a special character 😅"
	} else if len(user.Password) > 15 {
		return "too long password 🥲"
	}
	return ""
}

func isValidPassword(password string) bool {
	Upper := false
	Lower := false
	Number := false
	SpecialCar := false
	carSpe := "@*-_+.;:,$#~&^£!?|%/\\{}()"

	for _, char := range password {
		if unicode.IsUpper(char) {
			Upper = true
		}
		if unicode.IsLower(char) {
			Lower = true
		}
		if unicode.IsDigit(char) {
			Number = true
		}
		if strings.Contains(carSpe, string(char)) {
			SpecialCar = true
		}
	}

	if Upper && Lower && Number && SpecialCar {
		return true
	}
	return false
}

func IsValidPost(title, content string) *db.Post {
	post := &db.Post{
		Errors: make(map[string]string),
	}

	switch {
	case title == "":
		post.Errors["EmptyTitle"] = "Please enter a title!!"
	case content == "":
		post.Errors["EmptyContent"] = "Please enter a content!!"
	}

	return post
}

func JsonResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}