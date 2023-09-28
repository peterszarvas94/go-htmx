package handlers

import (
	"go-htmx/utils"
	"html/template"
	"net/http"
	"time"
)

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		base := "templates/base.html"
		signin := "templates/signin.html"
		error := "templates/error.html"
		tmpl, tmpl_err := template.ParseFiles(base, signin, error)
		if tmpl_err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		res_err := tmpl.Execute(w, nil)
		if res_err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		return
	}

	if r.Method == "POST" {
		r.ParseForm()
		user := r.FormValue("user")
		password := r.FormValue("password")

		exists, _ := utils.GetUserByUsernameOrEmail(user, password)
		if exists {
			jwt := utils.NewToken()
			expires := time.Unix(jwt.Expires, 0)

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")

			// set cokie
			http.SetCookie(w, &http.Cookie{
				Name:    "jwt",
				Value:   jwt.Token,
				Path:    "/",
				Expires: expires,
			})

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		base := "templates/base.html"
		signin := "templates/signin.html"
		error := "templates/error.html"
		tmpl, tmpl_err := template.ParseFiles(base, signin, error)
		if tmpl_err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)

		signinData := utils.Signin{
			User: user,
			Error: "Invalid username or password",
		}

		res_err := tmpl.Execute(w, signinData)
		if res_err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		return

	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
