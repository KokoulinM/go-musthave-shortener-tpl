package handlers

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/handlers/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/configs"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/workers"
)

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
			mockURL:   "Vq7zU8E5b7sLZo3qY82UKYRvQ-A=",
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
			mockURL:   "Vq7zU8E5b7sLZo3qY82UKYRvQ-A=",
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
			mockURL:   "Vq7zU8E5b7sLZo3qY82UKYRvQ-A=",
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
			mockURL:   "Vq7zU8E5b7sLZo3qY82UKYRvQ-A=",
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

			r := chi.NewRouter()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cfg := configs.New()

			wp := workers.New(context.Background(), cfg.Workers, cfg.WorkersBuffer)

			repoMock := NewMockRepository(ctrl)

			h := New(repoMock, cfg.BaseURL, wp)

			r.Post(tt.query, h.CreateShortURL)

			repoMock.EXPECT().AddURL(gomock.Any(), tt.body, tt.mockURL, "userID").Return(tt.mockError).AnyTimes()

			r.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), middlewares.UserIDCtxName, "userID")))

			response := w.Result()

			defer response.Body.Close()

			body, _ := ioutil.ReadAll(response.Body)

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.response, string(body))
		})
	}
}
