package handlers

import (
	"go-htmx/utils"
	"html"
	"html/template"
	"net/http"
	"net/mail"
)

func CheckUserHandler(w http.ResponseWriter, r *http.Request, pattern string) {
	// If the user is logged in, redirect them to the home page
	session := utils.CheckSession(r)
	if session.LoggedIn {
		utils.Log(utils.ERROR, "check/session", "User is already logged in")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	utils.Log(utils.INFO, "check/session", "User is not logged in")

	// Parse the "correct" and "incorrect" templates
	correct := "templates/correct.html"
	correctTmpl, correctTmplErr := template.ParseFiles(correct)
	if correctTmplErr != nil {
		utils.Log(utils.ERROR, "check/correct", correctTmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.Log(utils.INFO, "check/correct", "Template parsed successfully")

	incorrect := "templates/incorrect.html"
	incorrectTmpl, incorrectTmplErr := template.ParseFiles(incorrect)
	if incorrectTmplErr != nil {
		utils.Log(utils.ERROR, "check/incorrect", incorrectTmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.Log(utils.INFO, "check/incorrect", "Template parsed successfully")

	// Parse the form
	parseErr := r.ParseForm()
	if parseErr != nil {
		utils.Log(utils.ERROR, "check/parse", parseErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.Log(utils.INFO, "check/parse", "Form parsed successfully")

	// Check if username is taken
	username := html.EscapeString(r.FormValue("username"))

	if username != "" {
		_, userErr := utils.GetUserByUsername(username)
		if userErr != nil {
			utils.Log(utils.INFO, "check/usename", userErr.Error())

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			resErr := correctTmpl.Execute(w, nil)
			if resErr != nil {
				utils.Log(utils.ERROR, "check/username/correct", resErr.Error())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}

			return
		}

		utils.Log(utils.WARNING, "check/username", "User exists")

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		resErr := incorrectTmpl.Execute(w, nil)
		if resErr != nil {
			utils.Log(utils.ERROR, "check/username/incorrect", resErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		return
	}

	// Check if email is taken
	email := html.EscapeString(r.FormValue("email"))

	if email != "" {
		_, emailParseErr := mail.ParseAddress(email)
		_, userErr := utils.GetUserByEmail(email)
		if userErr != nil && emailParseErr == nil {
			utils.Log(utils.INFO, "check/email", userErr.Error())

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			resErr := correctTmpl.Execute(w, nil)
			if resErr != nil {
				utils.Log(utils.ERROR, "check/email/correct", resErr.Error())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}

			return
		}

		if emailParseErr != nil {
			utils.Log(utils.WARNING, "check/email", emailParseErr.Error())
		}

		if userErr == nil {
			utils.Log(utils.WARNING, "check/email", "User already exists")
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		resErr := incorrectTmpl.Execute(w, nil)
		if resErr != nil {
			utils.Log(utils.ERROR, "check/email/incorrect", resErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		utils.Log(utils.INFO, "check/email/incorrect", "Template rendered successfully")
		return
	}
}
