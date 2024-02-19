package handlers

import (
	"database/sql"
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"net/http"
)

func Getpost(r *http.Request, db *sql.DB) []database.Post {
	var Posts []database.Post
	if r.Method == "GET" {
		query := "SELECT post_id, user_id, title, PhotoURL, content, creation_date FROM Posts ORDER BY creation_date DESC"
		rows, err := db.Query(query)
		if err != nil {
			fmt.Println(err)
			return []database.Post{}
		}
		defer rows.Close()
		return ProcessPostRows(rows, db, Posts)
	}
	return []database.Post{}
}

func ProcessPostRows(rows *sql.Rows, db *sql.DB, Posts []database.Post) []database.Post {
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
		nbrcomments,_:=CountCommentsByPostID(db,post.PostID)
		post.Nbrcomments=nbrcomments

		Posts = append(Posts, post)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return []database.Post{}
	}
	return Posts
}

func GetUser(db *sql.DB) ([]database.User, error) {
	var Users []database.User
	query := "SELECT user_id, username, firstname, lastname, email, password_hash, registration_date FROM Users"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var CurrentUser database.User
		if err := rows.Scan(&CurrentUser.UserID, &CurrentUser.Username, &CurrentUser.Firstname, &CurrentUser.Lastname, &CurrentUser.Email, &CurrentUser.PasswordHash, &CurrentUser.RegistrationDate); err != nil {
			return nil, err
		}
		Users = append(Users, CurrentUser)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return Users, nil
}

func getAll(r *http.Request) (AllData, error) {
	Posts := GetPostsWithUser(r, db)
	DATA := AllData{
		Posts: Posts,
	}
	return DATA, nil
}

func GetPostsWithUser(r *http.Request, db *sql.DB) []PostWithUser {
	var postsWithUser []PostWithUser
	posts := Getpost(r, db)
	for _, post := range posts {
		// Fetch user for each post
		user, err := GetUserByID(db, post.UserID)
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

func GetUserByID(db *sql.DB, userID int) (database.User, error) {
	var user database.User
	query := "SELECT user_id, username, firstname, lastname, email, password_hash, registration_date FROM Users WHERE user_id = ?"
	err := db.QueryRow(query, userID).Scan(&user.UserID, &user.Username, &user.Firstname, &user.Lastname, &user.Email, &user.PasswordHash, &user.RegistrationDate)
	if err != nil {
		return user, err
	}
	return user, nil
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
		user, err := GetUserByID(db, post.UserID)
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
	return ProcessPostRows(rows, db, Posts)
}
