package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang/mock/gomock"
	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/handlers/middlewares"
	"github.com/stretchr/testify/assert"

	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/configs"
	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/workers"
)

func router(h *Handlers) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/", func(r chi.Router) {
		router.Post("/", h.CreateShortURL)
		router.Get("/{id}", h.RetrieveShortURL)
		router.Get("/ping", h.PingDB)
		router.Post("/api/shorten", h.ShortenURL)
		router.Get("/api/user/urls", h.GetUserURLs)
		router.Delete("/api/user/urls", h.DeleteBatch)
		router.Post("/api/shorten/batch", h.CreateBatch)
		router.Get("/api/internal/states", h.GetStates)
	})

	return router
}

func TestCreateShortURL(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name      string
		query     string
		body      string
		mockError error
		mockURL   string
		want      want
	}{
		{
			name:      "positive test",
			query:     "/",
			body:      "https://go.dev",
			mockError: nil,
			mockURL:   "http://localhost:8080/Vq7zU8E5b7sLZo3qY82UKYRvQ-A=",
			want: want{
				code:        http.StatusCreated,
				response:    "http://localhost:8080/Vq7zU8E5b7sLZo3qY82UKYRvQ-A=",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:      "empty body",
			query:     "/",
			body:      "",
			mockError: nil,
			mockURL:   "http://localhost:8080/Vq7zU8E5b7sLZo3qY82UKYRvQ-A=",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				response:    "the body cannot be an empty\n",
			},
		},
		{
			name:      "unexpected error when adding to the database",
			query:     "/",
			body:      "https://go.dev",
			mockError: errors.New("error"),
			mockURL:   "http://localhost:8080/Vq7zU8E5b7sLZo3qY82UKYRvQ-A=",
			want: want{
				code:        http.StatusInternalServerError,
				contentType: "text/plain; charset=utf-8",
				response:    "",
			},
		},
		{
			name:      "the url already exists in the database",
			query:     "/",
			body:      "https://go.dev",
			mockError: NewErrorWithDB(errors.New("UniqConstraint"), "UniqConstraint"),
			mockURL:   "http://localhost:8080/Vq7zU8E5b7sLZo3qY82UKYRvQ-A=",
			want: want{
				code:        http.StatusConflict,
				contentType: "text/plain; charset=utf-8",
				response:    "http://localhost:8080/Vq7zU8E5b7sLZo3qY82UKYRvQ-A=",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, tt.query, strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cfg := configs.New()

			wp := workers.New(context.Background(), cfg.Workers, cfg.WorkersBuffer)

			repoMock := NewMockURLServiceInterface(ctrl)

			h := New(repoMock, cfg.BaseURL, wp)

			r := router(h)

			repoMock.EXPECT().CreateURL(gomock.Any(), tt.body, "userID").Return(tt.mockURL, tt.mockError).AnyTimes()

			r.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), middlewares.UserIDCtxName, "userID")))

			response := w.Result()

			defer response.Body.Close()

			body, _ := ioutil.ReadAll(response.Body)

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.response, string(body))
		})
	}
}

func TestRetrieveShortURL(t *testing.T) {
	type want struct {
		code int
	}

	tests := []struct {
		name      string
		query     string
		mockError error
		mockID    string
		mockURL   string
		want      want
	}{
		{
			name:      "positive test",
			query:     "/Vq7zU8E5b7sLZo3qY82UKYRvQ-A=",
			mockError: nil,
			mockID:    "Vq7zU8E5b7sLZo3qY82UKYRvQ-A=",
			mockURL:   "https://go.dev",
			want: want{
				code: http.StatusTemporaryRedirect,
			},
		},
		{
			name:      "deleted",
			query:     "/Vq7zU8E5b7sLZo3qY82UKYRvQ-A=",
			mockError: NewErrorWithDB(errors.New("deleted"), "deleted"),
			mockID:    "Vq7zU8E5b7sLZo3qY82UKYRvQ-A=",
			mockURL:   "https://go.dev",
			want: want{
				code: http.StatusGone,
			},
		},
		{
			name:      "deleted",
			query:     "/Vq7zU8E5b7sLZo3qY82UKYRvQ-A=",
			mockError: errors.New(""),
			mockID:    "Vq7zU8E5b7sLZo3qY82UKYRvQ-A=",
			mockURL:   "https://go.dev",
			want: want{
				code: http.StatusNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, tt.query, nil)

			w := httptest.NewRecorder()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cfg := configs.New()

			wp := workers.New(context.Background(), cfg.Workers, cfg.WorkersBuffer)

			repoMock := NewMockURLServiceInterface(ctrl)

			h := New(repoMock, cfg.BaseURL, wp)

			r := router(h)

			repoMock.EXPECT().GetURL(gomock.Any(), tt.mockID).Return(tt.mockURL, tt.mockError).AnyTimes()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.want.code, w.Code)
		})
	}
}

