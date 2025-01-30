package handlers

import (
	"encoding/json"
	"fmt"
	"forum/helpers"
	"forum/internal/db"
	"log"
	"net/http"
	"time"
)

type RequestBodyPost struct {
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Categories []string `json:"selectedCategories"`
}

type PostResponse struct {
	Message string `json:"message"`
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {

	loggedUser, err := helpers.CheckCookie(r, h.DB)
	if err != nil {
		response := PostResponse{Message: "Unauthorized: Please log in to continue."}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	if r.Method != http.MethodPost {
		response := PostResponse{Message: "Oops! Method Not Allowed."}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	var requestBodyOfPost RequestBodyPost

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&requestBodyOfPost); err != nil {
		fmt.Println(err)
		response := PostResponse{Message: "Bad Request: Error decoding JSON."}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(response)
		return
	}

	if post := helpers.IsValidPost(requestBodyOfPost.Title, requestBodyOfPost.Content); len(post.Errors) != 0 || len(requestBodyOfPost.Categories) == 0 {
		response := PostResponse{Message: "Invalid input: Content, title and selectedCategories are required."}
		fmt.Println("2")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	post, err := h.DB.InsertPostWithCategories(
		h.Api.Users,
		loggedUser,
		requestBodyOfPost.Title,
		requestBodyOfPost.Content, time.Now().Format("Mon Jan 02 15:04:05 2006"),
		requestBodyOfPost.Categories,
	)
	if err != nil {
		log.Printf("Error inserting post: %v", err)
		response := PostResponse{Message: "Oops! Internal Server Error."}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return

	}

	h.Api.Users[post.User].Posts = append([]*db.Post{post}, h.Api.Users[post.User].Posts...)
	h.Api.Posts = append([]*db.Post{post}, h.Api.Posts...)

	// Respond with a JSON message
	response := PostResponse{Message: "post created successfully!"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
