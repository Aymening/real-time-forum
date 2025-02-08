package handlers

import (
	"encoding/json"
	"net/http"

	"forum/helpers"
	"forum/internal/db"
)

/*********************   Register *********************/

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		helpers.ExecuteTmpl(w, "register.html", 200, "", nil)
		return

	} else if r.Method == http.MethodPost {
		var user db.User
		err := json.NewDecoder(r.Body).Decode(&user)
		defer r.Body.Close()

		if err != nil {
			helpers.JsonResponse(w, http.StatusBadRequest, "Bad Request 🫤")
			return
		}

		existUser, existEmail, err := h.DB.DatabaseVerification(user.UserName, user.Email)
		if err != nil {
			helpers.JsonResponse(w, http.StatusInternalServerError, "Internal server error 😥")
			return
		}

		if existUser {
			helpers.JsonResponse(w, http.StatusConflict, "This username already exists 🧐")
			return
		}
		if existEmail {
			helpers.JsonResponse(w, http.StatusConflict, "This email already exists 🧐")
			return
		}

		if helpers.IsValidInput(user) == "" {
			err = h.DB.InsertUser(h.Api.Users, user.UserName, user.Email, user.Password)
			if err != nil {
				helpers.JsonResponse(w, http.StatusInternalServerError, "Internal server error 😥")
				return
			}
		} else {
			helpers.JsonResponse(w, http.StatusConflict, helpers.IsValidInput(user))
			return
		}
		
	} else {
		helpers.JsonResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed 😥")
		return
	}
}