package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed index", http.StatusMethodNotAllowed)
		return
	}

	index := "templates/index.html"
	todo := "templates/todo.html"
	todos := "templates/todos.html"
	tmpl, err := template.ParseFiles(index, todos, todo)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	data := Data{
		Todos: GetTodos(),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func TodosHandler(w http.ResponseWriter, r *http.Request) {
	rawPath := strings.Split(r.URL.Path, "/")
	var path []string
	for _, element := range rawPath {
		if element != "" {
			path = append(path, element)
		}
	}
	todo := "templates/todo.html"
	todos := "templates/todos.html"

	// get all todos
	if len(path) == 1 && r.Method == "GET" {
		tmpl, err := template.ParseFiles(todos, todo)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		data := Data{
			Todos: GetTodos(),
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// add new todo
	if len(path) == 1 && r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		text := r.FormValue("text")
		new_todo := AddTodo(text)

		tmpl, err := template.ParseFiles(todo)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = tmpl.Execute(w, new_todo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// get todo by id
	if len(path) == 2 && r.Method == "GET" {
		id, err := strconv.Atoi(path[1])
		if err != nil {
			http.Error(w, "Invalid Id", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles(todo)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		todo_data := GetTodoById(id)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = tmpl.Execute(w, todo_data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// delete todo by id
	if len(path) == 2 && r.Method == "DELETE" {
		id, err := strconv.Atoi(path[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		DeleteTodoById(id)

		return
	}
}
