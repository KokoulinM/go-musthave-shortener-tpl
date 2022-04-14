package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/configs"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/storage"
)

type Handler struct {
	storage storage.Repository
	config  configs.Config
}

type URL struct {
	URL string `json:"url"`
}

func New(c configs.Config) *Handler {
	h := &Handler{
		storage: storage.New(),
		config:  c,
	}

	if err := h.storage.Load(h.config); err != nil {
		panic(err)
	}

	return h
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "the parameter is missing", http.StatusBadRequest)
		return
	}

	url, err := h.storage.LinkBy(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) Save(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		http.Error(w, "the body cannot be an empty", http.StatusBadRequest)
		return
	}

	origin := string(body)

	short := string(h.storage.Save(origin))

	defer h.storage.Flush(h.config)

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	slURL := fmt.Sprintf("%s/%s", h.config.BaseURL, short)

	w.Write([]byte(slURL))
}

func (h *Handler) SaveJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	body, errReadAll := io.ReadAll(r.Body)

	if errReadAll != nil {
		http.Error(w, errReadAll.Error(), http.StatusInternalServerError)
		return
	}

	url := URL{}

	errUnmarshal := json.Unmarshal(body, &url)

	if errUnmarshal != nil {
		http.Error(w, "an unexpected error when unmarshaling JSON", http.StatusInternalServerError)
		return
	}

	if url.URL == "" {
		http.Error(w, "the URL property is missing", http.StatusBadRequest)
		return
	}

	sl := h.storage.Save(url.URL)

	defer h.storage.Flush(h.config)

	slURL := fmt.Sprintf("%s/%s", h.config.BaseURL, string(sl))

	result := struct {
		Result string `json:"result"`
	}{
		Result: slURL,
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	body, errMarshal := json.Marshal(result)

	if errMarshal != nil {
		http.Error(w, "an unexpected error when marshaling JSON", http.StatusInternalServerError)
		return
	}

	w.Write(body)
}
