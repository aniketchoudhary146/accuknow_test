package user

import (
	//"project/pkg/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {

	h := &handler{
		DB: db,
	}

	routes := router.Group("/v1/user")
	routes.POST("/signup", h.userSignup)
	routes.POST("/login", h.userLogin)
	routes.GET("/:sid", h.getUser) // session test endpoint

}
