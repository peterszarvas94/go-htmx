package utils

import (
	"net/http"
	"strings"
)

/*
GetPath returns a slice of strings from a url path
Example:
urlPath: /users/123
returns: []string{"users", "123"}
*/
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

/*
GetPathVariables returns a map of path variables from a url path and a pattern
Example:
urlPath: /users/123
pattern: /users/:id
returns: map[string]string{"id": "123"}
*/
func GetPathVariables(urlPath, pattern string) map[string]string {
	variables := make(map[string]string)

	urlParts := strings.Split(urlPath, "/")
	patternParts := strings.Split(pattern, "/")

	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") && i < len(urlParts) {
			variableName := strings.TrimPrefix(part, ":")
			variables[variableName] = urlParts[i]
		}
	}

	return variables
}
