package utils

import "net/http"

func CheckSession(r *http.Request) SessionData {
	cookies := r.Cookies()
	var token string
	for _, cookie := range cookies {
		if cookie.Name == "jwt" {
			token = cookie.Value
		}
	}

	claims, jwtErr := ValidateToken(token)
	userId, subErr := claims.GetSubject()
	user, dbErr := GetUserById(userId)

	loggedIn := jwtErr == nil && subErr == nil && dbErr == nil

	return SessionData{
		LoggedIn: loggedIn,
		User:     user,
	}
}
