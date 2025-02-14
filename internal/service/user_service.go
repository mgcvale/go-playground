package service

import (
	"awesomeProject/config"
	"awesomeProject/internal/models"
	"awesomeProject/internal/util"
	"errors"
	"fmt"
	"go/types"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
)

type authenticatedFunction[T any] func(user *models.User, data T, db *gorm.DB) error

func requireAuth[T any](decorated authenticatedFunction[T]) func(accessToken string, data T) error {
	return func(accessToken string, data T) error {
		db := config.GetApplicationState().GetDB()
		if db == nil {
			return util.NewAppError(util.BadRequestError, "Error getting GORM db. Try recalling InitializeDatabase().")
		}

		var user models.User
		result := db.Where("access_token=?", accessToken).First(&user)
		if result.Error != nil {
			return util.NewAppError(util.UnauthorizedError, "Error authenticating: Invalid Access Token")
		}

		return decorated(&user, data, db)
	}
}

func AuthUser(data models.AuthUserRequest) (*models.User, error) {
	// get user from username
	// first, we get the database
	db := config.GetApplicationState().GetDB()
	if db == nil {
		return nil, util.NewAppError(util.BadRequestError, "Error getting GORM db. Try recalling InitializeDatabase().")
	}

	// now, we get the user
	var user models.User
	result := db.Where("username = ?", data.Username).First(&user)
	if result.Error != nil {
		return nil, util.NewAppError(util.UnauthorizedError, "Unauthorized: Username not found")
	}

	// finally, we check if the hashes match
	if bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(data.Password)) != nil {
		return nil, util.NewAppError(util.UnauthorizedError, "Unauthorized: passwords do not match")
	}

	return &user, nil
}

func CreateUser(data models.CreateUserRequest) (*models.User, error) {
	// get database
	db := config.GetApplicationState().GetDB()
	if db == nil {
		return nil, util.NewAppError(util.BadRequestError, "Error getting GORM db. Try recalling InitializeDatabase().")
	}

	user, err := models.NewUser(data.Username, data.Password)
	if err != nil {
		return nil, err
	}

	result := db.Create(user)
	if result.Error != nil {
		fmt.Print("ERROR IN SERVICE: ", errors.Unwrap(result.Error))
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed") {
			return nil, util.NewAppError(util.ConflictError, "Username alredy exists")
		}
		return nil, util.NewAppError(util.InternalError, "Internal DB Error")
	}

	return user, nil
}

func deleteUserUndecorated(user *models.User, _ *types.Nil, db *gorm.DB) error {
	deleteResult := db.Delete(&user)

	if deleteResult.Error != nil {
		return util.NewAppError(util.UnauthorizedError, "Error deleting user: Invalid access token")
	}

	if deleteResult.RowsAffected == 0 {
		return util.NewAppError(util.NotFoundError, "No user was deleted: Invalid access token")
	}

	return nil
}

var DeleteUser = requireAuth(deleteUserUndecorated)

func updateUserUndecorated(user *models.User, request models.CreateUserRequest, db *gorm.DB) error {
	if request.Password != "" {
		if len(request.Password) < 8 {
			return util.NewAppError(util.ShortPasswordError, "Password must be at least 8 characters long")
		}
		err := models.UpdatePassword(user, request.Password)
		if err != nil {
			return err
		}
	}

	if request.Username != "" {
		user.Username = request.Username
	}

	result := db.Model(user).Updates(user)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed") {
			return util.NewAppError(util.ConflictError, "Username already exists")
		}
		return util.NewAppError(util.InternalError, "Internal error updating user")
	}
	return nil
}

var UpdateUser = requireAuth(updateUserUndecorated)
