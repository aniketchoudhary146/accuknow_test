package session

import (
	"os"

	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserSession struct {
	UserId int
	Name   string
	Email  string
}

func CreateSession(router *gin.Engine, db *gorm.DB) {

	key := os.Getenv("SESSION_SECRET")

	store := gormsessions.NewStore(db, true, []byte(key))

	router.Use(sessions.Sessions("mysession", store))

}
