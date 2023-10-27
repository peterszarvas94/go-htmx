package handlers

import (
	"go-htmx/utils"
	"net/http"
	"strconv"
	"time"
)

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request, pattern string) {
	cookies := r.Cookies()
	var token string
	for _, cookie := range cookies {
		if cookie.Name == "refresh" {
			token = cookie.Value
		}
	}

	// check if refresh token is valid
	claims, jwtErr := utils.ValidateToken(token)
	userId, subErr := claims.GetSubject()
	_, dbErr := utils.GetUserById(userId)

	loggedIn := jwtErr == nil && subErr == nil && dbErr == nil
	if !loggedIn {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userIdInt, userIdErr := strconv.Atoi(userId)
	if userIdErr != nil {
		utils.Log(utils.ERROR, "refersh/userid", userIdErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// get new refresh token as well
	refreshToken, refreshTokenErr := utils.NewToken(userIdInt, utils.REFRESH)
	if refreshTokenErr != nil {
		utils.Log(utils.ERROR, "refresh/token", refreshTokenErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.Log(utils.INFO, "refresh/token", "Refresh token created successfully")

	expires := time.Unix(refreshToken.Expires, 0)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh",
		Value:    refreshToken.Token,
		Path:     "/refresh",
		Expires:  expires,
		HttpOnly: true,
		// Secure:   true,
		// SameSite: http.SameSiteLaxMode,
	})

	// send access token
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}
