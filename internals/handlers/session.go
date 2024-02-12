package handlers

import (
	"database/sql"
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gofrs/uuid/v5"
)

func CreateCookie(w http.ResponseWriter) http.Cookie {
	Tokens, _ := uuid.NewV4()

	now := time.Now()
	expires := now.Add(time.Hour * 1)

	cookie := http.Cookie{
		Name:     "ForumCookie",
		Value:    Tokens.String(),
		Expires:  expires,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
	return cookie
}

func GetCookieHandler(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("ForumCookie")
	if err != nil {
		return ""
	}
	return (cookie.Value)
}

func CookieHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	donnees := Data{
		Status: "Logout",
	}
	if r.URL.Path == "/" && r.Method == "GET" {
		found := false
		ActualCookie := GetCookieHandler(w, r)
		if ActualCookie == "" {
			utils.FileService("home.html", w, nil)
			return
		}
		datas, err := database.Scan(db, "SELECT * FROM SESSIONS ", &database.Session{})
		if err != nil {
			fmt.Println("data")
			fmt.Println(err.Error())
			return
		}
		for _, data := range datas {
			u := data.(*database.Session)
			if u.Cookie_value == ActualCookie {
				donnees.Name = strconv.Itoa(u.UserID)
				utils.FileService("home.html", w, donnees)
				found = true
				return
			}
		}
		if !found {
			fmt.Println("test01")
			utils.FileService("home.html", w, nil)
			return
		}
	}
}