func TestShortenURL(t *testing.T) {
	type want struct {
		code     int
		response string
	}

	tests := []struct {
		name      string
		query     string
		body      string
		mockError error
		mockURL   string
		want      want
	}{
		{
			name:      "positive test",
			query:     "/api/shorten",
			body:      `{"url":"https://go.dev"}`,
			mockError: nil,
			mockURL:   "http://localhost:8080/Vq7zU8E5b7sLZo3qY82UKYRvQ-A=",
			want: want{
				code:     http.StatusCreated,
				response: `{"result":"http://localhost:8080/Vq7zU8E5b7sLZo3qY82UKYRvQ-A="}`,
			},
		},
		{
			name:      "when unmarshaling JSON",
			query:     "/api/shorten",
			body:      `{"url":"}`,
			mockError: nil,
			mockURL:   "http://localhost:8080/Vq7zU8E5b7sLZo3qY82UKYRvQ-A=",
			want: want{
				code:     http.StatusInternalServerError,
				response: "an unexpected error when unmarshaling JSON\n",
			},
		},
		{
			name:      "the URL property is missing",
			query:     "/api/shorten",
			body:      `{"url":""}`,
			mockError: nil,
			mockURL:   "http://localhost:8080/Vq7zU8E5b7sLZo3qY82UKYRvQ-A=",
			want: want{
				code:     http.StatusBadRequest,
				response: "the URL property is missing\n",
			},
		},
		{
			name:      "UniqConstraint",
			query:     "/api/shorten",
			body:      `{"url":"https://go.dev"}`,
			mockError: NewErrorWithDB(errors.New("UniqConstraint"), "UniqConstraint"),
			mockURL:   "http://localhost:8080/Vq7zU8E5b7sLZo3qY82UKYRvQ-A=",
			want: want{
				code:     http.StatusConflict,
				response: "{\"result\":\"http://localhost:8080/Vq7zU8E5b7sLZo3qY82UKYRvQ-A=\"}",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, tt.query, strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cfg := configs.New()

			wp := workers.New(context.Background(), cfg.Workers, cfg.WorkersBuffer)

			repoMock := NewMockURLServiceInterface(ctrl)

			h := New(repoMock, cfg.BaseURL, wp)

			r := router(h)

			url := URL{}

			_ = json.Unmarshal(bytes.NewBufferString(tt.body).Bytes(), &url)

			repoMock.EXPECT().CreateURL(gomock.Any(), url.URL, "userID").Return(tt.mockURL, tt.mockError).AnyTimes()

			r.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), middlewares.UserIDCtxName, "userID")))

			response := w.Result()

			defer response.Body.Close()

			body, _ := ioutil.ReadAll(response.Body)

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.response, string(body))
		})
	}
}

func TestGetUserURLs(t *testing.T) {
	type want struct {
		code     int
		response string
	}

	tests := []struct {
		name      string
		query     string
		mockError error
		mockURLs  []ResponseGetURL
		want      want
	}{
		{
			name:      "positive test",
			query:     "/api/user/urls",
			mockError: nil,
			mockURLs:  []ResponseGetURL{ResponseGetURL{ShortURL: "http://localhost:8080/Vq7zU8E5b7sLZo3qY82UKYRvQ-A=", OriginalURL: "https://go.dev"}},
			want: want{
				code:     http.StatusOK,
				response: `[{"short_url":"http://localhost:8080/Vq7zU8E5b7sLZo3qY82UKYRvQ-A=","original_url":"https://go.dev"}]`,
			},
		},
		{
			name:      "no content",
			query:     "/api/user/urls",
			mockError: nil,
			mockURLs:  []ResponseGetURL{},
			want: want{
				code:     http.StatusNoContent,
				response: "no content\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, tt.query, nil)
			w := httptest.NewRecorder()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cfg := configs.New()

			wp := workers.New(context.Background(), cfg.Workers, cfg.WorkersBuffer)

			repoMock := NewMockURLServiceInterface(ctrl)

			h := New(repoMock, cfg.BaseURL, wp)

			r := router(h)

			repoMock.EXPECT().GetUserURLs(gomock.Any(), "userID").Return(tt.mockURLs, tt.mockError).AnyTimes()

			r.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), middlewares.UserIDCtxName, "userID")))

			response := w.Result()

			defer response.Body.Close()

			body, _ := ioutil.ReadAll(response.Body)

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.response, string(body))
		})
	}
}

