package main

import (
	"log"
	"webapp/routes"
	"webapp/utils"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db := utils.SetupDatabase()
	r := routes.SetupRouter(db)
	r.Run()
}
