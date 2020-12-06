package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"truth/model"
)

// IsAdmin ...
func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*model.JwtClaims)
		if claims.User.RoleID != 1 {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}
