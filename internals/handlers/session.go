package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Sessions struct {
}

func CreateCookie() {
	Tokens, err := uuid.NewV4()
	if err != nil {
		return
	}
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
	fmt.Println(cookie)
}
