package handler

import (
	"html/template"
	"net/http"

	"main.go/internal/service"
)

type Handler struct {
	Mux      *http.ServeMux
	Tmpl     *template.Template
	Services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		Mux:      http.NewServeMux(),
		Tmpl:     template.Must(template.ParseGlob("./ui/temp/*.html")),
		Services: services,
	}
}

func (h *Handler) InitRouter() *http.ServeMux {
	h.Mux.HandleFunc("/", h.homePage)
	h.Mux.HandleFunc("/auth/signup", h.signUp)
	h.Mux.HandleFunc("/auth/signin", h.signIn)
	h.Mux.HandleFunc("/auth/logout", h.logOut)

	h.Mux.HandleFunc("/post/", h.post)
	h.Mux.HandleFunc("/post/create", h.createPost)
	h.Mux.HandleFunc("/post/like/", h.likePost)
	h.Mux.HandleFunc("/post/dislike/", h.dislikePost)
	h.Mux.HandleFunc("/post/categories/", h.filterPostCategories)

	h.Mux.HandleFunc("/comment/like/", h.likeComment)
	h.Mux.HandleFunc("/comment/dislike/", h.dislikeComment)

	h.Mux.HandleFunc("/profile/", h.userProfilePage)
	// h.Mux.HandleFunc("/likeposts/", h.userLikesPosts)

	// h.Mux.Handle("/ui/", http.StripPrefix("/ui", http.FileServer(http.Dir("./ui"))))
	h.Mux.Handle("/ui/static/", http.StripPrefix("/ui/static/", http.FileServer(http.Dir("./ui/static"))))
	h.Mux.Handle("/ui/image/", http.StripPrefix("/ui/image/", http.FileServer(http.Dir("./ui/image"))))

	return h.Mux
}
