package handler

import (
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/storage"
	"net/http"
	"path"
)

type Handler struct {
	storage storage.Repository
}

func New() *Handler {
	return &Handler{
		storage: storage.New(),
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	param := path.Base(r.URL.String())

	if param == "" {
		http.Error(w, "The parameter is missing", http.StatusBadRequest)

		w.WriteHeader(http.StatusNotFound)
		return
	}

	url, err := h.storage.LinkBy(param)
	if err == nil {
		w.Header().Set("Location", url)

		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

//func Handler(w http.ResponseWriter, r *http.Request) {
//	switch {
//	case r.Method == http.MethodGet:
//		GetHandler(w, r)
//		return
//	case r.Method == http.MethodPost:
//		w.Write([]byte("Hola, Mundo"))
//		return
//	default:
//		http.NotFound(w, r)
//		return
//	}
//}
