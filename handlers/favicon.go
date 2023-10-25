package handlers

import (
	"go-htmx/utils"
	"net/http"
)

func Favicon(w http.ResponseWriter, r *http.Request, pattern string) {
	utils.Notfound(w, r)
}
