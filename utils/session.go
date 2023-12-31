package utils

import (
	"net/http"
)

/*
CheckSession checks if the user is logged in.
*/
func CheckSession(r *http.Request) SessionData {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return SessionData{
			LoggedIn: false,
			User:     UserData{},
		}
	}

	token := authHeader[len("Bearer "):]

	claims, jwtErr := ValidateToken(token)
	userId, subErr := claims.GetSubject()
	user, dbErr := GetUserById(userId)

	loggedIn := jwtErr == nil && subErr == nil && dbErr == nil

	return SessionData{
		LoggedIn: loggedIn,
		User:     user,
	}
}
