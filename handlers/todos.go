package handlers

import (
	"go-htmx/utils"
	"html"
	"html/template"
	"net/http"
	"strconv"
)

func TodosHandler(w http.ResponseWriter, r *http.Request) {
	path := utils.GetPath(r)

	todoHtml := "templates/todo.html"
	todosHtml := "templates/todos.html"

	todosTmpl, todosTmplErr := template.ParseFiles(todosHtml, todoHtml)
	if todosTmplErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	todoTmpl, todoTmplErr := template.ParseFiles(todoHtml)
	if todoTmplErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	// get all todos
	if len(path) == 1 && r.Method == "GET" {
		todos, getTodosErr := utils.GetTodos()
		if getTodosErr != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		todosData := utils.TodosData{
			Todos: todos,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		resErr := todosTmpl.Execute(w, todosData)
		if resErr != nil {
			http.Error(w, "Internal server  error", http.StatusInternalServerError)
		}
		return
	}

	// add new todo
	if len(path) == 1 && r.Method == "POST" {
		session := utils.CheckSession(r)
		if !session.LoggedIn {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		formErr := r.ParseForm()
		if formErr != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		text := html.EscapeString(r.FormValue("text"))

		todoData, newTodoErr := utils.AddTodo(text, r)
		if newTodoErr != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		resErr := todoTmpl.Execute(w, todoData)
		if resErr != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		return
	}

	// get todo by id
	if len(path) == 2 && r.Method == "GET" {
		id, pathErr := strconv.Atoi(path[1])
		if pathErr != nil {
			http.Error(w, "Todo id should be a number", http.StatusBadRequest)
			return
		}

		todoData, getTodoErr := utils.GetTodoById(id)
		if getTodoErr != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		resErr := todoTmpl.Execute(w, todoData)
		if resErr != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// delete todo by id
	if len(path) == 2 && r.Method == "DELETE" {
		session := utils.CheckSession(r)
		if !session.LoggedIn {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		id, pathErr := strconv.Atoi(path[1])
		if pathErr != nil {
			http.Error(w, "Todo id should be a number", http.StatusBadRequest)
			return
		}

		utils.DeleteTodoById(id)

		return
	}
}
