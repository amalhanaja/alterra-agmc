package database

import (
	"alterra-agmc-day-2/config"
	"alterra-agmc-day-2/models"
)

func GetUsers() ([]models.User, error) {
	var users []models.User
	err := config.DB.Find(&users).Error
	return users, err
}

func GetUserById(id uint) (*models.User, error) {
	user := &models.User{}
	if err := config.DB.First(user, id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func CreateUser(user *models.User) (*models.User, error) {
	if err := config.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUser(user *models.User) (*models.User, error) {
	if err := config.DB.Model(&user).Updates(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func DeleteUser(id uint) error {
	return config.DB.Delete(&models.User{}, id).Error
}
