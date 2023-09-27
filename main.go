package main

import (
	"fmt"
	"net/http"
)

func main() {
	// handle static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// handle routes
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/todos/", TodosHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
