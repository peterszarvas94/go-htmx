package handlers

import (
	"go-htmx/utils"
	"net/http"
	"time"
)

func Signout(w http.ResponseWriter, r *http.Request, pattern string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh",
		Value:    "",
		Path:     "/refresh",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)

	utils.Log(utils.INFO, "signout/method", "User logged out and redirected to home page")
	return
}
