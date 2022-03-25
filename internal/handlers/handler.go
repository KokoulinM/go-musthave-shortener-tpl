package handler

import (
	"errors"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/storage"
	"io"
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
	if r.Method == http.MethodGet {
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
			return
		}
	}

	setBadResponse(w)
}

func (h *Handler) Save(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if r.Body != http.NoBody {
			body, err := io.ReadAll(r.Body)

			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusCreated)

			slURL := string(body)

			_, err = w.Write([]byte(slURL))
			if err == nil {
				return
			}
		}
	}

	setBadResponse(w)
}

func setBadResponse(w http.ResponseWriter) {
	http.Error(w, errors.New("bad request").Error(), http.StatusBadRequest)
}
