package utils

import (
	"fmt"
	"strings"
)

func GetPathVariables(urlPath, pattern string) map[string]string {
	variables := make(map[string]string)

	parts1 := strings.Split(urlPath, "/")
	parts2 := strings.Split(pattern, "/")

	for i, part := range parts2 {
		if strings.HasPrefix(part, ":") && i < len(parts1) {
			variableName := strings.TrimPrefix(part, ":")
			variables[variableName] = parts1[i]
		}
	}

	fmt.Println(variables)

	return variables
}
