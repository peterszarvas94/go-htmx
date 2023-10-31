package utils

import (
	// "fmt"
	"net/http"
	// "reflect"
	// "runtime"
	"strings"
)

/*
Router struct.
Maps: pattern -> method -> handlerFunc.
Eg.: /api/v1/user/:id -> GET -> GetUserHandler
*/
type Router struct {
	routes       map[string]map[string]HandlerFunc
	staticPrefix string
}

type HandlerFunc func(http.ResponseWriter, *http.Request, string)

/*
Utility function for creating a new router.
*/
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string]HandlerFunc),
	}
}

/*
SetStaticPath sets the static path for serving static files.
Static files should be put in the static folder.
*/
func (r *Router) SetStaticPath(prefix string) {
	r.staticPrefix = prefix
}

/*
Utility function for adding a new route to the router.
*/
func (r *Router) addRoute(method, pattern string, handler HandlerFunc) {
	// adds new route if it doesn't exist
	if _, ok := r.routes[pattern]; !ok {
		r.routes[pattern] = make(map[string]HandlerFunc)
	}
	// overwrites the route if it already exists
	r.routes[pattern][method] = handler
}

/*
Adds a new GET route to the router.
*/
func (r *Router) GET(pattern string, handler HandlerFunc) {
	r.addRoute("GET", pattern, handler)
}

/*
Adds a new POST route to the router.
*/
func (r *Router) POST(pattern string, handler HandlerFunc) {
	r.addRoute("POST", pattern, handler)
}

/*
Adds a new PUT route to the router.
*/
func (r *Router) PUT(pattern string, handler HandlerFunc) {
	r.addRoute("PUT", pattern, handler)
}

/*
Adds a new DELETE route to the router.
*/
func (r *Router) DELETE(pattern string, handler HandlerFunc) {
	r.addRoute("DELETE", pattern, handler)
}

/*
ServeHTTP is mandatory
*/
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method

	// handle static files
	if router.staticPrefix != "" && strings.HasPrefix(path, router.staticPrefix) {
		fs := http.FileServer(http.Dir("static"))
		staticHandler := http.StripPrefix(router.staticPrefix, fs)
		staticHandler.ServeHTTP(w, r)
		return
	}

	// check if the route and method are registered
	for pattern, handlers := range router.routes {
		if matches(pattern, path) {
			if handler, exists := handlers[method]; exists {
				handler(w, r, pattern)
				return
			}

			if method == "GET" {
				MethodNotAllowed(w, r)
				return
			}

			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}

	Notfound(w, r)
}

/*
Url pattern matching.
Eg.: /api/v1/user/:id -> /api/v1/user/1/
*/
func matches(pattern, path string) bool {

	// remove trailing slash
	pattern = removeTrailingslash(pattern)
	path = removeTrailingslash(path)

	patterns := strings.Split(pattern, "/")
	paths := strings.Split(path, "/")

	// check if the number of patterns and paths are the same
	if len(patterns) != len(paths) {
		return false
	}

	// check if the patterns and paths match
	for i, pattern := range patterns {
		if pattern == paths[i] || strings.HasPrefix(pattern, ":") {
			continue
		}

		return false
	}

	return true
}

// remove trailing slash
func removeTrailingslash(path string) string {
	if path != "/" && strings.HasSuffix(path, "/") {
		return path[:len(path)-1]
	}

	return path
}

// list all routes with path, method and handler function name
// func (router *Router) ListRoutes() {
// 	for path, handlers := range router.routes {
// 		for method, handlerFunc := range handlers {
// 			handlerFuncName := runtime.FuncForPC(reflect.ValueOf(handlerFunc).Pointer()).Name()
// 			fmt.Printf("%s %s -> %s\n", method, path, handlerFuncName)
// 		}
// 	}
// }
