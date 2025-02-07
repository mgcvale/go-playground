package service

import (
	"awesomeProject/config"
	"awesomeProject/models"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func AuthUser(username string, password string) (*models.User, error) {
	// get user from username
	// first, we get the database
	db := config.GetApplicationState().GetDB()
	if db == nil {
		return nil, fmt.Errorf("Error getting GORM db. Try recalling InitializeDatabase().")
	}

	// now, we get the user
	var user models.User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("UNAUTHORIZED")
	}

	// finally, we check if the hashes match
	if bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)) != nil {
		return nil, fmt.Errorf("UNAUTHORIZED")
	}

	return &user, nil
}

func CreateUser(username string, password string) (*models.User, error) {
	// get database
	db := config.GetApplicationState().GetDB()
	if db == nil {
		return nil, fmt.Errorf("Error getting database from application state. Call InitializeDatabase() first")
	}

	user, err := models.NewUser(username, password)
	if err != nil {
		return nil, err
	}

	db.Create(user)
	return user, nil
}
