package handlers

import (
	"fmt"
	"go-htmx/utils"
	"net/http"
	"time"
)

func SignoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		utils.Log(utils.INFO, "signout/method", "Method is POST")

		http.SetCookie(w, &http.Cookie{
			Name:    "jwt",
			Value:   "",
			Path:    "/",
			Expires: time.Now().Add(-1 * time.Hour),
			HttpOnly: true,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)

		utils.Log(utils.INFO, "signout/method", "User logged out and redirected to home page")
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	
	message := fmt.Sprintf("Method %s not allowed", r.Method)
	utils.Log(utils.ERROR, "signout/method", message)
}
