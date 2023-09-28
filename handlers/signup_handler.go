package handlers

import (
	"fmt"
	"go-htmx/utils"
	"html/template"
	"net/http"
	"os"
	"strings"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		base := "templates/base.html"
		signup := "templates/signup.html"
		error := "templates/error.html"
		tmpl, err := template.ParseFiles(base, signup, error)
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
		r.ParseForm()
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		userAdded, userErr := utils.NewUser(username, email, password)
		if userAdded {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		errorMsg := "Internal server error"

		if strings.Contains(userErr.Error(), "UNIQUE constraint failed: users.username") {
			errorMsg = "Username is taken"
		}

		if strings.Contains(userErr.Error(), "UNIQUE constraint failed: users.email") {
			errorMsg = "User with this email already exists"
		}

		base := "templates/base.html"
		signup := "templates/signup.html"
		error := "templates/error.html"
		tmpl, tmpl_err := template.ParseFiles(base, signup, error)
		if tmpl_err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		res_err := tmpl.Execute(w, errorMsg)
		if res_err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		return
	}

	http.Error(w, "Method not allowed auth", http.StatusMethodNotAllowed)
}
