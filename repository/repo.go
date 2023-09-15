package repository

import (
	"mini_project_p2/database"
	"mini_project_p2/models"
)

func GetUserByID(id float64) (models.User, error) {
	db := database.InitDB()
	var user models.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}
