package service

import (
	"awesomeProject/config"
	"awesomeProject/models"
	"awesomeProject/util"
	"golang.org/x/crypto/bcrypt"
)

func AuthUser(username string, password string) (*models.User, error) {
	// get user from username
	// first, we get the database
	db := config.GetApplicationState().GetDB()
	if db == nil {
		return nil, util.NewAppError(util.BadRequestError, "Error getting GORM db. Try recalling InitializeDatabase().")
	}

	// now, we get the user
	var user models.User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, util.NewAppError(util.UnauthorizedError, "Unauthorized: Username not found")
	}

	// finally, we check if the hashes match
	if bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)) != nil {
		return nil, util.NewAppError(util.UnauthorizedError, "Unauthorized: passwords do not match")
	}

	return &user, nil
}

func CreateUser(username string, password string) (*models.User, error) {
	// get database
	db := config.GetApplicationState().GetDB()
	if db == nil {
		return nil, util.NewAppError(util.BadRequestError, "Error getting GORM db. Try recalling InitializeDatabase().")
	}

	user, err := models.NewUser(username, password)
	if err != nil {
		return nil, err
	}

	db.Create(user)
	return user, nil
}

func DeleteUser(access_token string) error {
	db := config.GetApplicationState().GetDB()
	if db == nil {
		return util.NewAppError(util.BadRequestError, "Error getting GORM db. Try recalling InitializeDatabase().")
	}

	var user models.User
	result := db.Where("access_token = ?", access_token).First(&user)
	if result.Error != nil {
		return util.NewAppError(util.UnauthorizedError, "Error deleting user: Invalid Access Token")
	}

	deleteResult := db.Delete(&user)

	if deleteResult.Error != nil {
		return util.NewAppError(util.UnauthorizedError, "Error deleting user: Invalid access token")
	}

	if deleteResult.RowsAffected == 0 {
		return util.NewAppError(util.NotFoundError, "No user was deleted: Invalid access token")
	}

	return nil
}
