package models

import "net/http"

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
}
