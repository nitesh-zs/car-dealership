package middleware

import (
	"fmt"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check x-api-key in request header
		if r.Header.Get("x-api-key") != "nitesh-zs" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"error":{"code":"Authorization error","message":"A valid 'x-api-key' must be set in request headers"}}`)
			return
		}
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func RespHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// set Content-Type in response header
		w.Header().Set("Content-Type", "application/json")
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
