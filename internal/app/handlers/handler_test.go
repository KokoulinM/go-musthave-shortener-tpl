package handlers

import (
	"fmt"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_CommonHandler(t *testing.T) {
	h := New()

	type fields struct {
		storage *storage.Storage
	}
	type want struct {
		code        int
		response    string
		contentType string
	}
	type request struct {
		method string
		target string
		body   string
		path   string
	}
	tests := []struct {
		name    string
		want    want
		request request
	}{
		{
			name: "Save handler #1",
			want: want{
				code:        http.StatusCreated,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
			request: request{
				method: http.MethodPost,
				target: "/",
				body:   "http://site.ru/123",
				path:   "/",
			},
		},
		//{
		//	name: "Get handler #1",
		//	want: want{
		//		code:        http.StatusTemporaryRedirect,
		//		response:    "",
		//		contentType: "text/plain; charset=utf-8",
		//	},
		//	request: request{
		//		method: http.MethodGet,
		//		target: "/",
		//		body:   "http://site.ru",
		//		path:   "/",
		//	},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body io.Reader

			if len(tt.request.body) > 0 {
				body = strings.NewReader(tt.request.body)
			} else {
				body = nil
			}

			request := httptest.NewRequest(tt.request.method, tt.request.path, body)

			w := httptest.NewRecorder()

			h := http.HandlerFunc(h.CommonHandler)

			h.ServeHTTP(w, request)

			response := w.Result()

			assert.Equal(t, tt.want.code, response.StatusCode, "invalid response code")

			defer response.Body.Close()

			resultBody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Println(string(resultBody))

			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
		})
	}
}
