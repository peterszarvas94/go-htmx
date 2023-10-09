package handlers

import (
	"go-htmx/utils"
	"html"
	"html/template"
	"net/http"
)

func CheckHandler(w http.ResponseWriter, r *http.Request) {
	session := utils.CheckSession(r)
	if session.LoggedIn {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	correct := "templates/correct.html"
	correctTmpl, correctTmplErr := template.ParseFiles(correct)
	if correctTmplErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	incorrect := "templates/incorrect.html"
	incorrectTmpl, incorrectTmplErr := template.ParseFiles(incorrect)
	if incorrectTmplErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	parseErr := r.ParseForm()
	if parseErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	username := html.EscapeString(r.FormValue("username"))

	if username != "" {
		user, userErr := utils.GetUserByUsername(username)
		if userErr != nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			resErr := correctTmpl.Execute(w, nil)
			if resErr != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}

		if user.Username == username {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			resErr := incorrectTmpl.Execute(w, nil)
			if resErr != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}
	}

	email := html.EscapeString(r.FormValue("email"))

	if email != "" {
		user, userErr := utils.GetUserByEmail(email)
		if userErr != nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			resErr := correctTmpl.Execute(w, nil)
			if resErr != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}

		if user.Email == email {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			resErr := incorrectTmpl.Execute(w, nil)
			if resErr != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}
	}
}
