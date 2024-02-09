package main

import (
	"fmt"
	"forum/internals/handlers"
	"net/http"
	"os"
	// "fmt"
	// "forum/internals/handlers"
	// "net/http"
	// "os"
)

func main() {
	if len(os.Args) == 1 {
		service := http.FileServer(http.Dir("./web/static"))
		http.Handle("/static/", http.StripPrefix("/static", service))
		// http.HandleFunc("/", handlers.HomeHandler)
		// http.HandleFunc("/login", handlers.LoginHandler)
		// port := ":8080"
		// fmt.Println("Server listening at http://localhost" + port)
		// fmt.Println("press Ctrl+C to disconnect")
		// http.ListenAndServe(port, nil)
		for _, route := range handlers.Routes {
			http.Handle(route.Path, handlers.ErrorMiddleware(route.Handler))
		}
		fmt.Println("http://localhost" + handlers.Port)
		fmt.Println("press Ctrl+C to disconnect")
		http.ListenAndServe(handlers.Port, nil)
	} else {
		fmt.Println("usage: go run .")
	}
}

