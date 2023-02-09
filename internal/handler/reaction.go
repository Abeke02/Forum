package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"main.go/models"
)

func (h *Handler) likePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	postId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		h.errorPage(w, r, http.StatusNotFound, err.Error())
		return
	}
	_, err = h.Services.Post.GetPostById(postId)
	if err != nil {
		h.errorPage(w, r, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	user := h.userIdentity(w, r)
	if user == (models.User{}) {
		h.errorPage(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	if err := h.Services.Reaction.CreateReaction(postId, "post", "like", user.Username); err != nil {
		h.errorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/?id=%d", postId), http.StatusSeeOther)
}

func (h *Handler) dislikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	postId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		h.errorPage(w, r, http.StatusNotFound, err.Error())
		return
	}

	_, err = h.Services.Post.GetPostById(postId)
	if err != nil {
		h.errorPage(w, r, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	user := h.userIdentity(w, r)
	if user == (models.User{}) {
		h.errorPage(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	if err := h.Services.Reaction.CreateReaction(postId, "post", "dislike", user.Username); err != nil {
		h.errorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/?id=%d", postId), http.StatusSeeOther)
}

func (h *Handler) likeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	commentId, err := strconv.Atoi(r.URL.Query().Get("id_comment"))
	if err != nil {
		h.errorPage(w, r, http.StatusNotFound, err.Error())
		return
	}

	postId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		h.errorPage(w, r, http.StatusNotFound, err.Error())
		return
	}

	_, err = h.Services.Post.GetPostById(postId)
	if err != nil {
		h.errorPage(w, r, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	user := h.userIdentity(w, r)
	if user == (models.User{}) {
		h.errorPage(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	if err := h.Services.Reaction.CreateReaction(commentId, "comment", "like", user.Username); err != nil {
		h.errorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/?id=%d", postId), http.StatusSeeOther)
}

func (h *Handler) dislikeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	commentId, err := strconv.Atoi(r.URL.Query().Get("id_comment"))
	if err != nil {
		h.errorPage(w, r, http.StatusNotFound, err.Error())
		return
	}
	postId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		h.errorPage(w, r, http.StatusNotFound, err.Error())
		return
	}

	_, err = h.Services.Post.GetPostById(postId)
	if err != nil {
		h.errorPage(w, r, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	user := h.userIdentity(w, r)
	if user == (models.User{}) {
		h.errorPage(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	if err := h.Services.Reaction.CreateReaction(commentId, "comment", "dislike", user.Username); err != nil {
		h.errorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/?id=%d", postId), http.StatusSeeOther)
}
