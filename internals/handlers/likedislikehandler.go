package handlers

import (
	"database/sql"
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"net/http"
	"strconv"
)

func LikeDislikeHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		found := false
		usercorrespondance := 0
		actualcookie := GetCookieHandler(w, r)
		if actualcookie == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		} else {

			err := db.QueryRow("SELECT user_id FROM Sessions WHERE cookie_value =?", actualcookie).Scan(&usercorrespondance)
			if err != nil {
				if err == sql.ErrNoRows {
					http.Redirect(w, r, "/login", http.StatusSeeOther)
					return
				} else {
					fmt.Println("erreur at like dislike handler , with the query")
					fmt.Println(err.Error())
					utils.FileService("error.html", w, Err[500])
					return
				}

			}
			if usercorrespondance == 0 {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			} else {
				found = true
			}
		}
		if found {
			postid, err := strconv.Atoi(r.FormValue("postidouz"))
			if err != nil {
				
				w.WriteHeader(400)
				utils.FileService("error.html", w, Err[400])
				return
			}

			actionLike := r.FormValue("actionlike")
			actionDislike := r.FormValue("actiondislike")
			count := 0
			query := "SELECT COUNT(*) FROM LikesDislikes WHERE user_id = ? AND post_id = ?"
			err1 := db.QueryRow(query, usercorrespondance, postid).Scan(&count)
			if err1 != nil {
				w.WriteHeader(500)
				utils.FileService("error.html", w, Err[500])
				return
			}

			if count == 0 {
				database.Insert(db, "LikesDislikes", "(post_id, user_id)", postid, usercorrespondance)
			}

			if actionLike != "" {
				if actionLike == "false" || actionLike == "debut" {
					like := true

					err := updateLikeDislike(postid, like)
					if err != nil {
						w.WriteHeader(400)
						utils.FileService("error.html", w, Err[400])
						return
					}
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				} else if actionLike == "true" {
					like := false
					err := updateLikeOrDislike(postid, like, "liked")
					if err != nil {
						w.WriteHeader(400)
						utils.FileService("error.html", w, Err[400])
						return
					}
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				} else {
					w.WriteHeader(400)
					utils.FileService("error.html", w, Err[400])
					return
				}
			} else if actionDislike != "" {
				if actionDislike == "false" || actionDislike == "debut" {
					like := true
					err := updateLikeDislike(postid, !like)
					if err != nil {
						w.WriteHeader(400)
						utils.FileService("error.html", w, Err[400])
						return
					}
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				} else if actionDislike == "true" {
					dislike := false
					err := updateLikeOrDislike(postid, dislike, "disliked")
					if err != nil {
						w.WriteHeader(400)
						utils.FileService("error.html", w, Err[400])
						return
					}
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				} else {
					w.WriteHeader(400)
					utils.FileService("error.html", w, Err[400])
					return
				}
			} else {
				w.WriteHeader(400)
				utils.FileService("error.html", w, Err[400])
				return
			}
		}
	}
}

func updateLikeDislike(postID int, liked bool) error {
	query := "UPDATE LikesDislikes SET liked = ?, disliked = ? WHERE post_id = ?"
	_, err := db.Exec(query, liked, !liked, postID)
	if err != nil {
		return err
	}
	return nil
}

func updateLikeOrDislike(postID int, liked bool, column string) error {
	query := fmt.Sprintf("UPDATE LikesDislikes SET " + column + " = ? WHERE post_id = ?")
	_, err := db.Exec(query, liked, postID)
	if err != nil {
		return err
	}
	return nil
}

func GetStatus(db *sql.DB,status string, postID int , userID int)string{
	var etat bool
	var etatstr string
	query:="SELECT "+status+" FROM LikesDislikes WHERE post_id = ? AND user_id = ?"
	err:=db.QueryRow(query, postID, userID).Scan(&etat)
	if err != nil{
		if err == sql.ErrNoRows{
			etatstr="debut"
			return etatstr
		}else{
			fmt.Println("error at getstatus function")
			fmt.Println(err.Error())
			return ""
		}
	}
	if etat{
		etatstr="true"
	}else{
		etatstr="false"
	}
	return etatstr
}


func GetNbrStatus(db *sql.DB, status string, postID int) int {
    count := 0
    query := "SELECT COUNT(*) FROM LikesDislikes WHERE post_id = ? AND " + status + " = true"
    err := db.QueryRow(query, postID).Scan(&count)
    if err != nil {
        if err == sql.ErrNoRows {
            return 0
        } else {
            fmt.Println("Erreur lors de l'exécution de la requête dans GetNbrStatus:", err)
            return 0
        }
    }
    return count
}