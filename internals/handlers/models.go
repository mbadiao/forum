package handlers

import (
	"net/http"
)

type Error struct {
	Code    int
	Message string
}

var Err = map[int]Error{
	404: {
		http.StatusNotFound,
		http.StatusText(404),
	},
	500: {
		http.StatusInternalServerError,
		http.StatusText(500),
	},
	400: {
		http.StatusBadRequest,
		http.StatusText(400),
	},
}

type Route struct {
	Path    string
	Handler http.HandlerFunc
	Method  []string
}

var Routes = []Route{
	{
		Path:    "/login",
		Handler: LoginHandler,
		Method:  []string{"GET"},
	},
	{
		Path:    "/",
		Handler: HomeHandler,
		Method:  []string{"GET", "POST"},
	},
	// {
	// 	Path:    "/",
	// 	Handler: PostHandler,
	// 	Method:  []string{"GET", "POST"},
	// },
	{
		Path:    "/logout",
		Handler: LogoutHandler,
		Method:  []string{"GET", "POST"},
	},
}

var Port = ":8080"
