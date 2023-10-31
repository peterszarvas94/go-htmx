package utils

import (
	"html/template"
	"net/http"
)

/*
Notfound handles the 404 error.
*/
func Notfound(w http.ResponseWriter, r *http.Request) {
	baseHtml := "templates/base.html"
	notfoundHtml := "templates/404.html"

	tmpl, tmplErr := template.ParseFiles(baseHtml, notfoundHtml)
	if tmplErr != nil {
		Log(ERROR, "notfound/tmpl", tmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	Log(INFO, "notfound/tmpl", "Template parsed successfully")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)

	resErr := tmpl.Execute(w, nil)
	if resErr != nil {
		Log(ERROR, "notfound/res", resErr.Error())
		http.Error(w, "Internal server  error", http.StatusInternalServerError)
	}

	Log(INFO, "notfound/res", "Template rendered successfully")
}

/*
MethodNotAllowed handles the 405 error.
*/
func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	baseHtml := "templates/base.html"
	notallowedHtml := "templates/405.html"

	tmpl, tmplErr := template.ParseFiles(baseHtml, notallowedHtml)
	if tmplErr != nil {
		Log(ERROR, "notallowed/tmpl", tmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	Log(INFO, "notallowed/tmpl", "Template parsed successfully")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusMethodNotAllowed)
		
	resErr := tmpl.Execute(w, nil)
	if resErr != nil {
		Log(ERROR, "notallowed/res", resErr.Error())
		http.Error(w, "Internal server  error", http.StatusInternalServerError)
	}

	Log(INFO, "notallowed/res", "Template rendered successfully")
}
