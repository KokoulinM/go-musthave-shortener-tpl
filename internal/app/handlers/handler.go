package handlers

import (
	"errors"
	"fmt"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/storage"
	"io"
	"net/http"
	"path"
)

type Handler struct {
	storage storage.Repository
}

const Host = "http://localhost:8080"

func New() *Handler {
	return &Handler{
		storage: storage.New(),
	}
}

func (h *Handler) CommonHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
		return
	case http.MethodPost:
		h.Save(w, r)
		return
	default:
		setBadResponse(w)
		return
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id := path.Base(r.URL.String())

		if id == "" {
			http.Error(w, "The parameter is missing", http.StatusBadRequest)

			w.WriteHeader(http.StatusNotFound)
			return
		}

		url, err := h.storage.LinkBy(id)

		if err == nil {
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		return
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

			sl := h.storage.Save(string(body))

			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusCreated)

			slURL := fmt.Sprintf("%s/%s", Host, string(sl))

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
