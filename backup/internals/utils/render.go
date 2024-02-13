package utils

import (
	"fmt"
	"html/template"
	"net/http"
)

func FileService(str string, w http.ResponseWriter, data any) {
	tmpl, err := template.ParseFiles("./web/templates/" + str)
	if err != nil {
		fmt.Println("error while parsing the indicated templates")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println("error while executing the template")
		fmt.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
