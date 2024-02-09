package handlers

import "net/http"

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		if r.Method == "GET" || r.Method == "POST" {
			FileService("/home.html", w, nil)
		}
	}
}