func TestDeleteBatch(t *testing.T) {
	type want struct {
		code     int
		response string
	}

	tests := []struct {
		name      string
		query     string
		body      string
		mockError error
		mockURLs  []string
		want      want
	}{
		{
			name:      "positive test",
			query:     "/api/user/urls",
			body:      `["", ""]`,
			mockError: nil,
			mockURLs:  []string{"", ""},
			want: want{
				code: http.StatusAccepted,
			},
		},
		{
			name:      "unexpected end of JSON input",
			query:     "/api/user/urls",
			body:      `["]`,
			mockError: nil,
			mockURLs:  []string{"", ""},
			want: want{
				code:     http.StatusInternalServerError,
				response: "unexpected end of JSON input\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodDelete, tt.query, strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cfg := configs.New()

			wp := workers.New(context.Background(), cfg.Workers, cfg.WorkersBuffer)

			repoMock := NewMockURLServiceInterface(ctrl)

			h := New(repoMock, cfg.BaseURL, wp)

			r := router(h)

			repoMock.EXPECT().DeleteBatch(tt.mockURLs, "userID").AnyTimes()

			r.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), middlewares.UserIDCtxName, "userID")))

			response := w.Result()

			defer response.Body.Close()

			body, _ := ioutil.ReadAll(response.Body)

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.response, string(body))
		})
	}
}

func TestCreateBatch(t *testing.T) {
	type want struct {
		code     int
		response string
	}

	tests := []struct {
		name        string
		query       string
		body        string
		mockError   error
		mockReqURLs []RequestGetURLs
		mockResURLs []ResponseGetURLs
		want        want
	}{
		{
			name:        "positive test",
			query:       "/api/shorten/batch",
			body:        `[{"correlation_id":"","original_url":""}]`,
			mockError:   nil,
			mockReqURLs: []RequestGetURLs{RequestGetURLs{}},
			mockResURLs: []ResponseGetURLs{ResponseGetURLs{}},
			want: want{
				code:     http.StatusCreated,
				response: "[{\"correlation_id\":\"\",\"short_url\":\"\"}]",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, tt.query, strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cfg := configs.New()

			wp := workers.New(context.Background(), cfg.Workers, cfg.WorkersBuffer)

			repoMock := NewMockURLServiceInterface(ctrl)

			h := New(repoMock, cfg.BaseURL, wp)

			r := router(h)

			repoMock.EXPECT().CreateBatch(gomock.Any(), tt.mockReqURLs, "userID").Return(tt.mockResURLs, tt.mockError).AnyTimes()

			r.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), middlewares.UserIDCtxName, "userID")))

			response := w.Result()

			defer response.Body.Close()

			body, _ := ioutil.ReadAll(response.Body)

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.response, string(body))
		})
	}
}

func TestHandlers_GetStates(t *testing.T) {
	type want struct {
		code     int
		response string
	}
	tests := []struct {
		name           string
		query          string
		mockAllowedURL net.IP
		isAllowed      bool
		mockState      ResponseStates
		want           want
	}{
		{
			name:           "Test #1",
			query:          "/api/internal/states",
			mockAllowedURL: net.ParseIP("http://localhost:3333"),
			isAllowed:      true,
			mockState: ResponseStates{
				Urls:  7,
				Users: 3,
			},
			want: want{
				code:     http.StatusOK,
				response: `{"urls":7,"users":3}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, tt.query, nil)
			w := httptest.NewRecorder()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cfg := configs.New()

			wp := workers.New(context.Background(), cfg.Workers, cfg.WorkersBuffer)

			repoMock := NewMockURLServiceInterface(ctrl)

			h := New(repoMock, cfg.BaseURL, wp)

			r := router(h)

			repoMock.EXPECT().GetStates(gomock.Any(), tt.mockAllowedURL).Return(tt.isAllowed, tt.mockState, nil).AnyTimes()

			r.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), middlewares.UserIDCtxName, "userID")))

			response := w.Result()

			defer response.Body.Close()

			body, _ := ioutil.ReadAll(response.Body)

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.response, string(body))
		})
	}
}
