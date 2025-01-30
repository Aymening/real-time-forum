package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) AccessMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			response := CommentResponse{Message: "Unauthorized: Please log in to continue."}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}

		if _, err := h.DB.TokenVerification(cookie.Value); err != nil {
			http.Redirect(w, r, "/logout", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) RedirectMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("session")
		if err == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
