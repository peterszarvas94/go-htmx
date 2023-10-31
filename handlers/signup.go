package handlers

import (
	"fmt"
	"go-htmx/utils"
	"html"
	"html/template"
	"net/http"
	"net/mail"
)

/*
getSignupTmpl helper function to parse the signup template.
*/
func getSignupTmpl() (*template.Template, error) {
	baseHtml := "templates/base.html"
	signupHtml := "templates/signup.html"
	errorHtml := "templates/error.html"
	incorrectHtml := "templates/incorrect.html"
	correctHtml := "templates/correct.html"

	tmpl, tmplErr := template.ParseFiles(baseHtml, signupHtml, errorHtml, incorrectHtml, correctHtml)
	if tmplErr != nil {
		utils.Log(utils.ERROR, "signup/signupTmpl", tmplErr.Error())
		return nil, tmplErr
	}

	utils.Log(utils.INFO, "signup/signupTmpl", "Template parsed successfully")
	return tmpl, nil
}

/*
getErrorMessage helper function to get the error message for username and email.
*/
func getErrorMessage(invalidEmail, usernameExists, emailExists bool) string {
	if invalidEmail {
		return "Invalid email"
	}
	if usernameExists && emailExists {
		return "Username and email already exist"
	}
	if usernameExists {
		return "Username already exists"
	}
	if emailExists {
		return "Email already exists"
	}
	return "Something went wrong"
}

/*
SignupPageHandler handles the GET request to /signup.
*/
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

/*
NewUserHandler handles the POST request to /signup.
*/
func NewUserHandler(w http.ResponseWriter, r *http.Request, pattern string) {
	formErr := r.ParseForm()
	if formErr != nil {
		utils.Log(utils.ERROR, "signup/post/form", formErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.Log(utils.INFO, "signup/post/form", "Form parsed successfully")

	username := html.EscapeString(r.FormValue("username"))
	email := html.EscapeString(r.FormValue("email"))
	password := html.EscapeString(r.FormValue("password"))

	_, userErr := utils.AddUser(username, email, password)
	if userErr == nil {
		utils.Log(utils.INFO, "signup/post/user", "User added successfully")
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	utils.Log(utils.ERROR, "signup/post/user", userErr.Error())

	_, invalid := mail.ParseAddress(email)
	if invalid != nil {
		utils.Log(utils.ERROR, "signup/post/user", invalid.Error())
	}

	_, usernameQueryErr := utils.GetUserByUsername(username)
	if usernameQueryErr != nil {
		utils.Log(utils.ERROR, "signup/post/user", usernameQueryErr.Error())
	}

	_, emailQueryErr := utils.GetUserByEmail(email)
	if emailQueryErr != nil {
		utils.Log(utils.ERROR, "signup/post/user", emailQueryErr.Error())
	}

	emailInvalid := invalid != nil
	usernameExists := usernameQueryErr == nil
	emailExists := emailQueryErr == nil

	if emailInvalid {
		logMsg := fmt.Sprintf("Invalid email: %s", email)
		utils.Log(utils.ERROR, "signup/post/user", logMsg)
	}

	if usernameExists {
		logMsg := fmt.Sprintf("Username already exists: %s", username)
		utils.Log(utils.ERROR, "signup/post/user", logMsg)
	}

	if emailExists {
		logMsg := fmt.Sprintf("Email already exists: %s", email)
		utils.Log(utils.ERROR, "signup/post/user", logMsg)
	}

	errorMsg := getErrorMessage(emailInvalid, usernameExists, emailExists)

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
		utils.Log(utils.ERROR, "signup/post/signupTmpl", tmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	utils.Log(utils.INFO, "signup/post/signupTmpl", "Template parsed successfully")

	res_err := tmpl.Execute(w, signupData)
	if res_err != nil {
		utils.Log(utils.ERROR, "signup/post/res", res_err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	utils.Log(utils.INFO, "signup/post/res", "Template rendered successfully")

	return
}
