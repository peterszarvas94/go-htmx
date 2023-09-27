package main

import (
	"fmt"
	"net/http"
	"go-htmx/handlers"
)

func main() {
	// handle static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// handle routes
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/todos/", handlers.TodosHandler)
	http.HandleFunc("/auth/", handlers.AuthHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
