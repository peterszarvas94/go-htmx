package handlers

import (
	"fmt"
	"go-htmx/utils"
	"html/template"
	"net/http"
	"os"
	"time"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		base := "templates/base.html"
		auth := "templates/auth.html"
		tmpl, err := template.ParseFiles(base, auth)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	if r.Method == "POST" {
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
