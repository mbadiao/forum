package main

import (
	"fmt"
	"forum/internals/handlers"
	"net/http"
	"os"
)


func main() {
	if len(os.Args) == 1 {
		fs := http.FileServer(http.Dir("./web/static"))
		http.Handle("/static/", http.StripPrefix("/static/", fs))
		for _, route := range handlers.Routes {
			http.Handle(route.Path, handlers.ErrorMiddleware(route.Handler))
		}
		fmt.Println("http://localhost" + handlers.Port)
		http.ListenAndServe(handlers.Port, nil)
	}
}
