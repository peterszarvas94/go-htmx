package handlers

import (
	"go-htmx/utils"
	"net/http"
)

func RefreshToken(w http.ResponseWriter, r *http.Request, pattern string) {
	cookies := r.Cookies()
	var token string
	for _, cookie := range cookies {
		if cookie.Name == "refresh" {
			token = cookie.Value
		}
	}

	claims, jwtErr := utils.ValidateToken(token)
	userId, subErr := claims.GetSubject()
	_, dbErr := utils.GetUserById(userId)

	loggedIn := jwtErr == nil && subErr == nil && dbErr == nil
	if (!loggedIn) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Authorization", "Bearer " + token)
	w.WriteHeader(http.StatusOK)
}
