package handlers

import (
	"go-htmx/utils"
	"net/http"
)

func FaviconHandler(w http.ResponseWriter, r *http.Request, pattern string) {
	utils.Notfound(w, r)
}
