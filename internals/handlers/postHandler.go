package handlers

import (
	"database/sql"
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Category struct {
	CategoryID   int
	CategoryName string
}

var db = database.CreateTable()
var postID int64 // Declare postID outside the if statement

func PostHandler(w http.ResponseWriter, r *http.Request, userid database.User) {
	// var Posts []database.Post
	var PostS database.Post
	// var MyPostCAt PostCAt
	if r.Method == "POST" {
		err := r.ParseMultipartForm(20 << 20)
		fmt.Println(err)
		if err != nil {
			utils.FileService("error.html", w, Err[400])
			// http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		CheckboxValues := r.Form["checkbox"]
		title := r.FormValue("title")
		PhotoURL := uploadHandler(w, r)
		fmt.Println(CheckboxValues)
		fmt.Println(title)
		if PhotoURL == "NoPhoto" {
			fmt.Println(PhotoURL)
		}
		if PhotoURL != "NoPhoto" {
			if PhotoURL == "err400" {
				utils.FileService("error.html", w, Err[400])
				return
			}
			if PhotoURL == "err500" {
				utils.FileService("error.html", w, Err[500])
				return
			} else {
				PhotoURL = PhotoURL[5:]
			}
		}
		thread := r.FormValue("thread")
		fmt.Println()
		if len(CheckboxValues) == 0 || title == "" || thread == "" {
			return
		}
		a := Checkcategory(CheckboxValues)
		if !a {
			return
		}
		PostS = database.Post{
			UserID:   userid.UserID,
			Title:    title,
			PhotoURL: PhotoURL,
			Content:  thread,
		}
		database.Insert(db, "Posts", "(user_id, title, PhotoURL, content)", PostS.UserID, PostS.Title, PostS.PhotoURL, PostS.Content)
		err = db.QueryRow("SELECT last_insert_rowid()").Scan(&postID)
		if err != nil {
			log.Fatal(err)
		}
		CategoriesId := getCategory(CheckboxValues)
		for _, v := range CategoriesId {
			database.Insert(db, "PostCategories", "(post_id, category_id)", postID, v)
		}
		fmt.Println("here")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func getCategory(CheckboxValues []string) []int {
	var CategoriesId []int
	for _, v := range CheckboxValues {
		if v == "All" {
			CategoriesId = []int{1, 2, 3, 4, 5}
		}
		if v == "Tech" {
			CategoriesId = append(CategoriesId, 1)
		}
		if v == "Actu" {
			CategoriesId = append(CategoriesId, 2)
		}
		if v == "Mode" {
			CategoriesId = append(CategoriesId, 3)
		}
		if v == "Sport" {
			CategoriesId = append(CategoriesId, 4)
		}
		if v == "Edu" {
			CategoriesId = append(CategoriesId, 5)
		}
	}
	return CategoriesId
}

func Checkcategory(category []string) bool {
	Mapcategory := map[string]bool{
		"All":   true,
		"Tech":  true,
		"Actu":  true,
		"Mode":  true,
		"Sport": true,
		"Edu":   true,
	}
	found := true
	for _, v := range category {
		if !CompareCategory(Mapcategory, v) {
			found = false
		}
	}
	return found
}

func CompareCategory(categoriesMap map[string]bool, categoryToCheck string) bool {
	_, found := categoriesMap[categoryToCheck]
	return found
}

// Function to retrieve categories of a post
func GetPostCategories(db *sql.DB, postID int) ([]string, error) {
	query := `
    SELECT PostCategories.post_id, Categories.name
    FROM PostCategories
    INNER JOIN Categories ON PostCategories.category_id = Categories.category_id
    WHERE PostCategories.post_id =` + strconv.Itoa(postID)
	// Call the Scan function with the PostCategory struct
	data, err := database.Scan(db, query, &database.Category{})
	if err != nil {
		return nil, err
	}
	// Extract category names from the result
	var categories []string
	for _, d := range data {
		category := d.(*database.Category)
		categories = append(categories, category.Name)
	}
	return categories, nil
}

func uploadHandler(w http.ResponseWriter, r *http.Request) string {
	var photoURL = "NoPhoto"
	// ParseMultipartForm analyse la demande pour extraire les fichiers
	err := r.ParseMultipartForm(20 << 20) // taille maximale autorisée pour les fichiers (20MB)
	if err != nil {
		fmt.Println("taille")
		return "err400"
	}
	// Obtenez le fichier téléchargé à partir de la clé du formulaire
	file, handler, err := r.FormFile("file")
	if file != nil {
		if err != nil {
			fmt.Println("fichier" + err.Error())
			return "err400"
		}
		defer file.Close()
		if !IsValidImage(file, handler) {
			return "err400"
		}
		// Vérifiez la taille du fichier
		if handler.Size > 20<<20 {
			fmt.Println("heres")
			return "err400"
		}
		// Créez un fichier dans le répertoire temporaire du serveur
		tempFile, err := os.CreateTemp("./web/static/upload", "upload-*"+filepath.Ext(handler.Filename))
		if err != nil {
			fmt.Println("répertoire" + err.Error())
			return "err500"
		}
		defer tempFile.Close()
		// Copiez le contenu du fichier téléchargé dans le fichier temporaire
		_, err = io.Copy(tempFile, file)
		if err != nil {
			fmt.Println("Copiez" + err.Error())
			return "err500"
		}
		// Générez une URL pour le fichier téléchargé (par exemple, l'URL peut être le chemin relatif vers le fichier temporaire)
		photoURL = tempFile.Name()
		fmt.Println("avant", photoURL)
	}
	return photoURL
}

func IsValidImage(file multipart.File, handler *multipart.FileHeader) bool {
	//contentType := handler.Header.Get("Content-Type")
	buff := make([]byte, 512)
	_, err := file.Read(buff)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return false
	}

	// Reset file offset to beginning for subsequent reads
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Println("Error seeking file:", err)
		return false
	}

	// Detect content type based on the first 512 bytes
	contentType := http.DetectContentType(buff)

	switch contentType {
	case "image/svg+xml", "image/jpeg", "image/gif", "image/png":
		// OK, continuez le traitement
	default:
		fmt.Println(contentType)
		fmt.Println("type")
		return false
	}
	return true
}
