package handlers

import (
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"net/http"
	"strconv"
	"strings"
)

func DisplayComment(w http.ResponseWriter, r *http.Request) []database.Comment {
	var CommentData []database.Comment
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return nil
	}
	rows, err := db.Query("SELECT * FROM Comments WHERE post_id=? ORDER BY creation_date DESC", id)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	for rows.Next() {
		var comment database.Comment
		err = rows.Scan(&comment.CommentID, &comment.PostID, &comment.UserID, &comment.Content, &comment.Username, &comment.CreationDate)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		CommentData = append(CommentData, comment)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return CommentData
}

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if !CheckId(id) {
		utils.FileService("error.html", w, Err[404])
		return
	}
	if err != nil {
		return
	}

	if strings.TrimSpace(r.FormValue("comment")) != "" {
		RecordComment(w, r)
	}
	utils.FileService("comment.html", w, DisplayComment(w, r))
}

func RecordComment(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("ForumCookie")
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusUnauthorized)
		return
	}
	var userId int
	err = db.QueryRow("SELECT user_id FROM Sessions WHERE cookie_value=?", cookie.Value).Scan(&userId)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		fmt.Println(err.Error())
		return
	}
	var username string
	errScanUser := db.QueryRow("SELECT userName FROM Users WHERE user_id=?", userId).Scan(&username)
	if errScanUser != nil {
		fmt.Println(errScanUser.Error())
		return
	}
	comment := database.Comment{
		PostID:   id,
		UserID:   userId,
		Username: username,
		Content:  r.FormValue("comment"),
	}
	fmt.Println("Take comment")
	database.Insert(db, "Comments", "(post_id, user_id, userName ,content)", comment.PostID, comment.UserID, comment.Username, comment.Content)
}

func CheckId(id int) bool {
	var userid int
	err := db.QueryRow("SELECT post_id FROM Posts WHERE post_id=?", id).Scan(&userid)
	return err == nil
}
