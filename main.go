package main

import (
	"fmt"
	"forum/internals/database"
	"forum/internals/handlers"
	"net/http"
	"os"
)

// type Users struct {
// 	Id           int
// 	Username     string
// 	FirstName    string
// 	LastName     string
// 	Email        string
// 	Password     string
// 	Registration string
// }

func main() {
	datas, err := database.Scan(database.CreateTable(), "Users", database.Users{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range datas {
		fmt.Println(v)
	}
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
