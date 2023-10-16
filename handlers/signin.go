package handlers

import (
	"fmt"
	"go-htmx/utils"
	"html"
	"html/template"
	"net/http"
	"time"
)

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	session := utils.CheckSession(r)
	if session.LoggedIn {
		utils.Log(utils.ERROR, "signin/session", "User is already logged in")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	utils.Log(utils.INFO, "signin/session", "User is not logged in")

	baseHtml := "templates/base.html"
	signinHtml := "templates/signin.html"
	errorHtml := "templates/error.html"
	tmpl, tmpl_err := template.ParseFiles(baseHtml, signinHtml, errorHtml)
	if tmpl_err != nil {
		utils.Log(utils.ERROR, "signin/tmpl", tmpl_err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	utils.Log(utils.INFO, "signin/tmpl", "Template parsed successfully")

	if r.Method == "GET" {
		utils.Log(utils.INFO, "signin/method", "Method is GET")

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		res_err := tmpl.Execute(w, nil)
		if res_err != nil {
			utils.Log(utils.ERROR, "signin/res", res_err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		utils.Log(utils.INFO, "signin/res", "Template rendered successfully")
		return
	}

	if r.Method == "POST" {
		utils.Log(utils.INFO, "signin/method", "Method is POST")

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
				Name:    "jwt",
				Value:   jwt.Token,
				Path:    "/",
				Expires: expires,
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

		resErr := tmpl.Execute(w, signinData)
		if resErr != nil {
			utils.Log(utils.ERROR, "signin/res", resErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		utils.Log(utils.INFO, "signin/res", "Template rendered successfully")
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	message := fmt.Sprintf("Method %s not allowed", r.Method)
	utils.Log(utils.ERROR, "signin/method", message)
}
