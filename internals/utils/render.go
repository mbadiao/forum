package utils

import (
	"html/template"
	"log"
	"net/http"
)

var TempPath = "./web/templates/"

func Render(w http.ResponseWriter, tmplName string, data interface{}) {
	tmpl, err := template.ParseFiles(TempPath + tmplName + ".html")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		log.Println("Error parsing template:", err)
		return
	}
	tmpl.Execute(w, data)
}
