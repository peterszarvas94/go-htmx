package handlers

import (
	"fmt"
	"go-htmx/utils"
	"html/template"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request, pattern string) {
	utils.Log(utils.INFO, "index/path", r.URL.Path)

	fmt.Printf(r.Header.Get("Authorization"))

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
