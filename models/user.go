package models

import (
	"awesomeProject/util"
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const DefaultTokenSize int = 32

type User struct {
	gorm.Model

	Username     string `gorm:"unique;not null" json:"username"`
	PasswordHash []byte `json:"password_hash"`
	AccessToken  string `json:"access_token"`
}

type AuthUserRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	AccessToken string `json:"access_token"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func generateAccessToken(size int) (string, error) {
	tokenBytes := make([]byte, size)

	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)
	return token, nil
}

func HashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, util.NewAppError(util.InternalError, "Error generating hash form password")
	}

	return hash, nil
}

func NewUser(username string, password string) (*User, error) {
	hash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	accessToken, err := generateAccessToken(DefaultTokenSize)
	if err != nil {
		return nil, util.NewAppError(util.InternalError, "Error creating random access token for user creation")
	}

	user := &User{
		Username:     username,
		PasswordHash: hash,
		AccessToken:  accessToken,
	}
	return user, nil
}
