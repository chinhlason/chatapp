package api

import (
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strings"
)

// CheckLogin middleware func to CheckLogin by check the flag in the context
func CheckLogin(rd *redis.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// check in redis
			token := c.Request().Header.Get("Authorization")
			token = strings.TrimPrefix(token, "Bearer ")
			username, id, err := ValidateToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, Response{
					Code:    http.StatusUnauthorized,
					Message: "Token is invalid or expired",
					Data:    err.Error(),
				})
			}
			tokenInRedis := rd.Get(c.Request().Context(), username).Val()
			if tokenInRedis == token {
				return c.JSON(http.StatusUnauthorized, Response{
					Code:    http.StatusUnauthorized,
					Message: "You have been logged out",
					Data:    nil,
				})
			}
			c.Set("id", id)
			return next(c)
		}
	}
}
