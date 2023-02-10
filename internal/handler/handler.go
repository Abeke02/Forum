package handler

import (
	"html/template"
	"net/http"

	"main.go/internal/service"
	"main.go/models"
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
	h.Mux.HandleFunc("/post/category/", h.filterPostCategory)

	h.Mux.HandleFunc("/comment/like/", h.likeComment)
	h.Mux.HandleFunc("/comment/dislike/", h.dislikeComment)

	h.Mux.HandleFunc("/profile/", h.userProfilePage)
	// h.Mux.HandleFunc("/likeposts/", h.userLikesPosts)

	// h.Mux.Handle("/ui/", http.StripPrefix("/ui", http.FileServer(http.Dir("./ui"))))
	h.Mux.Handle("/ui/static/", http.StripPrefix("/ui/static/", http.FileServer(http.Dir("./ui/static"))))
	h.Mux.Handle("/ui/image/", http.StripPrefix("/ui/image/", http.FileServer(http.Dir("./ui/image"))))

	return h.Mux
}

func (h *Handler) filterPostCategory(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/category/" {
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

	// categories, err := h.Services.Post.GetAllCategories()
	// if err != nil {
	// 	h.errorPage(w, r, http.StatusInternalServerError, err.Error())
	// }

	cat := (r.URL.Query().Get("cat"))
	catID, err := h.Services.GetIdPostsByCategory(cat)
	if len(catID) == 0 {
		h.errorPage(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if err != nil {
		h.errorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}
	var posts []models.Post
	for _, i := range catID {
		post, err := h.Services.GetPostById(i)
		if err != nil {
			h.errorPage(w, r, http.StatusNotFound, err.Error())
			return
		}
		posts = append(posts, post)
	}

	info := models.Info{
		Posts:    posts,
		ThatUser: user,
		Category: cat,
	}
	if err := h.Tmpl.ExecuteTemplate(w, "category.html", info); err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}
