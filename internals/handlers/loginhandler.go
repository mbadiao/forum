package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func FileService(str string, w http.ResponseWriter, data any) {
	tmpl, err := template.ParseFiles("./web/templates/" + str)
	if err != nil {
		fmt.Println("error why parsing the indicated templates")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println("error why executing the template")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/login" {
		if r.Method == "GET"  {
			FileService("login.html", w, nil)
		}/*else if r.Method == "POST"{
			
		 }*/
	} else {
		FileService("error.html", w, nil)
	}
}

