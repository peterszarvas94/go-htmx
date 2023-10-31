package handlers

import (
	"go-htmx/utils"
	"html"
	"html/template"
	"net/http"
	"time"
	// "time"
)

/*
getSigninTmpl helper function to parse the signin template.
*/
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

/*
SigninPageHandler handles the GET request to /signin.
*/
func SigninPageHandler(w http.ResponseWriter, r *http.Request, pattern string) {
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

/*
SigninHandler handles the POST request to /signin.
*/
func SigninHandler(w http.ResponseWriter, r *http.Request, pattern string) {
	formErr := r.ParseForm()
	if formErr != nil {
		utils.Log(utils.ERROR, "signin/post/parse", formErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.Log(utils.INFO, "signin/post/parse", "Form parsed successfully")

	user := html.EscapeString(r.FormValue("user"))
	password := html.EscapeString(r.FormValue("password"))

	userData, userErr := utils.LoginUser(user, password)
	if userErr == nil {
		accessToken, accessTokenErr := utils.NewToken(userData.Id, utils.Access)
		if accessTokenErr != nil {
			utils.Log(utils.ERROR, "signin/post/access", accessTokenErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		utils.Log(utils.INFO, "signin/post/access", "Access token created successfully")

		refreshToken, refreshTokenErr := utils.NewToken(userData.Id, utils.Refresh)
		if refreshTokenErr != nil {
			utils.Log(utils.ERROR, "signin/post/refresh", refreshTokenErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		utils.Log(utils.INFO, "signin/post/refresh", "Refresh token created successfully")

		expires := time.Unix(refreshToken.Expires, 0)

		// set the refresh token cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh",
			Value:    refreshToken.Token,
			Path:     "/refresh",
			Expires:  expires,
			HttpOnly: true,
			// Secure:   true,
			// SameSite: http.SameSiteLaxMode,
		})

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		// set the access token header
		w.Header().Set("Authorization", "Bearer "+accessToken.Token)

		// for client-side redirection
		w.Header().Set("HX-Redirect", "/")

		utils.Log(utils.INFO, "signin/post/res", "Redirected to /")
		return
	}

	utils.Log(utils.WARNING, "signin/post/errorUser", userErr.Error())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)

	signinData := utils.SigninData{
		User:  user,
		Error: "Invalid username or password",
	}

	tmpl, tmplErr := getSigninTmpl()
	if tmplErr != nil {
		utils.Log(utils.ERROR, "signin/post/errorTmpl", tmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	utils.Log(utils.INFO, "signin/post/errorTmpl", "Template parsed successfully")

	resErr := tmpl.Execute(w, signinData)
	if resErr != nil {
		utils.Log(utils.ERROR, "signin/post/errorRes", resErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	utils.Log(utils.INFO, "signin/post/errorRes", "Template rendered successfully")
	return
}
