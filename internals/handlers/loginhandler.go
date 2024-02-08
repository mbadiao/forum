package handlers

import (
	"forum/internals/utils"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	utils.Render(w, "login", nil)
}
