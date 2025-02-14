package main

import (
	"awesomeProject/config"
	"awesomeProject/internal/route"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDatabse()
	err := config.Migrate()

	if err != nil {
		panic("Error migrating database")
	}

	r := gin.Default()
	route.RegisterUserRoutes(r)
	r.Run(":8000")
}
