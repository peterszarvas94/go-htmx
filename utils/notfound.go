package utils

import (
	"html/template"
	"net/http"
)

func NotfoundHandler(w http.ResponseWriter, r *http.Request) {
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

	resErr := tmpl.Execute(w, nil)
	if resErr != nil {
		Log(ERROR, "notfound/res", resErr.Error())
		http.Error(w, "Internal server  error", http.StatusInternalServerError)
	}

	Log(INFO, "notfound/res", "Template rendered successfully")
}
