package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) errorPage(w http.ResponseWriter, r *http.Request, status int, msg string) {
	w.WriteHeader(status)
	data := struct {
		Status  int
		Message string
	}{
		Status:  status,
		Message: msg,
	}

	if err := h.Tmpl.ExecuteTemplate(w, "error.html", data); err != nil {
		fmt.Fprintf(w, "%d - %s\n", data.Status, data.Message)
	}
}
