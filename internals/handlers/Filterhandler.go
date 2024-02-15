package handlers

import (
	"fmt"
	"forum/internals/utils"
	"net/http"
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		checkboxfilter := r.Form["Category"]
		if len(checkboxfilter) == 0 {
			fmt.Println("Empty checkbox filter")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		if !utils.CheckCategory(checkboxfilter) {
			fmt.Println("Bad request: Invalid category")
			w.WriteHeader(400)
			utils.FileService("error.html", w, Err[400])
			return
		}

		// cette partie correspond au session
		Isconnected := false
		categorypost, createdlikedpost, foundAll := utils.SplitFilter(checkboxfilter)


		query,err := utils.QueryFilter(categorypost, createdlikedpost, foundAll, Isconnected)

		if err == "err" {
			// page login
			fmt.Println("trouve toi un compte")
			http.Redirect(w, r, "./login", http.StatusSeeOther)
			return
		}

		fmt.Println("filter by", query)
	} else {
		fmt.Println("filter method different de POST")
		w.WriteHeader(405)
		utils.FileService("error.html", w, Err[405])
		return
	}
}
