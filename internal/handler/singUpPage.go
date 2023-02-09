package handler

import (
	"errors"
	"net/http"
	"strings"

	"main.go/internal/service"
	"main.go/models"
)

type data struct {
	Message string
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/signup" {
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
		if err := h.Tmpl.ExecuteTemplate(w, "singUp.html", nil); err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		email, ok := r.Form["email"]
		emailCheck := strings.Trim(email[0], " ")
		if !ok || email[0] == "" || emailCheck == "" {
			h.errorPage(w, r, http.StatusBadRequest, "email field not found")
			return
			// msg := data{
			// 	Message: "email field not found",
			// }
			// if err := h.Tmpl.ExecuteTemplate(w, "singUp.html", msg); err != nil {
			// 	h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			// 	return
			// }
			// return
		}
		username, ok := r.Form["username"]
		usernameCheck := strings.Trim(username[0], " ")
		if !ok || username[0] == "" || usernameCheck == "" {
			h.errorPage(w, r, http.StatusBadRequest, "username field not found")
			return
			// msg := data{
			// 	Message: "username field not found",
			// }
			// if err := h.Tmpl.ExecuteTemplate(w, "singUp.html", msg); err != nil {
			// 	h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			// 	return
			// }
			// return
		}
		password, ok := r.Form["password"]
		passwordCheck := strings.Trim(password[0], " ")
		if !ok || password[0] == "" || passwordCheck == "" {
			h.errorPage(w, r, http.StatusBadRequest, "password field not found")
			return
			// msg := data{
			// 	Message: "password field not found",
			// }
			// if err := h.Tmpl.ExecuteTemplate(w, "singUp.html", msg); err != nil {
			// 	h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			// 	return
			// }
			// return
		}
		verifyPassword, ok := r.Form["verifyPassword"]
		verifyPasswordCheck := strings.Trim(verifyPassword[0], " ")
		if !ok || verifyPassword[0] == "" || verifyPasswordCheck == " " {
			h.errorPage(w, r, http.StatusBadRequest, "verifyPassword field not found")
			return
			// msg := data{
			// 	Message: "verifyPassword field not found",
			// }
			// if err := h.Tmpl.ExecuteTemplate(w, "singUp.html", msg); err != nil {
			// 	h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			// 	return
			// }
			// return
		}

		user := models.User{
			Email:          email[0],
			Username:       username[0],
			Password:       password[0],
			VerifyPassword: verifyPassword[0],
		}
		if err := h.Services.Auth.CreateUser(user); err != nil {
			if errors.Is(err, service.ErrInvalidUserName) ||
				errors.Is(err, service.ErrPasswordDontMatch) ||
				errors.Is(err, service.ErrInvalidEmail) ||
				errors.Is(err, service.ErrInvalidPassword) ||
				errors.Is(err, service.ErrUserExist) {
				h.errorPage(w, r, http.StatusBadRequest, err.Error())
				return
			}
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}
