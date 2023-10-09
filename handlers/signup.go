package handlers

import (
	"go-htmx/utils"
	"html"
	"html/template"
	"net/http"
	"net/mail"
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
	incorrectHtml := "templates/incorrect.html"
	correctHtml := "templates/correct.html"

	signupTmpl, signupTmplErr := template.ParseFiles(baseHtml, signupHtml, errorHtml, incorrectHtml, correctHtml)
	if signupTmplErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	// signup page
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		resErr := signupTmpl.Execute(w, nil)
		if resErr != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		return
	}

	// signup new user
	if r.Method == "POST" {
		formErr := r.ParseForm()
		if formErr != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		username := html.EscapeString(r.FormValue("username"))
		email := html.EscapeString(r.FormValue("email"))
		password := html.EscapeString(r.FormValue("password"))

		_, userErr := utils.AddUser(username, email, password)
		if userErr == nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		errorMsg := ""

		_, invalid := mail.ParseAddress(email)
		if invalid != nil {
			errorMsg = "Invalid email"
		}

		var usernameExists bool
		_, usernameQueryErr := utils.GetUserByUsername(username)
		if usernameQueryErr != nil {
			usernameExists = false
		} else {
			usernameExists = true
			if errorMsg == "" {
				errorMsg = "Username"
			}
		}

		var emailExists bool
		_, emailQueryErr := utils.GetUserByEmail(email)
		if emailQueryErr != nil {
			emailExists = false
		} else {
			emailExists = true
			if errorMsg == "" {
				errorMsg = "Email"
			}
			if errorMsg == "Username" {
				errorMsg += " and email"
			}
		}

		if errorMsg != "Invalid email" {
			errorMsg += " already exists"
		}

		if errorMsg == "" {
			errorMsg = "Something went wrong"
		}

		signupData := utils.SignupData{
			Username: username,
			Email:    email,
			Error:    errorMsg,
			Exists: utils.ExistsData{
				Username: usernameExists,
				Email:    emailExists,
			},
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusConflict)

		res_err := signupTmpl.Execute(w, signupData)
		if res_err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		return
	}

	http.Error(w, "Method not allowed auth", http.StatusMethodNotAllowed)
}
