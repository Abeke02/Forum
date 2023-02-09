package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"main.go/models"
)

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.errorPage(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	user := h.userIdentity(w, r)
	if err := r.ParseForm(); err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error()) //
		return
	}
	filter, _ := r.Form["filter"]
	if len(filter) == 0 {
		filter = append(filter, "")
	}
	if filter[0] != "More Liked" && filter[0] != "Newest" && filter[0] != "" && filter[0] != "More Disliked" {
		h.errorPage(w, r, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	posts, err := h.Services.Post.GetAllPosts(filter[0])
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	}
	info := models.Info{
		ThatUser: user,
		Posts:    posts,
	}
	if err := h.Tmpl.ExecuteTemplate(w, "index.html", info); err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
	}
}
