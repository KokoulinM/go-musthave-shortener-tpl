package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/handlers/middlewares"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/models"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/shortener"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/workers"
)

type Repository interface {
	AddURL(ctx context.Context, longURL models.LongURL, shortURL models.ShortURL, user models.UserID) error
	GetURL(ctx context.Context, shortURL models.ShortURL) (models.ShortURL, error)
	GetUserURLs(ctx context.Context, user models.UserID) ([]ResponseGetURL, error)
	DeleteMultipleURLs(ctx context.Context, user models.UserID, urls ...string) error
	Ping(ctx context.Context) error
	AddMultipleURLs(ctx context.Context, user models.UserID, urls ...RequestGetURLs) ([]ResponseGetURLs, error)
}

type Handlers struct {
	repo    Repository
	baseURL string
	wp      *workers.WorkerPool
}

type URL struct {
	URL string `json:"url"`
}

type ResponseGetURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type RequestGetURLs struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ResponseGetURLs struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type ErrorWithDB struct {
	Err   error
	Title string
}

func (err *ErrorWithDB) Error() string {
	return fmt.Sprintf("%v", err.Err)
}

func (err *ErrorWithDB) Unwrap() error {
	return err.Err
}

func NewErrorWithDB(err error, title string) error {
	return &ErrorWithDB{
		Err:   err,
		Title: title,
	}
}

func New(repo Repository, baseURL string, wp *workers.WorkerPool) *Handlers {
	return &Handlers{
		repo:    repo,
		baseURL: baseURL,
		wp:      wp,
	}
}

func (h *Handlers) RetrieveShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "the parameter is missing", http.StatusBadRequest)
		return
	}

	url, err := h.repo.GetURL(r.Context(), id)
	if err != nil {
		var dbErr *ErrorWithDB

		if errors.As(err, &dbErr) && dbErr.Title == "deleted" {
			w.WriteHeader(http.StatusGone)
			return
		}

		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Location", url)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handlers) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		http.Error(w, "the body cannot be an empty", http.StatusBadRequest)
		return
	}

	userIDCtx := r.Context().Value(middlewares.UserIDCtxName)

	userID := "default"

	if userIDCtx != nil {
		userID = userIDCtx.(string)
	}

	longURL := models.LongURL(body)
	shortURL := shortener.ShorterURL(longURL)

	err = h.repo.AddURL(r.Context(), longURL, shortURL, userID)
	if err != nil {
		var dbErr *ErrorWithDB

		if errors.As(err, &dbErr) && dbErr.Title == "UniqConstraint" {
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")

			w.WriteHeader(http.StatusConflict)

			slURL := fmt.Sprintf("%s/%s", h.baseURL, shortURL)

			_, err = w.Write([]byte(slURL))
			if err != nil {
				http.Error(w, "unexpected error when writing the response body", http.StatusInternalServerError)
			}

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	slURL := fmt.Sprintf("%s/%s", h.baseURL, shortURL)

	_, err = w.Write([]byte(slURL))
	if err != nil {
		http.Error(w, "unexpected error when writing the response body", http.StatusInternalServerError)
	}
}

func (h *Handlers) ShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	result := map[string]string{}

	body, errReadAll := io.ReadAll(r.Body)
	if errReadAll != nil {
		http.Error(w, errReadAll.Error(), http.StatusInternalServerError)
		return
	}

	url := URL{}

	err := json.Unmarshal(body, &url)
	if err != nil {
		http.Error(w, "an unexpected error when unmarshaling JSON", http.StatusInternalServerError)
		return
	}

	if url.URL == "" {
		http.Error(w, "the URL property is missing", http.StatusBadRequest)
		return
	}

	userIDCtx := r.Context().Value(middlewares.UserIDCtxName)

	userID := "default"

	if userIDCtx != nil {
		userID = userIDCtx.(string)
	}

	shortURL := shortener.ShorterURL(url.URL)

	slURL := fmt.Sprintf("%s/%s", h.baseURL, shortURL)

	err = h.repo.AddURL(r.Context(), url.URL, shortURL, userID)
	if err != nil {
		var dbErr *ErrorWithDB
		if errors.As(err, &dbErr) && dbErr.Title == "UniqConstraint" {
			result["result"] = slURL

			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusConflict)

			body, err = json.Marshal(result)
			if err != nil {
				http.Error(w, "an unexpected error when marshaling JSON", http.StatusInternalServerError)
				return
			}

			_, err = w.Write(body)
			if err != nil {
				http.Error(w, "unexpected error when writing the response body", http.StatusInternalServerError)
				return
			}

			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result["result"] = slURL

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	body, err = json.Marshal(result)
	if err != nil {
		http.Error(w, "an unexpected error when marshaling JSON", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		http.Error(w, "unexpected error when writing the response body", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDCtx := r.Context().Value(middlewares.UserIDCtxName)

	userID := "default"

	if userIDCtx != nil {
		userID = userIDCtx.(string)
	}

	urls, err := h.repo.GetUserURLs(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(urls) == 0 {
		http.Error(w, errors.New("no content").Error(), http.StatusNoContent)
		return
	}

	body, err := json.Marshal(urls)

	if err == nil {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		w.WriteHeader(http.StatusOK)

		_, err = w.Write(body)
		if err == nil {
			return
		}
	}
}

func (h *Handlers) DeleteBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "only DELETE requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDCtx := r.Context().Value(middlewares.UserIDCtxName)

	userID := "default"

	if userIDCtx != nil {
		userID = userIDCtx.(string)
	}

	defer r.Body.Close()

	var data []string

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var sliceData [][]string

	for i := 10; i <= len(data); i += 10 {
		sliceData = append(sliceData, data[i-10:i])
	}

	rem := len(data) % 10
	if rem > 0 {
		sliceData = append(sliceData, data[len(data)-rem:])
	}

	for _, item := range sliceData {
		func(taskData []string) {
			h.wp.Push(func(ctx context.Context) error {
				err := h.repo.DeleteMultipleURLs(ctx, userID, taskData...)
				return err
			})
		}(item)
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handlers) CreateBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var data []RequestGetURLs

	userIDCtx := r.Context().Value(middlewares.UserIDCtxName)

	userID := "default"

	if userIDCtx != nil {
		userID = userIDCtx.(string)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	urls, err := h.repo.AddMultipleURLs(r.Context(), userID, data...)
	if err != nil {
		log.Println("err.Error(): ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err = json.Marshal(urls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(body)
	if err != nil {
		http.Error(w, "unexpected error when writing the response body", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) PingDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	err := h.repo.Ping(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
