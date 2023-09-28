package handlers

import (
	"fmt"
	"go-htmx/utils"
	"html/template"
	"net/http"
	"os"
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

	var LoggedIn bool
	if token == "" {
		LoggedIn = false
	} else {
		_, err := utils.ValidateToken(token)
		if err != nil {
			LoggedIn = false
		} else {
			LoggedIn = true
		}
	}

	base := "templates/base.html"
	index := "templates/index.html"
	todo := "templates/todo.html"
	todos := "templates/todos.html"
	tmpl, err := template.ParseFiles(base, index, todos, todo)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	data := utils.Data{
		Todos: utils.GetTodos(),
		LoggedIn: LoggedIn,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
