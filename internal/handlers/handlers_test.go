package handlers

import "testing"

func TestCreateShortURL(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name   string
		body   string
		result string
		want   want
	}{
		{
			name:   "positive test",
			body:   "https://go.dev",
			result: "suNu2gXD_LK3aL0z5AsfO6ywMEM=",
			want: want{
				code:        201,
				response:    "http://localhost:8080/suNu2gXD_LK3aL0z5AsfO6ywMEM=",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

//func TestGetHandler(t *testing.T) {
//	c := configs.New()
//	database := database.New(new(sql.DB))
//	h := New(database, c)
//	s := database.MockStorage{
//		Data: make(map[database.UserID]database.ShortLinks),
//	}
//
//	s.GenerateMockData()
//
//	h.database = &s
//
//	type want struct {
//		code        int
//		response    string
//		contentType string
//	}
//	type request struct {
//		method string
//		target string
//		path   string
//	}
//	tests := []struct {
//		name    string
//		want    want
//		request request
//	}{
//		{
//			name: "simple test Get handler #1",
//			want: want{
//				code:        http.StatusTemporaryRedirect,
//				response:    "<a href=\"https://go.dev\">Temporary Redirect</a>.\n\n",
//				contentType: "text/plain; charset=utf-8",
//			},
//			request: request{
//				method: http.MethodGet,
//				target: "https://go.dev/GMWJGSAPGA_test_1",
//				path:   "/{id}",
//			},
//		},
//		{
//			name: "negative test Get handler without param",
//			want: want{
//				code:        http.StatusBadRequest,
//				response:    "<a href=\"https://go.dev\">Temporary Redirect</a>.\n\n",
//				contentType: "text/plain; charset=utf-8",
//			},
//			request: request{
//				method: http.MethodGet,
//				target: "https://go.dev",
//				path:   "/",
//			},
//		},
//		{
//			name: "negative test Get handler empty row in the database",
//			want: want{
//				code:        http.StatusNotFound,
//				response:    "<a href=\"https://go.dev\">Temporary Redirect</a>.\n\n",
//				contentType: "text/plain; charset=utf-8",
//			},
//			request: request{
//				method: http.MethodGet,
//				target: "https://vk.com/123",
//				path:   "/{id}",
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			request := httptest.NewRequest(tt.request.method, tt.request.target, nil)
//
//			w := httptest.NewRecorder()
//
//			router := chi.NewRouter()
//
//			router.Get(tt.request.path, h.Get)
//
//			router.ServeHTTP(w, request.WithContext(context.WithValue(request.Context(), middlewares.UserIDCtxName, "userID")))
//
//			response := w.Result()
//
//			defer response.Body.Close()
//
//			assert.Equal(t, tt.want.code, response.StatusCode, "invalid response code")
//
//			_, err := ioutil.ReadAll(response.Body)
//
//			if err != nil {
//				t.Fatal(err)
//			}
//		})
//	}
//}
//
//func TestSaveHandler(t *testing.T) {
//	c := configs.New()
//	database := database.New(new(sql.DB))
//	h := New(database, c)
//	s := database.MockStorage{
//		Data: make(map[database.UserID]database.ShortLinks),
//	}
//
//	s.GenerateMockData()
//
//	h.database = &s
//
//	type want struct {
//		code        int
//		response    string
//		contentType string
//	}
//	type request struct {
//		method string
//		target string
//		path   string
//	}
//	tests := []struct {
//		name    string
//		want    want
//		request request
//	}{
//		{
//			name: "simple test Post handler #1",
//			want: want{
//				code:        http.StatusCreated,
//				response:    "https://go.dev/GMWJGSAPGA",
//				contentType: "text/plain; charset=utf-8",
//			},
//			request: request{
//				method: http.MethodPost,
//				target: "https://go.dev",
//				path:   "/",
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			request := httptest.NewRequest(tt.request.method, tt.request.target, strings.NewReader(tt.request.target))
//
//			w := httptest.NewRecorder()
//
//			router := chi.NewRouter()
//
//			router.Post(tt.request.path, h.Save)
//
//			router.ServeHTTP(w, request)
//
//			response := w.Result()
//
//			defer response.Body.Close()
//
//			assert.Equal(t, tt.want.code, response.StatusCode, "invalid response code")
//
//			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"), "invalid response Content-Type")
//
//			_, err := ioutil.ReadAll(response.Body)
//
//			if err != nil {
//				t.Fatal(err)
//			}
//		})
//	}
//}
//
//func TestHandlerSaveJSON(t *testing.T) {
//	c := configs.New()
//	database := database.New(new(sql.DB))
//	h := New(database, c)
//	s := database.MockStorage{
//		Data: make(map[database.UserID]database.ShortLinks),
//	}
//
//	s.GenerateMockData()
//
//	h.database = &s
//
//	type want struct {
//		code        int
//		response    string
//		contentType string
//		body        string
//	}
//
//	type request struct {
//		method string
//		target string
//		path   string
//		body   string
//	}
//
//	tests := []struct {
//		name    string
//		want    want
//		request request
//	}{
//		{
//			name: "simple test Post handler #1",
//			want: want{
//				code:        http.StatusCreated,
//				response:    "{\"result\":\"http://localhost:8080/CAKKMYDSJD_test_14\"}",
//				contentType: "application/json; charset=utf-8",
//			},
//			request: request{
//				method: http.MethodPost,
//				target: "/",
//				path:   "/",
//				body:   "{\"url\":\"https://go.dev/123\"}",
//			},
//		},
//		{
//			name: "negative test Post handler",
//			want: want{
//				code:        http.StatusBadRequest,
//				response:    "the URL property is missing\n",
//				contentType: "text/plain; charset=utf-8",
//			},
//			request: request{
//				method: http.MethodPost,
//				target: "/",
//				path:   "/",
//				body:   "{\"url2\":\"https://go.dev/123\"}",
//			},
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			reader := strings.NewReader(tt.request.body)
//
//			request := httptest.NewRequest(tt.request.method, tt.request.target, reader)
//
//			w := httptest.NewRecorder()
//
//			router := chi.NewRouter()
//
//			router.Post(tt.request.path, h.SaveJSON)
//
//			router.ServeHTTP(w, request)
//
//			response := w.Result()
//
//			defer response.Body.Close()
//
//			body, _ := ioutil.ReadAll(response.Body)
//
//			assert.Equal(t, tt.want.code, response.StatusCode, "invalid response code")
//
//			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"), "invalid response Content-Type")
//
//			assert.Equal(t, tt.want.response, string(body), "invalid response body")
//		})
//	}
//}
