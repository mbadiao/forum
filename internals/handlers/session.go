package handlers

import (
	"database/sql"
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
)

type AllData struct {
	Posts []PostWithUser
}

type PostWithUser struct {
	Post database.Post
	User database.User
}

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
	var CurrentUser database.User
	if r.URL.Path == "/" && r.Method == "GET" {
		found := false
		ActualCookie := GetCookieHandler(w, r)
		if ActualCookie == "" {
			AllData, err := getAll(r)
			if err != nil {
				fmt.Println(err)
				return
			}
			donnees := Data{
				Alldata: AllData,
			}
			utils.FileService("home.html", w, donnees)
			return
		}
		datas, err := database.Scan(db, "SELECT * FROM SESSIONS ", &database.Session{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		for _, data := range datas {
			u := data.(*database.Session)
			if u.Cookie_value == ActualCookie {
				found = true
				CurrentUser := database.User{}
				query := "SELECT user_id, username, firstname, lastname, email, password_hash, registration_date FROM Users WHERE user_id=?"
				err := db.QueryRow(query, u.UserID).Scan(&CurrentUser.UserID, &CurrentUser.Username, &CurrentUser.Firstname, &CurrentUser.Lastname, &CurrentUser.Email, &CurrentUser.PasswordHash, &CurrentUser.RegistrationDate)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				AllData, err := getAll(r)
				if err != nil {
					fmt.Println(err)
					return
				}
				donnees := Data{
					Status:      "logout",
					ActualUser:  CurrentUser,
					Isconnected: true,
					Alldata:     AllData,
				}
				utils.FileService("home.html", w, donnees)
				return
			}
		}
		if !found {
			AllData, err := getAll(r)
			if err != nil {
				fmt.Println(err)
				return
			}
			donnees := Data{
				Isconnected: false,
				Alldata:     AllData,
			}
			utils.FileService("home.html", w, donnees)
			return
		}
	}
	
	if r.URL.Path == "/" && r.Method == "POST" {
		ActualCookie := GetCookieHandler(w, r)
		datas, err := database.Scan(db, "SELECT * FROM SESSIONS ", &database.Session{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		Found := false
		for _, data := range datas {
			u := data.(*database.Session)
			if u.Cookie_value == ActualCookie {
				Found = true
				CurrentUser = database.User{}
				query := "SELECT user_id, username, firstname, lastname, email, password_hash, registration_date FROM Users WHERE user_id=?"
				err := db.QueryRow(query, u.UserID).Scan(&CurrentUser.UserID, &CurrentUser.Username, &CurrentUser.Firstname, &CurrentUser.Lastname, &CurrentUser.Email, &CurrentUser.PasswordHash, &CurrentUser.RegistrationDate)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}
			PostHandler(w, r, CurrentUser)
			if !Found {
				utils.FileService("login.html", w, nil)
			}
		}
	}

	if r.URL.Path == "/filter" && r.Method == "POST" {
		ActualCookie := GetCookieHandler(w, r)
		datas, err := database.Scan(db, "SELECT * FROM SESSIONS ", &database.Session{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		Found := false
		for _, data := range datas {
			u := data.(*database.Session)
			if u.Cookie_value == ActualCookie {
				Found = true
				CurrentUser = database.User{}
				query := "SELECT user_id, username, firstname, lastname, email, password_hash, registration_date FROM Users WHERE user_id=?"
				err := db.QueryRow(query, u.UserID).Scan(&CurrentUser.UserID, &CurrentUser.Username, &CurrentUser.Firstname, &CurrentUser.Lastname, &CurrentUser.Email, &CurrentUser.PasswordHash, &CurrentUser.RegistrationDate)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}
			FilterHandler(w, r, CurrentUser)
			if !Found {
				utils.FileService("login.html", w, nil)
			}
		}
	} 
}


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
		return Posts
	}
	return []database.Post{}
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
