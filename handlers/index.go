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

	session := utils.CheckSession(r)

	baseHtml := "templates/base.html"
	indexHtml := "templates/index.html"
	todoHtml := "templates/todo.html"
	todosHtml := "templates/todos.html"

	tmpl, tmplErr := template.ParseFiles(baseHtml, indexHtml, todosHtml, todoHtml)
	if tmplErr != nil {
		http.Error(w, "Intenal server error at tmpl", http.StatusInternalServerError)
		return
	}

	todos, todosErr := utils.GetTodos()
	if todosErr != nil {
		http.Error(w, "Internal server error at todos", http.StatusInternalServerError)
		return
	}

	indexData := utils.IndexData{
		Session: session,
		Todos:    todos,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	resErr := tmpl.Execute(w, indexData)
	if resErr != nil {
		http.Error(w, "Internal server error at res", http.StatusInternalServerError)
	}

}
