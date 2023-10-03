package handlers

import (
	"go-htmx/utils"
	"html/template"
	"net/http"
	"strings"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	session := utils.CheckSession(r)
	if session.LoggedIn {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	baseHtml := "templates/base.html"
	signupHtml := "templates/signup.html"
	errorHtml := "templates/error.html"
	tmpl, tmplErr := template.ParseFiles(baseHtml, signupHtml, errorHtml)
	if tmplErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		resErr := tmpl.Execute(w, nil)
		if resErr != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		return
	}

	if r.Method == "POST" {
		formErr := r.ParseForm()
		if formErr != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		_, userErr := utils.AddUser(username, email, password)
		if userErr == nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		errorMsg := "Internal server error"

		if strings.Contains(userErr.Error(), "UNIQUE constraint failed: users.email") {
			errorMsg = "User with this email already exists"
		}

		if strings.Contains(userErr.Error(), "UNIQUE constraint failed: users.username") {
			errorMsg = "Username is taken"
		}

		signupData := utils.SignupData{
			Username: username,
			Email:    email,
			Error:    errorMsg,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusConflict)
		res_err := tmpl.Execute(w, signupData)
		if res_err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		return
	}

	http.Error(w, "Method not allowed auth", http.StatusMethodNotAllowed)
}
