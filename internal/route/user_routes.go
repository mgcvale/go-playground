package route

import (
	"awesomeProject/internal/models"
	"awesomeProject/internal/service"
	"awesomeProject/internal/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func RegisterUserRoutes(r *gin.Engine) error {
	group := r.Group("/user")

	group.POST("", createUserRoute)
	group.DELETE("", deleteUserRoute)
	group.PUT("", updateUserRoute)
	group.GET("", getTokenRoute)

	return nil
}

func getAuth(c *gin.Context) (string, error) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UNAUTHORIZED", "info": "Missing bearer token"})
		return "", util.Break
	}
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		fmt.Println("not found bearer token")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UNAUTHORIZED", "info": "Invalid bearer token"})
		return "", util.Break
	}
	return parts[1], nil
}

func createUserRoute(c *gin.Context) {
	var userdata models.CreateUserRequest
	if err := c.ShouldBindJSON(&userdata); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BAD_REQUEST"})
		return
	}
	if len(userdata.Username) == 0 || len(userdata.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BAD_REQUEST", "info": "username is empty or password has less than 8 characters"})
		return
	}

	user, err := service.CreateUser(userdata)
	if err != nil {
		if errors.Is(err, util.ConflictError) {
			c.JSON(http.StatusConflict, gin.H{"message": "CONFLICT"})
			return
		}
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "INTERNAL_SERVER_ERROR"})
		return
	}

	c.JSON(http.StatusCreated, user)
	return
}

func deleteUserRoute(c *gin.Context) {
	token, err := getAuth(c)
	if err != nil {
		return
	}

	error := service.DeleteUser(token, nil)
	if error != nil {
		fmt.Println(error)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UNAUTHORIZED", "info": "Invalid bearer token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func updateUserRoute(c *gin.Context) {
	token, err := getAuth(c)
	if err != nil {
		return
	}

	var data models.CreateUserRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BAD_REQUEST"})
		return
	}

	if err := service.UpdateUser(token, data); err != nil {
		if errors.Is(err, util.ConflictError) {
			c.JSON(http.StatusConflict, gin.H{"message": "CONFLICT"})
			return
		}
		if errors.Is(err, util.UnauthorizedError) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "UNAUTHORIZED", "info": "Invalid Bearer token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "INTERNAL_ERROR"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func getTokenRoute(c *gin.Context) {
	var data models.AuthUserRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BAD_REQUEST", "info": "Missing username and password fields"})
		return
	}

	user, err := service.AuthUser(data)
	if err != nil {
		if errors.Is(err, util.UnauthorizedError) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "UNAUTHORIZED", "info": "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "INTERNAL_SERVER_ERROR"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
	return
}
