package utils

import (
	"net/http"
	"strings"
)

func GetPath(r *http.Request) []string {
	rawPath := strings.Split(r.URL.Path, "/")
	var path []string
	for _, element := range rawPath {
		if element != "" {
			path = append(path, element)
		}
	}
	return path
}
