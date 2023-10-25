package main

import (
	"go-htmx/handlers"
	"go-htmx/utils"
	"net/http"
	"os"
)

func main() {
	r := utils.NewRouter()

	r.GET("/favicon.ico", handlers.Favicon)
	r.GET("/", handlers.HomePage)
	r.GET("/signup", handlers.SignupPage)
	r.POST("/signup", handlers.NewUser)
	r.GET("/signin", handlers.SigninPage)
	r.POST("/signin", handlers.Signin)
	r.POST("/signout", handlers.Signout)
	r.POST("/todos", handlers.NewTodo)
	r.DELETE("/todos/:id", handlers.DeleteTodo)
	r.GET("/check", handlers.CheckUser)

	r.SetStaticPath("/static")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		utils.Log(utils.FATAL, "main/listen", err.Error())
		os.Exit(1)
	}
}
