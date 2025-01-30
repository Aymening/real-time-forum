package handlers

import (
	"encoding/json"
	"forum/helpers"
	"log"
	"net/http"
	"time"
)

type RequestBodyComment struct {
	IdPost  int    `json:"postId"`
	Content string `json:"content"`
}

type CommentResponse struct {
	Message string `json:"message"`
}

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	
	loggedUser, err := helpers.CheckCookie(r, h.DB)
	if err != nil {
		response := CommentResponse{Message: "Unauthorized: Please log in to continue."}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	if r.Method != http.MethodPost {
		response := CommentResponse{Message: "Oops! Method Not Allowed."}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	var requestBodyOfComment RequestBodyComment

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&requestBodyOfComment); err != nil {
		response := CommentResponse{Message: "Bad Request: Error decoding JSON."}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if requestBodyOfComment.Content == "" || requestBodyOfComment.IdPost <= 0 {
		response := CommentResponse{Message: "Invalid input: Content and postId are required."}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	comments, err := h.DB.InsertComment(
		requestBodyOfComment.Content,
		time.Now().Format("Mon Jan 02 15:04:05 2006"),
		loggedUser, h.Api.Users,
		requestBodyOfComment.IdPost,
	)

	if err != nil {
		log.Printf("Error inserting comment: %v", err)
		response := CommentResponse{Message: "Oops! Internal Server Error."}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	h.Api.Comments[requestBodyOfComment.IdPost] = append(h.Api.Comments[requestBodyOfComment.IdPost], comments)

	response := CommentResponse{Message: "Comment created successfully!"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
