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
	storage *storage.Storage
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
		id := path.Base(r.URL.String())

		if id == "" {
			http.Error(w, "The parameter is missing", http.StatusBadRequest)

			w.WriteHeader(http.StatusNotFound)
			return
		}

		url, err := h.storage.LinkBy(id)

		if err == nil {
			w.Header().Set("Location", url)

			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		return
	case http.MethodPost:
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
		return
	default:
		setBadResponse(w)
		return
	}
}

//func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
//	if r.Method == http.MethodGet {
//		param := path.Base(r.URL.String())
//
//		if param == "" {
//			http.Error(w, "The parameter is missing", http.StatusBadRequest)
//
//			w.WriteHeader(http.StatusNotFound)
//			return
//		}
//
//		url, err := h.storage.LinkBy(param)
//		if err != nil {
//			w.WriteHeader(http.StatusNotFound)
//			return
//		}
//
//		w.Header().Set("Location", url)
//
//		w.WriteHeader(http.StatusTemporaryRedirect)
//	}
//
//	setBadResponse(w)
//}
//
//func (h *Handler) Save(w http.ResponseWriter, r *http.Request) {
//	if r.Method == http.MethodPost {
//		if r.Body != http.NoBody {
//			body, err := io.ReadAll(r.Body)
//
//			if err != nil {
//				http.Error(w, err.Error(), 500)
//				return
//			}
//
//			sl := h.storage.Save(string(body))
//
//			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
//			w.WriteHeader(http.StatusCreated)
//
//			slURL := fmt.Sprintf("%s/%s", Host, string(sl))
//
//			_, err = w.Write([]byte(slURL))
//			if err == nil {
//				return
//			}
//		}
//	}
//
//	setBadResponse(w)
//}

func StatusHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	// намеренно сделана ошибка в JSON
	rw.Write([]byte(`{"status":"ok"}`))
}

func setBadResponse(w http.ResponseWriter) {
	http.Error(w, errors.New("bad request").Error(), http.StatusBadRequest)
}
