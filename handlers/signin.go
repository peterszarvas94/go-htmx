package handlers

import (
	"go-htmx/utils"
	"html"
	"html/template"
	"net/http"
	"time"
)

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	session := utils.CheckSession(r)
	if session.LoggedIn {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	baseHtml := "templates/base.html"
	signinHtml := "templates/signin.html"
	errorHtml := "templates/error.html"
	tmpl, tmpl_err := template.ParseFiles(baseHtml, signinHtml, errorHtml)
	if tmpl_err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		res_err := tmpl.Execute(w, nil)
		if res_err != nil {
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

		user := html.EscapeString(r.FormValue("user"))
		password := html.EscapeString(r.FormValue("password"))

		userData, userErr := utils.LoginUser(user, password)
		if userErr == nil {
			jwt, jwtErr := utils.NewToken(userData.Id)
			if jwtErr != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			expires := time.Unix(jwt.Expires, 0)

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")

			http.SetCookie(w, &http.Cookie{
				Name:    "jwt",
				Value:   jwt.Token,
				Path:    "/",
				Expires: expires,
				HttpOnly: true,
			})

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)

		signinData := utils.SigninData{

			User:  user,
			Error: "Invalid username or password",
		}

		resErr := tmpl.Execute(w, signinData)
		if resErr != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		return

	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
