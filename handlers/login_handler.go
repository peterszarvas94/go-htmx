package handlers

import (
	"go-htmx/utils"
	"net/http"
	"time"
)

// plan B

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		jwt := utils.NewToken()
		expires := time.Unix(jwt.Expires, 0)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		// set cokie
		http.SetCookie(w, &http.Cookie{
			Name:    "jwt",
			Value:   jwt.Token,
			Path:    "/",
			Expires: expires,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	http.Error(w, "Method not allowed auth", http.StatusMethodNotAllowed)
}
