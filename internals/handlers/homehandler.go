package handlers

import (
	"forum/internals/database"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	db := database.CreateTable()
	CookieHandler(w, r, db)
}
