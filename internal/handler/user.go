package handler

import (
	"net/http"
	"strings"

	"main.go/models"
)

func (h *Handler) userProfilePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	username := strings.TrimPrefix(r.URL.Path, "/profile/")

	userI := h.userIdentity(w, r)

	if userI.Username == username {

		postsUser, err := h.Services.Post.GetPostsByUsername(username)
		if err != nil {
			h.errorPage(w, r, http.StatusBadRequest, err.Error())
			return
		}

		postLiked := []models.Post{}

		ReactionPost, err := h.Services.GetIDPostsByUsername(username)
		// fmt.Println(ReactionPost)
		for i := 0; i < len(ReactionPost); i++ {
			post, err := h.Services.GetPostById(ReactionPost[i].PostId)
			// fmt.Println(ReactionPost[i].PostId)
			if err != nil {
				h.errorPage(w, r, http.StatusInternalServerError, err.Error())
				return
			}
			postLiked = append(postLiked, post)
		}
		// post:=h.Services.GetPostById(ReactionPost...)

		info := models.Info{
			ThatUser:  userI,
			Posts:     postsUser,
			PostsLike: postLiked,
		}

		if err := h.Tmpl.ExecuteTemplate(w, "user.html", info); err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}

	} else {

		user, err := h.Services.User.GetUserByUsername(username)
		if err != nil {
			h.errorPage(w, r, http.StatusNotFound, err.Error())
			return
		}
		postsUser, err := h.Services.Post.GetPostsByUsername(username)
		if err != nil {
			h.errorPage(w, r, http.StatusBadRequest, err.Error())
			return
		}

		info := models.Info{
			User:     user,
			Posts:    postsUser,
			ThatUser: userI,
		}
		if err := h.Tmpl.ExecuteTemplate(w, "user.html", info); err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}

	}
}
