package middleware

import (
	"net/http"

	handlers "github.com/aldysp34/mnc_test/rest/handler"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

var jwtSecret = []byte("your-secret-key")

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Missing token"})
		}

		tokenStr := authHeader[len("Bearer "):]

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenStr, &handlers.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
		}

		// Token is valid, proceed to the next handler
		return next(c)
	}
}
