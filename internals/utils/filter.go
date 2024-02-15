package utils

import (
	"fmt"
	"strings"
)

func CheckCategory(category []string) bool {
	Mapcategory := map[string]bool{
		"All":     true,
		"Tech":    true,
		"Actu":    true,
		"Mode":    true,
		"Sport":   true,
		"Edu":     true,
		"Like":    true,
		"Created": true,
	}
	found := true
	for _, v := range category {
		if !CompareCategorys(Mapcategory, v) {
			found = false
		}
	}
	return found
}

func CompareCategorys(categoriesMap map[string]bool, categoryToCheck string) bool {
	_, found := categoriesMap[categoryToCheck]
	return found
}

func QueryFilter(categorypost, createdlikedpost []string, foundAll bool, Isconnected bool) (string, string) {
	query := ""
	FoundQuery := CheckQuery(categorypost, createdlikedpost, foundAll)

	if FoundQuery == "category" {
		query = "SELECT * FROM Posts p INNER JOIN PostCategories pc ON p.post_id = pc.post_id INNER JOIN Categories c ON pc.category_id = c.category_id WHERE c.name IN ("
		placeholders := make([]string, len(categorypost))
		for i := range categorypost {
			placeholders[i] = "?"
		}
		query += strings.Join(placeholders, ",") + ")"
	} else if FoundQuery == "createlike" {
		if Isconnected {
			query = "SELECT DISTINCT p.* FROM Posts p LEFT JOIN LikesDislikes ld ON p.post_id = ld.post_id LEFT JOIN Users u ON p.user_id = u.user_id WHERE p.user_id = <id_utilisateur> OR ld.user_id = <id_utilisateur> ORDER BY p.creation_date DESC;"
		} else {
			return "", "err"
		}
	} else if FoundQuery == "likecategory" {
		if Isconnected {
			query = "SELECT DISTINCT p.* FROM Posts p INNER JOIN PostCategories pc ON p.post_id = pc.post_id INNER JOIN Categories c ON pc.category_id = c.category_id LEFT JOIN LikesDislikes ld ON p.post_id = ld.post_id WHERE c.name IN ("
			placeholders := make([]string, len(categorypost))
			for i := range categorypost {
				placeholders[i] = "?"
			}
			query += strings.Join(placeholders, ",") + ")"
			query += " AND (p.user_id = ? OR ld.user_id = ?) AND ld.like_dislike_type = 'like'"
		} else {
			return "", "err"
		}
	} else if FoundQuery == "like" {
		if Isconnected {
			query = "SELECT DISTINCT p.* FROM Posts p JOIN LikesDislikes ld ON p.post_id = ld.post_id WHERE ld.user_id = <id_utilisateur> AND ld.like_dislike_type = 'like';"
		} else {
			return "", "err"
		}
	} else if FoundQuery == "create" {
		if Isconnected {
			query = "SELECT * FROM Posts WHERE user_id = <id_utilisateur>;"
		} else {
			return "", "err"
		}
	} else {
		query = "SELECT post_id, user_id, title, PhotoURL, content, creation_date FROM Posts ORDER BY creation_date DESC"
	}

	return query, ""
}

func SplitFilter(checkedvalue []string) ([]string, []string, bool) {
	var categorypost []string
	var createdlikedpost []string
	foundAll := false
	for _, v := range checkedvalue {
		if v == "All" {
			foundAll = true
		} else if v == "Like" || v == "Created" {
			createdlikedpost = append(createdlikedpost, v)
		} else {
			categorypost = append(categorypost, v)
		}
	}
	return categorypost, createdlikedpost, foundAll
}

func CheckQuery(categorypost, createdlikedpost []string, foundAll bool) string {
	if foundAll {
		// fmt.Println("trie all")
		return "all"
	} else {
		if len(categorypost) == 0 {
			if len(createdlikedpost) == 2 {
				fmt.Println("trie sur post creer et liké")
				return "createlike"
			} else if createdlikedpost[0] == "Like" {
				fmt.Println("trie sur like")
				return "like"
			} else {
				fmt.Println("trie sur creer")
				return "create"
			}
		} else if len(createdlikedpost) == 0 {
			// fmt.Println("trie sur post category")
			return "category"
		} else {
			// fmt.Println("trie sur post creer ou liké et category")
			return "likecategory"
		}
	}
}
