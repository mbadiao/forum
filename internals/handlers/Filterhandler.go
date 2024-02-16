package handlers

import (
	"database/sql"
	"fmt"
	"forum/internals/database"
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

		query, err := utils.QueryFilter(categorypost, createdlikedpost, foundAll, Isconnected)

		if err == "err" {
			data := Data{
				Page: "signin",
			}

			fmt.Println("filtre sans login")

			utils.FileService("login.html", w, data)
			return
		}

		fmt.Println("filter by", query)
		post := Getpostbyfilter(r, db, query, categorypost)
		if post == nil {
			utils.FileService("error.html", w, Err[0])
			return
		} else {
			utils.FileService("home.html", w, post)
			return
		}

	} else {
		fmt.Println("filter method different de POST")
		w.WriteHeader(405)
		utils.FileService("error.html", w, Err[405])
		return
	}
}

func Getpostbyfilter(r *http.Request, db *sql.DB, query string, categorypost []string) []database.Post {
    var Posts []database.Post
    var categoryInterfaces []interface{}
    for _, cat := range categorypost {
        categoryInterfaces = append(categoryInterfaces, cat)
    }

    rows, err := db.Query(query, categoryInterfaces...)
    if err != nil {
        fmt.Println(err)
        return []database.Post{}
    }
    defer rows.Close()
    for rows.Next() {
        var post database.Post
        if err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.PhotoURL, &post.Content, &post.CreationDate); err != nil {
            fmt.Println(err)
            return []database.Post{}
        }
        categories, err := GetPostCategories(db, post.PostID)
        if err != nil {
            fmt.Println(err)
            return []database.Post{}
        }
        post.Categories = categories
        Posts = append(Posts, post)
    }
    if err := rows.Err(); err != nil {
        fmt.Println(err)
        return []database.Post{}
    }
    fmt.Println("posts", Posts)
    return Posts
}


