package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")
			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, "Token is missing")
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte("secret-key"), nil // Gantilah dengan secret key yang sama dengan yang Anda gunakan di Login
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, "Invalid token")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, "Invalid token claims")
			}

			// Mendapatkan ID dari klaim
			userID := claims["id"].(float64) // ID adalah float64 dalam klaim

			// Menambahkan klaim ID ke konteks
			c.Set("user_id", int(userID))

			return next(c)
		}
	}
}
