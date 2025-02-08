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
		helpers.JsonResponse(w, http.StatusUnauthorized, "Unauthorized: Please log in to continue.")
		return
	}

	if r.Method != http.MethodPost {
		helpers.JsonResponse(w, http.StatusMethodNotAllowed, "Oops! Method Not Allowed.")
		return
	}

	var requestBodyOfComment RequestBodyComment

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&requestBodyOfComment); err != nil {
		helpers.JsonResponse(w, http.StatusBadRequest, "Bad Request: Error decoding JSON.")
		return
	}

	if requestBodyOfComment.Content == "" || requestBodyOfComment.IdPost <= 0 {
		helpers.JsonResponse(w, http.StatusBadRequest, "Invalid input: Content and postId are required.")
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
		helpers.JsonResponse(w, http.StatusInternalServerError, "Oops! Internal Server Error.")
		return
	}

	h.Api.Comments[requestBodyOfComment.IdPost] = append(h.Api.Comments[requestBodyOfComment.IdPost], comments)

	helpers.JsonResponse(w, http.StatusOK, "Comment created successfully!")
}
