package middleware

import (
	"mini_project_p2/repository"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.ErrUnauthorized
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return echo.ErrUnauthorized
		}

		tokenString := tokenParts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan Anda menggunakan secret key yang sama saat membuat dan memverifikasi token
			return []byte(os.Getenv("secret")), nil
		})

		if err != nil {
			return echo.ErrUnauthorized
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userIDFloat64, ok := claims["id"].(float64)
			if !ok {
				return echo.ErrUnauthorized
			}

			// Konversi userIDFloat64 ke uint
			userID := uint(userIDFloat64)

			// Dapatkan data pengguna dari database berdasarkan userID
			user, err := repository.GetUserByID(float64(userID))
			if err != nil {
				return echo.ErrUnauthorized
			}

			// Simpan data pengguna dalam konteks
			c.Set("user", user)
			return next(c)
		}

		return echo.ErrUnauthorized
	}
}
