package handlers

import (
	"forum/internals/utils"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	utils.Render(w, "home", nil)
}
