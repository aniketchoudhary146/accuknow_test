package notes

import (
	//"project/pkg/auth"

	"project/common/auth"

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

	routes := router.Group("/v1/notes", auth.UserAuth)
	routes.POST("/", h.createNote)
	routes.GET("/", h.getNotes)
	routes.DELETE("/", h.deleteNote)

}
