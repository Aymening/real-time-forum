package handlers

import (
	"net/http"

	"forum/helpers"
)

func (h *Handler) UserPosts(w http.ResponseWriter, r *http.Request) {
	loggedUser, err := helpers.CheckCookie(r, h.DB)
	if err != nil {
		helpers.ExecuteTmpl(w, "error.html", http.StatusInternalServerError, "Oops! Internal server error.", nil)
		return
	}

	helpers.ExecuteTmpl(w, "user_posts.html", http.StatusOK, "", h.Api.Users[loggedUser])
}
