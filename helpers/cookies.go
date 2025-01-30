package helpers

import (
	"forum/internal/db"
	"net/http"
	"time"
)

func DeleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func AddCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func CheckCookie(r *http.Request, d *db.Database) (string, error) {
	session, err := r.Cookie("session")
	if err != nil {
		return "", err
	}

	userName, err := d.TokenVerification(session.Value)
	if err != nil {
		return "", err
	}
	return userName, nil
}
