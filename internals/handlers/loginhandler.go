package handlers

import (
	"fmt"
	"forum/internals/database"
	"html/template"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
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
		if r.Method == "GET" {
			FileService("login.html", w, nil)
		} else if r.Method == "POST" {
			db := database.CreateTable()
			firstname, err1 := IsEmpty(r.FormValue("firstname"))
			lastname, err2 := IsEmpty(r.FormValue("lastname"))
			username, err3 := IsEmpty(r.FormValue("username"))
			email, err4 := IsEmpty(r.FormValue("signup-email"))
			password, err5 := IsEmpty(r.FormValue("signup-password"))

			if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil{
				FileService("/login", w, nil)
				return
			} else {
				var nbremail, nbrusername int
				err := db.QueryRow("SELECT COUNT(*) FROM Users WHERE email=?", email).Scan(&nbremail)
				err1 := db.QueryRow("SELECT COUNT(*) FROM Users WHERE username=?", username).Scan(&nbrusername)
				http.Error(w, err1.Error(), http.StatusInternalServerError)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				if err1 != nil || err != nil {
					http.Error(w, "database error", 500)
					return
				}
				if nbremail > 0 {
					fmt.Println("email already used")
					http.Error(w, "Email already exist", http.StatusBadRequest)
					FileService("/login", w, nil)
					return
				}
				if nbrusername > 0 {
					fmt.Println("Username already used")
					http.Error(w, "Username already exist", http.StatusBadRequest)
					FileService("/login", w, nil)
					return
				}

				hashedpassword,errr:=bcrypt.GenerateFromPassword([]byte(password), 5)
				if errr!=nil{
					fmt.Println("failed to generate password")
					return
				}
				database.Insert(db, "Users", ("?,?,?,?,?"),username, firstname, lastname, email, string(hashedpassword))
			}
		}
	} else {
		FileService("error.html", w, nil)
		return
	}
}

func IsEmpty(str string) (string, error) {
	if strings.TrimSpace(str) == "" {
		return "", fmt.Errorf("all fields must be completed")
	}
	return str, nil
}
