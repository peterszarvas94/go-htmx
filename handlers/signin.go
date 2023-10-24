package handlers

import (
	"go-htmx/utils"
	"html"
	"html/template"
	"net/http"
	"time"
)

func getSigninTmpl() (*template.Template, error) {
	baseHtml := "templates/base.html"
	signinHtml := "templates/signin.html"
	errorHtml := "templates/error.html"

	tmpl, tmplErr := template.ParseFiles(baseHtml, signinHtml, errorHtml)
	if tmplErr != nil {
		return nil, tmplErr
	}

	return tmpl, nil
}

func SigninPage(w http.ResponseWriter, r *http.Request, pattern string) {
	session := utils.CheckSession(r)
	if session.LoggedIn {
		utils.Log(utils.INFO, "signin/checkSession", "Already logged in, redirecting to index")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	utils.Log(utils.INFO, "signin/checkSession", "Not logged in")

	tmpl, tmplErr := getSigninTmpl()
	if tmplErr != nil {
		utils.Log(utils.ERROR, "signin/signinTmpl", tmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	utils.Log(utils.INFO, "signin/signinTmpl", "Template parsed successfully")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	resErr := tmpl.Execute(w, nil)
	if resErr != nil {
		utils.Log(utils.ERROR, "signin/get/res", resErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	utils.Log(utils.INFO, "signin/get/res", "Template rendered successfully")
	return
}

func Signin(w http.ResponseWriter, r *http.Request, pattern string) {
	formErr := r.ParseForm()
	if formErr != nil {
		utils.Log(utils.ERROR, "signin/parse", formErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.Log(utils.INFO, "signin/parse", "Form parsed successfully")

	user := html.EscapeString(r.FormValue("user"))
	password := html.EscapeString(r.FormValue("password"))

	userData, userErr := utils.LoginUser(user, password)
	if userErr == nil {
		jwt, jwtErr := utils.NewToken(userData.Id)
		if jwtErr != nil {
			utils.Log(utils.ERROR, "signin/jwt", jwtErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		utils.Log(utils.INFO, "signin/jwt", "JWT created successfully")

		expires := time.Unix(jwt.Expires, 0)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    jwt.Token,
			Path:     "/",
			Expires:  expires,
			HttpOnly: true,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)

		utils.Log(utils.INFO, "signin/res", "Redirected to /")
		return
	}

	utils.Log(utils.WARNING, "signin/user", userErr.Error())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)

	signinData := utils.SigninData{
		User:  user,
		Error: "Invalid username or password",
	}

	tmpl, tmplErr := getSigninTmpl()
	if tmplErr != nil {
		utils.Log(utils.ERROR, "signin/signinTmpl", tmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	utils.Log(utils.INFO, "signin/signinTmpl", "Template parsed successfully")

	resErr := tmpl.Execute(w, signinData)
	if resErr != nil {
		utils.Log(utils.ERROR, "signin/res", resErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	utils.Log(utils.INFO, "signin/res", "Template rendered successfully")
	return
}
