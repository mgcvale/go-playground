package main

import (
	"awesomeProject/config"
)

func main() {
	config.InitDatabse()
	err := config.Migrate()

	if err != nil {
		panic("Error migrating database")
	}

}
