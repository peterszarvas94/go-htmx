package handlers

import (
	"go-htmx/utils"
	"net/http"
	"strconv"
	"time"
)

/*
RefreshTokenHandler gets a new access token and refresh token
if the refresh token is valid.
*/
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request, pattern string) {
	cookies := r.Cookies()
	var token string
	for _, cookie := range cookies {
		if cookie.Name == "refresh" {
			token = cookie.Value
		}
	}

	if token == "" {
		utils.Log(utils.WARNING, "refresh/token", "No refresh token found")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	utils.Log(utils.INFO, "refresh/token", "Refresh token received")

	// check if refresh token is valid
	claims, jwtErr := utils.ValidateToken(token)
	userId, subErr := claims.GetSubject()
	_, dbErr := utils.GetUserById(userId)

	if jwtErr != nil {
		utils.Log(utils.ERROR, "refresh/jwt", jwtErr.Error())
	}

	if subErr != nil {
		utils.Log(utils.ERROR, "refresh/sub", subErr.Error())
	}

	if dbErr != nil {
		utils.Log(utils.ERROR, "refresh/db", dbErr.Error())
	}

	loggedIn := jwtErr == nil && subErr == nil && dbErr == nil
	if !loggedIn {
		utils.Log(utils.ERROR, "refresh/token", "Refresh token invalid")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userIdInt, userIdErr := strconv.Atoi(userId)
	if userIdErr != nil {
		utils.Log(utils.ERROR, "refersh/userid", userIdErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// get new access token
	accessToken, accessTokenErr := utils.NewToken(userIdInt, utils.Access)
	if accessTokenErr != nil {
		utils.Log(utils.ERROR, "refresh/token", accessTokenErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// get new refresh token as well
	refreshToken, refreshTokenErr := utils.NewToken(userIdInt, utils.Refresh)
	if refreshTokenErr != nil {
		utils.Log(utils.ERROR, "refresh/token", refreshTokenErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.Log(utils.INFO, "refresh/token", "Refresh token created successfully")

	expires := time.Unix(refreshToken.Expires, 0)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// set the new refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh",
		Value:    refreshToken.Token,
		Path:     "/refresh",
		Expires:  expires,
		HttpOnly: true,
		// Secure:   true,
		// SameSite: http.SameSiteLaxMode,
	})

	// send access token in header
	w.Header().Set("Authorization", "Bearer "+accessToken.Token)
	w.WriteHeader(http.StatusOK)
}
