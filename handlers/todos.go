package handlers

import (
	"go-htmx/utils"
	"html"
	"html/template"
	"net/http"
	"strconv"
)

func NewTodoHandler(w http.ResponseWriter, r *http.Request, pattern string) {
	todoHtml := "templates/todo.html"
	deleteHtml := "templates/delete.html"
	newTodoHtml := "templates/new-todo.html"

	tmpl, tmplErr := template.ParseFiles(newTodoHtml, todoHtml, deleteHtml)
	if tmplErr != nil {
		utils.Log(utils.ERROR, "todos/add/tmpl", tmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	utils.Log(utils.INFO, "todos/add/tmpl", "Template parsed successfully")

	session := utils.CheckSession(r)
	if !session.LoggedIn {
		utils.Log(utils.ERROR, "todos/add/checkSession", "Unauthorized")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
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
}

func DeleteTodoHandler(w http.ResponseWriter, r *http.Request, pattern string) {
	variables := utils.GetPathVariables(r.URL.Path, pattern)
	idStr, exists := variables["id"]
	if !exists {
		utils.Log(utils.ERROR, "todos/delete/path", "Todo id not found")
		http.Error(w, "Todo id not found", http.StatusBadRequest)
		return
	}

	session := utils.CheckSession(r)
	if !session.LoggedIn {
		utils.Log(utils.ERROR, "todos/delete/checkSession", "Unauthorized")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	utils.Log(utils.INFO, "todos/delete/checkSession", "Authorized")

	id, pathErr := strconv.Atoi(idStr)
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

	w.WriteHeader(http.StatusNoContent)
}
