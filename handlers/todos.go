package handlers

import (
	"go-htmx/utils"
	"html"
	"html/template"
	"net/http"
	"strconv"
)

func TodosHandler(w http.ResponseWriter, r *http.Request) {
	utils.Log(utils.INFO, "todos/path", r.URL.Path)
	path := utils.GetPath(r)

	if len(path) == 1 {
		// get all todos
		if r.Method == "GET" {
			utils.Log(utils.INFO, "todos/get/method", "Method is GET")

			todosHtml := "templates/todos.html"
			todoHtml := "templates/todo.html"
			deleteHtml := "templates/delete.html"

			tmpl, tmplErr := template.ParseFiles(todosHtml, todoHtml, deleteHtml)
			if tmplErr != nil {
				utils.Log(utils.ERROR, "todos/get/tmpl", tmplErr.Error())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}

			utils.Log(utils.INFO, "todos/get/tmpl", "Template parsed successfully")

			todos, getTodosErr := utils.GetTodos()
			if getTodosErr != nil {
				utils.Log(utils.ERROR, "todos/get/todos", getTodosErr.Error())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			utils.Log(utils.INFO, "todos/get/todos", "Todos retrieved successfully")

			todosData := utils.TodosData{
				Session: utils.SessionData{},
				Todos:   todos,
			}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")

			resErr := tmpl.Execute(w, todosData)
			if resErr != nil {
				utils.Log(utils.ERROR, "todos/get/res", resErr.Error())
				http.Error(w, "Internal server  error", http.StatusInternalServerError)
			}

			utils.Log(utils.INFO, "todos/get/res", "Template rendered successfully")
			return
		}

		// add new todo
		if r.Method == "POST" {
			utils.Log(utils.INFO, "todos/add/method", "Method is POST")

			todoHtml := "templates/todo.html"
			deleteHtml := "templates/delete.html"
			newTodoHtml := "templates/new-todo.html"

			tmpl, tmplErr := template.ParseFiles(newTodoHtml, todoHtml, deleteHtml)
			if tmplErr != nil {
				utils.Log(utils.ERROR, "todos/add/tmpl", tmplErr.Error())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			utils.Log(utils.INFO, "todos/add/tmpl", "Template parsed successfully")

			session := utils.CheckSession(r)
			if !session.LoggedIn {
				utils.Log(utils.ERROR, "todos/add/checkSession", "Unauthorized")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}

			utils.Log(utils.INFO, "todos/add/checkSession", "Authorized")

			formErr := r.ParseForm()
			if formErr != nil {
				utils.Log(utils.ERROR, "todos/add/parseForm", formErr.Error())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			utils.Log(utils.INFO, "todos/add/parseForm", "Form parsed successfully")

			text := html.EscapeString(r.FormValue("text"))

			newTodo, newTodoErr := utils.AddTodo(text, r)
			if newTodoErr != nil {
				utils.Log(utils.ERROR, "todos/add/newTodo", newTodoErr.Error())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			utils.Log(utils.INFO, "todos/add/newTodo", "Todo added successfully")

			w.Header().Set("Content-Type", "text/html; charset=utf-8")

			resErr := tmpl.Execute(w, newTodo)
			if resErr != nil {
				utils.Log(utils.ERROR, "todos/add/res", resErr.Error())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}

			utils.Log(utils.INFO, "todos/add/res", "Template rendered successfully")

			return
		}

		utils.Log(utils.ERROR, "todos/path1", "Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if len(path) == 2 {

		// delete todo by id
		if r.Method == "DELETE" {
			utils.Log(utils.INFO, "todos/delete/method", "Method is DELETE")

			session := utils.CheckSession(r)
			if !session.LoggedIn {
				utils.Log(utils.ERROR, "todos/delete/checkSession", "Unauthorized")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			utils.Log(utils.INFO, "todos/delete/checkSession", "Authorized")

			id, pathErr := strconv.Atoi(path[1])
			if pathErr != nil {
				utils.Log(utils.ERROR, "todos/delete/path", pathErr.Error())
				http.Error(w, "Todo id should be a number", http.StatusBadRequest)
				return
			}

			utils.Log(utils.INFO, "todos/delete/path", "Path parsed successfully")

			deleteErr := utils.DeleteTodoById(id)
			if deleteErr != nil {
				utils.Log(utils.ERROR, "todos/delete/deleteTodo", deleteErr.Error())
			}

			utils.Log(utils.INFO, "todos/delete/deleteTodo", "Todo deleted successfully")
			return
		}

		if r.Method == "GET" {
			baseHtml := "templates/base.html"
			notfoundHtml := "templates/404.html"

			tmpl, tmplErr := template.ParseFiles(baseHtml, notfoundHtml)
			if tmplErr != nil {
				utils.Log(utils.ERROR, "todos/path2/notfound/tmpl", tmplErr.Error())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			resErr := tmpl.Execute(w, nil)
			if resErr != nil {
				utils.Log(utils.ERROR, "todos/path2/notfound/res", resErr.Error())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}

			utils.Log(utils.INFO, "todos/path2/notfound/res", "Template rendered successfully")
			return
		}

		utils.Log(utils.ERROR, "todos/path2", "Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	utils.Log(utils.WARNING, "todos/notfound/path", "Path not found")

	baseHtml := "templates/base.html"
	notfoundHtml := "templates/404.html"

	tmpl, tmplErr := template.ParseFiles(baseHtml, notfoundHtml)
	if tmplErr != nil {
		utils.Log(utils.ERROR, "todos/notfound/tmpl", tmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.Log(utils.INFO, "todos/notfound/tmpl", "Template parsed successfully")

	resErr := tmpl.Execute(w, nil)
	if resErr != nil {
		utils.Log(utils.ERROR, "index/notfound/res", resErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	utils.Log(utils.INFO, "index/notfound/res", "Template rendered successfully")
	return
}
