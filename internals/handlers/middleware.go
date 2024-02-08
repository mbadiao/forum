package handlers

import (
	"forum/internals/utils"
	"net/http"
)

func ErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Path
		found := false
		for _, route := range Routes {
			if route.Path == page {
				found = true
			}
		}
		if !found {
			w.WriteHeader(404)
			utils.Render(w, "error", Err[404])
		}
		next.ServeHTTP(w, r)
	})
}
