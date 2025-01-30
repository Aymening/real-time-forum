package handlers

import (
	"net/http"

	"forum/helpers"
	"forum/internal/api"
	"forum/internal/db"
)

type Handler struct {
	Api *api.Api
	DB  *db.Database
}

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	
	if r.URL.Path != "/" {
		helpers.ExecuteTmpl(w, "error.html", 404, "Oops! Not found", nil)
		return
	}

	userName := ""
	session, _ := r.Cookie("session")
	if session != nil {
		tem, err := h.DB.TokenVerification(session.Value)
		userName = tem
		if err != nil {
			helpers.DeleteCookie(w)
			helpers.ExecuteTmpl(w, "home.html", 200, "", nil)
			return
		} else {
			// Update
			helpers.ExecuteTmpl(w, "home.html", 200, "", h.Api.Users[userName])
			return
		}
	} else {
		helpers.ExecuteTmpl(w, "home.html", 200, "", nil)
		return
	}
}
