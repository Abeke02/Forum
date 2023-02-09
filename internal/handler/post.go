package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"main.go/models"
)

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/" {
		h.errorPage(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		h.errorPage(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	post, err := h.Services.Post.GetPostById(id)
	if err != nil {
		h.errorPage(w, r, http.StatusNotFound, err.Error())
		return
	}
	comments, err := h.Services.GetCommentsByIdPost(id)
	if err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	for i := 0; i < len(comments); i++ {
		comments[i].Likes, comments[i].Dislikes, err = h.Services.Reaction.GetCounts(comments[i].Id, "comment")
		if err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		if err := h.Services.Post.UpdateCountsReactions("comment", comments[i].Likes, comments[i].Dislikes, comments[i].Id); err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	}
	user := h.userIdentity(w, r)
	if err := r.ParseForm(); err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	post.Likes, post.Dislikes, err = h.Services.Reaction.GetCounts(post.Id, "post")
	if err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.Services.Post.UpdateCountsReactions("post", post.Likes, post.Dislikes, post.Id); err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:

		info := models.Info{
			Post:     post,
			Comments: comments,
			ThatUser: user,
		}

		if err := h.Tmpl.ExecuteTemplate(w, "post.html", info); err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		user := h.userIdentity(w, r)
		if user == (models.User{}) {
			h.errorPage(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}
		comment, ok := r.Form["comment"]
		commentCheck := strings.Trim(comment[0], " ")
		if !ok || commentCheck == "" {
			h.errorPage(w, r, http.StatusBadRequest, "comment field not found")
			return
		}

		if comment[0] == "" {
			h.errorPage(w, r, http.StatusBadRequest, "comment field not found")
			return
		}

		if err := h.Services.Comment.CreateComment(id, user, comment[0]); err != nil {
			h.errorPage(w, r, http.StatusBadRequest, "comment field not found")
			return
		}

		http.Redirect(w, r, r.URL.Path+fmt.Sprintf("/?id=%d", id), http.StatusSeeOther)
	default:
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/create" {
		h.errorPage(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	user := h.userIdentity(w, r)
	if user == (models.User{}) {
		h.errorPage(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	if err := r.ParseForm(); err != nil {
		h.errorPage(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	info := models.Info{
		ThatUser: user,
	}
	switch r.Method {
	case http.MethodGet:
		if err := h.Tmpl.ExecuteTemplate(w, "createPost.html", info); err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		title, ok := r.Form["title"]
		if !ok || title[0] == "" {
			h.errorPage(w, r, http.StatusBadRequest, "title field not found")
			return
		}
		titl := strings.Trim(title[0], " ")
		if titl == "" {
			h.errorPage(w, r, http.StatusBadRequest, "title field not found")
			return
		}
		description, ok := r.Form["description"]

		if !ok || description[0] == "" {
			h.errorPage(w, r, http.StatusBadRequest, "description field not found")
			return
		}
		des := strings.Trim(description[0], " ")
		if des == "" {
			h.errorPage(w, r, http.StatusBadRequest, "description field not found")
			return
		}
		category, ok := r.Form["category"]

		if !ok || category[0] == "" {
			h.errorPage(w, r, http.StatusBadRequest, "category field not found")
			return
		}
		cat := strings.Trim(category[0], " ")
		if cat == "" {
			h.errorPage(w, r, http.StatusBadRequest, "category field not found")
			return
		}
		post := models.Post{
			Title:       title[0],
			Description: description[0],
			Category:    category,
		}
		user := h.userIdentity(w, r)

		if user == (models.User{}) {
			h.errorPage(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}
		if err := h.Services.Post.CreatePost(post, user); err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func (h *Handler) filterPostCategories(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/categories/" {
		h.errorPage(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	user := h.userIdentity(w, r)
	if err := r.ParseForm(); err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	categories, err := h.Services.Post.GetAllCategories()
	if err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
	}
	posts, err := h.Services.Post.GetAllPosts("")
	if err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	info := models.Info{
		Categories: categories,
		Posts:      posts,
		ThatUser:   user,
	}
	if err := h.Tmpl.ExecuteTemplate(w, "categories.html", info); err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}
