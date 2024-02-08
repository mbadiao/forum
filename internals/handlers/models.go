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
		"Not Found",
	},
	500: {
		http.StatusInternalServerError,
		"Internal Server Error",
	},
	400: {
		http.StatusBadRequest,
		"Bad Request",
	},
	0: {
		0,
		"result found",
	},
}

type Route struct {
	Path    string
	Handler http.HandlerFunc
	Method  []string
}

var Routes = []Route{
	{
		Path:    "login",
		Handler: LoginHandler,
		Method:  []string{"GET"},
	},
	{
		Path:    "home",
		Handler: HomeHandler,
		Method:  []string{"GET", "POST"},
	},
}

var Port = "8081"
