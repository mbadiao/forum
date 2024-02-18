package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"net/http"
	"strconv"
)

func LikeDislikeHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		found := false
		usercorrespondance, err := getUserIDFromCookie(w, r)
		if err != nil {
			handleError(w, r, err)
			return
		}
		found = true
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
				if err1 == sql.ErrNoRows {
					fmt.Println("no rows returned")
					count = 0
				} else {
					w.WriteHeader(500)
					utils.FileService("error.html", w, Err[500])
					return
				}
			}

			if count == 0 {
				database.Insert(db, "LikesDislikes", "(post_id, user_id)", postid, usercorrespondance)
			}

			request := "SELECT liked, disliked FROM LikesDislikes WHERE user_id = ? AND post_id = ?"
			row := db.QueryRow(request, usercorrespondance, postid)
			var liked, disliked bool
			err2 := row.Scan(&liked, &disliked)
			if err2 != nil {
				if err2 == sql.ErrNoRows {
					fmt.Println("no row found")
					return
				} else {
					handleError(w, r, err2)
					return
				}
			}
			if actionLike != "" && actionDislike != "" {
				fmt.Println("c'est mort")
				handleBadRequest(w, r, errors.New("impossible"))
				return
			}
			if actionLike != "" {
				if (!liked && !disliked) || (!liked && disliked) {
					query := "UPDATE LikesDislikes SET liked = ?, disliked = ? WHERE user_id = ? AND post_id = ?"
					_, err := db.Exec(query, true, false, usercorrespondance, postid)
					if err != nil {
						fmt.Println("error to uptade like or dislike")
						return
					}
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				}

				if liked && !disliked {
					query := "UPDATE LikesDislikes SET liked = ?, disliked = ? WHERE user_id = ? AND post_id = ?"
					_, err := db.Exec(query, false, false, usercorrespondance, postid)
					if err != nil {
						fmt.Println("error to uptade like or dislike")
						return
					}
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				}

				handleBadRequest(w, r, errors.New("no action specified"))
				return
			}

			if actionDislike != "" {
				if (!liked && !disliked) || (liked && !disliked) {
					query := "UPDATE LikesDislikes SET liked = ?, disliked = ? WHERE user_id = ? AND post_id = ?"
					_, err := db.Exec(query, false, true, usercorrespondance, postid)
					if err != nil {
						fmt.Println("error to uptade like or dislike")
						return
					}
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				}
				if !liked && disliked {
					query := "UPDATE LikesDislikes SET liked = ?, disliked = ? WHERE user_id = ? AND post_id = ?"
					_, err := db.Exec(query, false, false, usercorrespondance, postid)
					if err != nil {
						fmt.Println("error to uptade like or dislike")
						return
					}
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				}
				handleBadRequest(w, r, errors.New("no action specified"))
				return
			}
			handleBadRequest(w, r, errors.New("no action specified"))
			return
		}
	}
}

func getUserIDFromCookie(w http.ResponseWriter, r *http.Request) (int, error) {
	actualCookie := GetCookieHandler(w, r)
	if actualCookie == "" {
		return 0, nil
	}

	var userID int
	err := db.QueryRow("SELECT user_id FROM Sessions WHERE cookie_value = ?", actualCookie).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return userID, nil
}

func handleBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(400)
	utils.FileService("error.html", w, Err[400])
	fmt.Println("Bad request:", err)
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(500)
	utils.FileService("error.html", w, Err[500])
	fmt.Println("Internal server error:", err)
}

func GetStatus(db *sql.DB, status string, postID int, userID int) string {
	var etat bool
	var etatstr string
	query := "SELECT " + status + " FROM LikesDislikes WHERE post_id = ? AND user_id = ?"
	err := db.QueryRow(query, postID, userID).Scan(&etat)
	if err != nil {
		if err == sql.ErrNoRows {
			etatstr = "debut"
			return etatstr
		} else {
			fmt.Println("error at getstatus function")
			fmt.Println(err.Error())
			return ""
		}
	}
	if etat {
		etatstr = "true"
	} else {
		etatstr = "false"
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
