package main

import (
	"awesomeProject/config"
	"awesomeProject/service"
	"fmt"
)

func main() {
	config.InitDatabse()
	err := config.Migrate()

	if err != nil {
		panic("Error migrating database")
	}

	_, _ = service.CreateUser("teteu", "1234")

	user, uerr := service.AuthUser("teteu", "1234")
	if uerr != nil {
		fmt.Println("Error authenticating user: ", uerr)
	} else {
		fmt.Println(user)
	}
}
