package middlewares

import "net/http"

func InternalMiddleware(trustedSubnet string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ts := r.Header.Get("X-Real-IP")

			if ts == "" || ts != trustedSubnet {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
