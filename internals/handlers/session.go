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

	if r.URL.Path == "/" && r.Method == "GET" {
		found := false
		ActualCookie := GetCookieHandler(w, r)
		if ActualCookie == "" {
			Post := Getpost(r, db)
			donnees := Data{
				UserPost: Post,
			}
			fmt.Println("aza")
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
				
				Post := Getpost(r, db)
				
				donnees := Data{
					Status:     "logout",
					ActualUser: CurrentUser,
					UserPost:   Post,
				}
				
				fmt.Println("aze")
				utils.FileService("home.html", w, donnees)

				return
			}
		}
		if !found {
			fmt.Println("test01")
			Post := Getpost(r, db)
			donnees:=Data{
				UserPost:Post ,
			}
			fmt.Println("aze")
			utils.FileService("home.html", w, donnees)
			return
		}
	}else if r.URL.Path == "/" && r.Method == "POST"{
		ActualCookie := GetCookieHandler(w, r)
		datas, err := database.Scan(db, "SELECT * FROM SESSIONS ", &database.Session{})
		if err != nil {
			fmt.Println("data")
			fmt.Println(err.Error())
			return
		}
		for _, data := range datas {
			u := data.(*database.Session)
			if u.Cookie_value == ActualCookie {
				CurrentUser := database.User{}
				query := "SELECT user_id, username, firstname, lastname, email, password_hash, registration_date FROM Users WHERE user_id=?"
				err := db.QueryRow(query, u.UserID).Scan(&CurrentUser.UserID, &CurrentUser.Username, &CurrentUser.Firstname, &CurrentUser.Lastname, &CurrentUser.Email, &CurrentUser.PasswordHash, &CurrentUser.RegistrationDate)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				PostHandler(w,r,CurrentUser)
			}else{
				utils.FileService("login.html",w,nil)
			}
		}
	}
}

// func Getpost(r *http.Request, db *sql.DB) []database.Post {
// 	var Posts []database.Post
// 	if r.Method == "GET" {
// 		query := "SELECT post_id, user_id, title, PhotoURL, content, creation_date FROM Posts ORDER BY creation_date DESC"
// 		if data, err := database.Scan(db, query, &database.Post{}); err == nil {
// 			for _, item := range data {
// 				if post, ok := item.(*database.Post); ok {
// 					categories, err1 := GetPostCategories(db, int(post.PostID))
// 					post.Categories = categories
// 					if err1 != nil {
// 						fmt.Println(err)
// 						return []database.Post{}
// 					}
// 					Posts = append(Posts, *post)
// 				}
// 			}
// 		} else {
// 			fmt.Println(err)
// 			return []database.Post{}
// 		}
		
// 		return Posts
// 	}
// 	return []database.Post{}
// }

// func Getpost(r *http.Request, db *sql.DB) []database.Post {
//     var Posts []database.Post
//     if r.Method == "GET" {
//         query := "SELECT post_id, user_id, title, PhotoURL, content, creation_date FROM Posts ORDER BY creation_date DESC"
//         rows, err := db.Query(query)
//         if err != nil {
//             fmt.Println(err)
//             return []database.Post{}
//         }
//         defer rows.Close()

//         for rows.Next() {
//             var post database.Post
//             if err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.PhotoURL, &post.Content, &post.CreationDate); err != nil {
//                 fmt.Println(err)
//                 return []database.Post{}
//             }
//             // Vous pouvez ajouter d'autres opérations ici si nécessaire, comme GetPostCategories
//             Posts = append(Posts, post)
//         }
//         if err := rows.Err(); err != nil {
//             fmt.Println(err)
//             return []database.Post{}
//         }
//         return Posts
//     }
//     return []database.Post{}
// }

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


