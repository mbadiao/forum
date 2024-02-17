package handlers

import (
	"database/sql"
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"net/http"
)

func FilterHandler(w http.ResponseWriter, r *http.Request, CurrentUser database.User) {
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

		Isconnected := utils.Isconnected(CurrentUser)

		categorypost, createdlikedpost, foundAll := utils.SplitFilter(checkboxfilter)

		query, noquery := utils.QueryFilter(categorypost, createdlikedpost, foundAll, Isconnected, CurrentUser)

		if noquery == "err" {
			data := Data{
				Page: "signin",
			}
			fmt.Println("filtre sans login")
			utils.FileService("login.html", w, data)
			return
		}

		fmt.Println("filter by", query)

		AllData, err1 := getAllFilter(w, r, query, categorypost)
		if len(AllData.Posts) == 0 {
			utils.FileService("error.html", w, Err[0])
			return
		}
		if err1 != nil {
			w.WriteHeader(400)
			utils.FileService("error.html", w, Err[400])
			return
		}

		var donnees Data
		if Isconnected {
			donnees = Data{
				Status:      "logout",
				ActualUser:  CurrentUser,
				Isconnected: true,
				Alldata:     AllData,
			}
		} else {
			donnees = Data{
				Status:      "login",
				Isconnected: false,
				Alldata:     AllData,
			}
		}

		utils.FileService("home.html", w, donnees)
		return
	} else {
		fmt.Println("filter method different de POST")
		w.WriteHeader(405)
		utils.FileService("error.html", w, Err[405])
		return
	}
}

func getAllFilter(w http.ResponseWriter, r *http.Request, query string, categorypost []string) (AllData, error) {
	Posts := GetFilterWithUser(w, r, db, query, categorypost)
	DATA := AllData{
		Posts: Posts,
	}
	return DATA, nil
}

func GetFilterWithUser(w http.ResponseWriter, r *http.Request, db *sql.DB, query string, categorypost []string) []PostWithUser {
	var postsWithUser []PostWithUser
	posts := Getpostbyfilter(r, db, query, categorypost)
	if posts == nil {
		return []PostWithUser{}
	}
	for _, post := range posts {
		// Fetch user for each post
		user, err := GetFilterUserByID(db, post.UserID)
		if err != nil {
			fmt.Println("Error fetching user for post:", err)
			continue
		}
		postsWithUser = append(postsWithUser, PostWithUser{
			Post: post,
			User: user,
		})
	}
	return postsWithUser
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
		post.FormatedDate = utils.FormatTimeAgo(post.CreationDate)

		liked := GetStatus(db, "liked", post.PostID, post.UserID)
		post.StatusLiked = liked
		disliked := GetStatus(db, "disliked", post.PostID, post.UserID)
		post.StatusDisliked = disliked

		post.Nbrlike = GetNbrStatus(db, "liked", post.PostID)
		post.Nbrdislike = GetNbrStatus(db, "disliked", post.PostID)

		Posts = append(Posts, post)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return []database.Post{}
	}
	return Posts
}

func GetFilterUserByID(db *sql.DB, userID int) (database.User, error) {
	var user database.User
	query := "SELECT user_id, username, firstname, lastname, email, password_hash, registration_date FROM Users WHERE user_id = ?"
	err := db.QueryRow(query, userID).Scan(&user.UserID, &user.Username, &user.Firstname, &user.Lastname, &user.Email, &user.PasswordHash, &user.RegistrationDate)
	if err != nil {
		return user, err
	}
	return user, nil
}
