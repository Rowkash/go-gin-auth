package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	CookieDomain = "localhost"
	CookiePath   = "/"
	CookieSecure = true
)

func SetCookie(c *gin.Context, accessToken string, refreshToken string) {
	accessTokenMaxAge := int((30 * time.Hour).Seconds())
	refreshTokenMaxAge := int((30 * 24 * time.Hour).Seconds())

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"accessToken",
		accessToken,
		accessTokenMaxAge,
		CookiePath,
		CookieDomain,
		CookieSecure,
		true, // HttpOnly
	)

	c.SetCookie(
		"refreshToken",
		refreshToken,
		refreshTokenMaxAge,
		CookiePath,
		CookieDomain,
		CookieSecure,
		true, // HttpOnly
	)
}

func ClearCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie("accessToken", "", -1, CookiePath, CookieDomain, CookieSecure, true)
	c.SetCookie("refreshToken", "", -1, CookiePath, CookieDomain, CookieSecure, true)
}
