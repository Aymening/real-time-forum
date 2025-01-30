package handlers

import (
	"net/http"

	"forum/helpers"
)

func (h *Handler) Static(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.ExecuteTmpl(w, "error.html", 405, "Oops! Method Not Allowed!", nil)
		return
	}

	allowedFiles := map[string]bool{
		"css/home.css":           true,
		"css/post.css":           true,
		"css/styles.css":         true,
		"css/comments.css":       true,
		"css/error.css":          true,
		"js/register.js":         true,
		"js/login.js":            true,
		"js/permissionDenied.js": true,
		"js/displayPosts.js":     true,
		"js/createPost.js":       true,
		"js/createComment.js":    true,
		"js/reaction.js":         true,
	}
	path := r.URL.Path[len("/static/"):]
	if allowedFiles[path] {
		http.StripPrefix("/static/", http.FileServer(http.Dir("./assets/static/"))).ServeHTTP(w, r)
	} else {
		helpers.ExecuteTmpl(w, "error.html", 403, "Oops! Access Forbidden", nil)
		return
	}
}
