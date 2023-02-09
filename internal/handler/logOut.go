package handler

import (
	"net/http"
	"time"
)

func (h *Handler) logOut(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/logout" {
		h.errorPage(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			h.errorPage(w, r, http.StatusUnauthorized, err.Error())
			return
		}
		h.errorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.Services.DeleteSessionToken(c.Value); err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
