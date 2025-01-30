package handlers

import (
	"forum/helpers"
	"net/http"
)

/********************** Logout ********************/
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.ExecuteTmpl(w, "error.html", 405, "Oops! Method not allowed.", nil)
		return
	}
	helpers.DeleteCookie(w)
	http.Redirect(w, r, "/login", http.StatusFound)
}
