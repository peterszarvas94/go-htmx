package main

import (
	"go-htmx/handlers"
	"go-htmx/utils"
	"net/http"
	"os"
)

func main() {
	// handle static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	//handle favicon
	http.Handle("/favicon.ico", http.NotFoundHandler())

	// handle routes
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/todos/", handlers.TodosHandler)
	http.HandleFunc("/signup/", handlers.SignupHandler)
	http.HandleFunc("/signin/", handlers.SigninHandler)
	http.HandleFunc("/signout/", handlers.SignoutHandler)
	http.HandleFunc("/check/", handlers.CheckHandler)
	http.HandleFunc("/404/", handlers.NotfoundHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		utils.Log(utils.FATAL, "main/listen", err.Error())
		os.Exit(1)
	}
}
