package handlers

import (
	"go-htmx/utils"
	"html"
	"html/template"
	"net/http"
	"net/mail"
)

func getSignupTmpl() (*template.Template, error) {
	baseHtml := "templates/base.html"
	signupHtml := "templates/signup.html"
	errorHtml := "templates/error.html"
	incorrectHtml := "templates/incorrect.html"
	correctHtml := "templates/correct.html"

	tmpl, tmplErr := template.ParseFiles(baseHtml, signupHtml, errorHtml, incorrectHtml, correctHtml)
	if tmplErr != nil {
		return nil, tmplErr
	}

	return tmpl, nil
}

func SignupPageHandler(w http.ResponseWriter, r *http.Request, pattern string) {
	session := utils.CheckSession(r)
	if session.LoggedIn {
		utils.Log(utils.INFO, "signup/checkSession", "Already logged in, redirecting to index")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	utils.Log(utils.INFO, "signup/checkSession", "Not logged in")

	tmpl, tmplErr := getSignupTmpl()
	if tmplErr != nil {
		utils.Log(utils.ERROR, "signup/signupTmpl", tmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	utils.Log(utils.INFO, "signup/signupTmpl", "Template parsed successfully")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	resErr := tmpl.Execute(w, nil)
	if resErr != nil {
		utils.Log(utils.ERROR, "signup/get/res", resErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	utils.Log(utils.INFO, "signup/get/res", "Template rendered successfully")
	return
}

func NewUserHandler(w http.ResponseWriter, r *http.Request, pattern string) {
	formErr := r.ParseForm()
	if formErr != nil {
		utils.Log(utils.ERROR, "signup/add/form", formErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.Log(utils.INFO, "signup/add/form", "Form parsed successfully")

	username := html.EscapeString(r.FormValue("username"))
	email := html.EscapeString(r.FormValue("email"))
	password := html.EscapeString(r.FormValue("password"))

	_, userErr := utils.AddUser(username, email, password)
	if userErr == nil {
		utils.Log(utils.INFO, "signup/add/user", "User added successfully")
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	utils.Log(utils.ERROR, "signup/add/user", userErr.Error())

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

	tmpl, tmplErr := getSignupTmpl()
	if tmplErr != nil {
		utils.Log(utils.ERROR, "signup/signupTmpl", tmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	utils.Log(utils.INFO, "signup/signupTmpl", "Template parsed successfully")

	res_err := tmpl.Execute(w, signupData)
	if res_err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	return
}
