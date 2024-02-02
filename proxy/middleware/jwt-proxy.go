package middleware

import (
	"github.com/go-chi/jwtauth/v5"
	"log"
	"net/http"
)

func TokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}
		token, ok := claims["id"]
		log.Println(token)
		if !ok {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
