package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/helpers"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/helpers/encryptor"
	"github.com/google/uuid"
)

const CookieUserIDName = "user_id"

type ContextType string

const UserIDCtxName ContextType = "ctxUserId"

func CookieMiddleware(key []byte) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookieUserID, _ := r.Cookie(CookieUserIDName)
			encryptor, err := encryptor.New(key)

			log.Println("cookieUserID.Value: ", cookieUserID.Value)

			if err != nil {
				return
			}

			if cookieUserID != nil {
				userID, err := encryptor.Decode(cookieUserID.Value)

				if err == nil {
					next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserIDCtxName, userID)))
					return
				}
			}
			userID := uuid.New().String()
			encoded := encryptor.Encode([]byte(userID))
			cookie := helpers.CreateCookie(CookieUserIDName, encoded)

			http.SetCookie(w, cookie)
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserIDCtxName, userID)))
		})
	}
}
