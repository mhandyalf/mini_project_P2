package handlers

import (
	"mini_project_p2/models"
	"time"

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

	return c.JSON(201, user)
}

func (h *Auth) Login(c echo.Context) error {
	var user models.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(400, err)
	}

	var foundUser models.User

	if err := h.DB.Where("email = ?", user.Email).First(&user).Error; err != nil {
		return c.JSON(404, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(foundUser.Password)); err != nil {
		return c.JSON(401, err)
	}

	// Buat token JWT
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID                                // Menambahkan klaim ID ke dalam token
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Menambahkan waktu kadaluarsa

	// Tandatangani token dengan secret key (gantilah dengan secret key yang kuat)
	tokenString, err := token.SignedString([]byte("secret-key"))
	if err != nil {
		return c.JSON(500, err)
	}

	// Mengembalikan token JWT
	return c.JSON(200, map[string]string{
		"token": tokenString,
	})

}
