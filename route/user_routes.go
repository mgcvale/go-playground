package route

import "github.com/gin-gonic/gin"

func RegisterUserRoutes(r *gin.Engine) error {
	group := r.Group("/user")

	group.POST("/create", createUserRoute)

	return nil
}

func createUserRoute(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "success",
	})
}
