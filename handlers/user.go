package handlers

import (
	"mini_project_p2/middleware"
	"mini_project_p2/models"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Auth struct {
	DB *gorm.DB
}

func NewAuth(db *gorm.DB) *Auth {
	return &Auth{DB: db}
}

func (h *Auth) Register(c echo.Context) error {
	var user models.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(400, err)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(400, err)
	}

	user.Password = string(hashPassword)

	if err := h.DB.Create(&user).Error; err != nil {
		return err
	}

	// Mengirim email konfirmasi
	if err := middleware.SendConfirmationEmail(&user); err != nil {
		return err
	}

	return c.JSON(201, user)
}

func (a *Auth) Login(e echo.Context) error {
	var user models.User

	if err := e.Bind(&user); err != nil {
		return e.JSON(http.StatusBadRequest, err)
	}

	var foundUser models.User
	if err := a.DB.Where("email = ?", user.Email).First(&foundUser).Error; err != nil {
		return e.JSON(http.StatusNotFound, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
		return e.JSON(http.StatusUnauthorized, err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": float64(foundUser.ID),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("secret")))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"message": "gagal membuat token JWT",
		})
	}

	return e.JSON(http.StatusOK, map[string]string{
		"token": tokenString,
	})

}
