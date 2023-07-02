package main

import (
	"log"
	"project/common/db"
	"project/pkg/notes"
	"project/pkg/user"

	mysession "project/common/session"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	router := gin.Default()

	db.ConnectDatabase()
	database := db.DB

	mysession.CreateSession(router, database)
	user.RegisterRoutes(router, database)
	notes.RegisterRoutes(router, database)

	router.Run()
}
