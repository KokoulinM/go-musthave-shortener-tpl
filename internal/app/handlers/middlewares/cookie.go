package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/helpers"
)

const CookieUserIDName = "user_id"
const UserIDCtxName = "ctxUserId"

func CookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID string

		if cookieUserID, err := r.Cookie(CookieUserIDName); err == nil {
			err := helpers.Decode(cookieUserID.Value, &userID)

			if err != nil {
				fmt.Printf("error: %v\n", err)
			}
		}

		encoded, err := helpers.Encode(userID)

		if err == nil {
			cookie := helpers.CreateCookie(CookieUserIDName, encoded)
			http.SetCookie(w, cookie)
		} else {
			fmt.Printf("error: %v\n", err)
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserIDCtxName, userID)))
	})
}
