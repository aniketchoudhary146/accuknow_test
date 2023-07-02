package auth

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func UserAuth(c *gin.Context) {
	session := sessions.Default(c)

	sid := c.Query("sid")
	if sid == "" {
		var body_sid struct {
			Sid string `json:"sid" form:"sid"`
		}
		err := c.ShouldBindBodyWith(&body_sid, binding.JSON)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Message": "Sid missing",
			})
			c.Abort()
			return
		}
		sid = body_sid.Sid
	}

	session_data := session.Get(sid)
	if session_data == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Please login again",
		})
		c.Abort()
	}
	c.Next()
}
