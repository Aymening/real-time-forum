package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"forum/helpers"
	"forum/internal/db"
)

func (h *Handler) PostReactions(w http.ResponseWriter, r *http.Request) {
	loggedUser, err := helpers.CheckCookie(r, h.DB)
	if err != nil {
		helpers.JsonResponse(w, http.StatusUnauthorized, "Unauthorized: Please log in to continue.")
		return
	}

	if r.Method != http.MethodPost {
		helpers.JsonResponse(w, http.StatusMethodNotAllowed, "Oops! Method Not Allowed.")
		return
	}

	var userReaction db.UserReaction

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&userReaction); err != nil {
		helpers.JsonResponse(w, http.StatusBadRequest, "Bad Request:Error decoding JSON.")
		return
	}

	postId, err := strconv.Atoi(userReaction.PostId)
	if err != nil {
		helpers.JsonResponse(w, http.StatusBadRequest, "Invalid: postId is not correct.")
		return
	}

	lastReaction, err := h.DB.GetPostReactionFromDB(loggedUser, userReaction)
	if err != nil && err != sql.ErrNoRows {
		helpers.JsonResponse(w, http.StatusInternalServerError, "Oops! Internal Server Error.")
		return
	}

	if lastReaction == "" {
		if err := h.DB.InsertPostReactionInDB(loggedUser, userReaction); err != nil {
			helpers.JsonResponse(w, http.StatusInternalServerError, "Oops! Internal Server Error.")
		    return
		}
	} else {
		if strings.HasPrefix(userReaction.Reaction, "un") {
			if err := h.DB.DeletePostReactionInDB(loggedUser, userReaction); err != nil {
				helpers.JsonResponse(w, http.StatusInternalServerError, "Oops! Internal Server Error.")
		        return
			}
		} else {
			if err := h.DB.UpdatePostReactionInDB(loggedUser, userReaction); err != nil {
				helpers.JsonResponse(w, http.StatusInternalServerError, "Oops! Internal Server Error.")
		        return
			}
		}
		h.Api.DeletePostReactionInApi(loggedUser, postId, userReaction)
	}

	if !strings.HasPrefix(userReaction.Reaction, "un") {
		h.Api.AddPostReactionInApi(loggedUser, postId, userReaction)
	}
}
