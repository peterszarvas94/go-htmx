package handlers

import (
	"go-htmx/utils"
	"html/template"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed index", http.StatusMethodNotAllowed)
		return
	}

	cookies := r.Cookies()
	var token string
	for _, cookie := range cookies {
		if cookie.Name == "jwt" {
			token = cookie.Value
		}
	}

	claims, jwtErr := utils.ValidateToken(token)
	userId, subErr := claims.GetSubject()
	user, dbErr := utils.GetUserById(userId)

	loggedIn := jwtErr == nil && subErr == nil && dbErr == nil

	baseHtml := "templates/base.html"
	indexHtml := "templates/index.html"
	todoHtml := "templates/todo.html"
	todosHtml := "templates/todos.html"
	tmpl, tmplErr := template.ParseFiles(baseHtml, indexHtml, todosHtml, todoHtml)
	if tmplErr != nil {
		http.Error(w, "Intenal server error", http.StatusInternalServerError)
	}

	todos, todosErr := utils.GetTodos()
	if todosErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	todosData := utils.IndexData{
		User:     user,
		Todos:    todos,
		LoggedIn: loggedIn,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	resErr := tmpl.Execute(w, todosData)
	if resErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

}
