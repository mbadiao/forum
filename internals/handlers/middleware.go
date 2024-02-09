package handlers

import (
	"forum/internals/utils"
	"net/http"
)

func ErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		found := false
		for _, route := range Routes {
			if route.Path == r.URL.Path {
				found = true
			}
		}
		if !found {
			w.WriteHeader(404)
			utils.Render(w, "error", Err[404])
			return
		}
		next.ServeHTTP(w, r)
	})
}
