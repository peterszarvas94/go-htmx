package utils

import (
	"fmt"
	"net/http"
)

func CheckSession(r *http.Request) SessionData {
	//rewrite for auth header, not cookie :)

	// cookies := r.Cookies()
	// var token string
	// for _, cookie := range cookies {
	// 	if cookie.Name == "jwt" {
	// 		token = cookie.Value
	// 	}
	// }

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return SessionData{
			LoggedIn: false,
			User:     UserData{},
		}
	}

	token := authHeader[len("Bearer "):]
	fmt.Printf(token)

	claims, jwtErr := ValidateToken(token)
	userId, subErr := claims.GetSubject()
	user, dbErr := GetUserById(userId)

	loggedIn := jwtErr == nil && subErr == nil && dbErr == nil

	return SessionData{
		LoggedIn: loggedIn,
		User:     user,
	}
}
