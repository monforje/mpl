package middleware

import (
	"auth/internal/service"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func JWTMiddleware(authService *service.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "токен не предоставлен"})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "невалидный формат токена"})
			}

			token := parts[1]
			userID, err := authService.ValidateToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "невалидный токен"})
			}

			c.Set("user_id", userID.String())
			return next(c)
		}
	}
}
