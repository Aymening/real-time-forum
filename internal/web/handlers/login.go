package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"forum/helpers"
	"forum/internal/db"

	"github.com/google/uuid"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		helpers.ExecuteTmpl(w, "login.html", 200, "", nil)
		return
	} else if r.Method == http.MethodPost {
		var user db.User
		err := json.NewDecoder(r.Body).Decode(&user)
		defer r.Body.Close()

		if err != nil {
			helpers.JsonResponse(w, http.StatusBadRequest, "Bad Request 🫤")
			return
		}

		id, err := h.DB.Authenticate(
			user.Email, user.Password,
		)
		if err != nil {
			if err == sql.ErrNoRows || strings.Contains(err.Error(), "hashedPassword") {
				helpers.JsonResponse(w, http.StatusConflict, "UserName or Password is invalid 🧐")
				return
			} else {
				helpers.JsonResponse(w, http.StatusInternalServerError, "Internal server error 😥")
				return
			}
		}

		token := uuid.New().String()

		if err := h.DB.InsertToken(id, token); err != nil {
			helpers.JsonResponse(w, http.StatusInternalServerError, "Internal server error 😥")
			return
		}

		helpers.AddCookie(w, token)

		// w.WriteHeader(302) // 302 Found
	} else {
		helpers.JsonResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed 😥")
		return
	}
}