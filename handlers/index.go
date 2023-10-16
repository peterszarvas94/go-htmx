package handlers

import (
	"fmt"
	"go-htmx/utils"
	"html/template"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// if method is not GET, 405
	if r.Method != "GET" {
		message := fmt.Sprintf("Method %s not allowed", r.Method)

		utils.Log(utils.ERROR, "index/method", message)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	utils.Log(utils.INFO, "index/path", r.URL.Path)
	path := utils.GetPath(r)

	// if path is not empty, 404
	if len(path) > 0 {
		message := fmt.Sprintf("Path %s not found", r.URL.Path)
		utils.Log(utils.WARNING, "index/notfound/path", message)
		
		baseHtml := "templates/base.html"
		notfoundHtml := "templates/404.html"

		tmpl, tmplErr := template.ParseFiles(baseHtml, notfoundHtml)
		if tmplErr != nil {
			utils.Log(utils.ERROR, "index/notfound/tmpl", tmplErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		resErr := tmpl.Execute(w, nil)
		if resErr != nil {
			utils.Log(utils.ERROR, "index/notfound/res", resErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		utils.Log(utils.INFO, "index/notfound/res", "Template parsed successfully")
		return
	}

	session := utils.CheckSession(r)

	if session.LoggedIn {
		sessionId := fmt.Sprint(session.User.Id)
		sessionLog := fmt.Sprintf("Session found: %s", sessionId)
		utils.Log(utils.INFO, "index/checkSession", sessionLog)
	} else {
		utils.Log(utils.INFO, "index/checkSession", "No session")
	}

	baseHtml := "templates/base.html"
	indexHtml := "templates/index.html"
	todosHtml := "templates/todos.html"
	todoHtml := "templates/todo.html"
	deleteHtml := "templates/delete.html"

	tmpl, tmplErr := template.ParseFiles(baseHtml, indexHtml, todosHtml, todoHtml, deleteHtml)
	if tmplErr != nil {
		utils.Log(utils.ERROR, "index/tmpl", tmplErr.Error())
		http.Error(w, "Intenal server error at tmpl", http.StatusInternalServerError)
		return
	}

	utils.Log(utils.INFO, "index/tmpl", "Template parsed successfully")

	todos, todosErr := utils.GetTodos()
	if todosErr != nil {
		utils.Log(utils.ERROR, "index/todos", todosErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.Log(utils.INFO, "index/todos", "Todos retrieved successfully")

	data := utils.TodosData{
		Session: session,
		Todos:   todos,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	resErr := tmpl.Execute(w, data)
	if resErr != nil {
		utils.Log(utils.ERROR, "index/res", resErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	utils.Log(utils.INFO, "index/res", "Template rendered successfully")
}
