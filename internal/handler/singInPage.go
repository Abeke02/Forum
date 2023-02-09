package handler

import (
	"errors"
	"net/http"

	"main.go/internal/service"
	"main.go/models"
)

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/signin" {
		h.errorPage(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	user := h.userIdentity(w, r)
	if user != (models.User{}) {
		h.errorPage(w, r, http.StatusOK, "you are already logged in")
		return
	}
	switch r.Method {
	case http.MethodGet:
		if err := h.Tmpl.ExecuteTemplate(w, "singin.html", nil); err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		username, ok := r.Form["username"]
		if !ok {
			h.errorPage(w, r, http.StatusBadRequest, "username field not found")
			return
		}
		password, ok := r.Form["password"]
		if !ok {
			h.errorPage(w, r, http.StatusBadRequest, "password field not found")
			return
		}
		sessionToken, expiresAt, err := h.Services.Auth.GenerateSessionToken(username[0], password[0])
		if err != nil {
			if errors.Is(err, service.ErrUserNotFound) {
				h.errorPage(w, r, http.StatusBadRequest, err.Error())
				return
			}
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: expiresAt,
			Path:    "/",
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}
