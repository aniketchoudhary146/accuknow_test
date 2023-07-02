package user

import (
	"log"
	"net/http"
	"project/pkg/user/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserSignup struct {
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// func (h handler) getUser(ctx *gin.Context) {
// 	ctx.JSON(http.StatusOK, gin.H{"Data": "John Doe"})
// }

func (h handler) userSignup(ctx *gin.Context) {
	var user_data UserSignup
	err := ctx.ShouldBind(&user_data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Message":      "Invalid request params",
			"ErrorMessage": err.Error(),
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user_data.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	var hashed_user models.User
	hashed_user.Password = string(hashedPassword)
	hashed_user.Email = user_data.Email
	hashed_user.Name = user_data.Name

	results := h.DB.Create(&hashed_user)

	if results.Error != nil {
		ctx.JSON(http.StatusInternalServerError, results.Error.Error())
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"Message": "user created",
		"User": gin.H{
			"Name":  hashed_user.Name,
			"Email": hashed_user.Email,
		},
	})
}

func (h handler) userLogin(ctx *gin.Context) {
	var user_login_data UserLogin
	err := ctx.ShouldBind(&user_login_data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Message":      "Invalid request params",
			"ErrorMessage": err.Error(),
		})
		return
	}

	var user models.User
	result := h.DB.Where("email = ?", user_login_data.Email).First(&user)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Login failed, please try again",
			"Error":   result.Error.Error(),
		})
		return
	}
	if result.RowsAffected < 1 {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Email or password is incorrect",
		})
		return

	} else {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(user_login_data.Password))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Message": "Email or password is incorrect",
			})
			return
		}

		session := sessions.Default(ctx)

		var session_uuid string = uuid.NewString()

		session.Set(session_uuid, user.ID)

		session_val := session.Get(session_uuid)
		err = session.Save()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Message":      "Session error",
				"ErrorMessage": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"Message":      "Logged in successfully",
			"session_data": session_val,
			"sid":          session_uuid,
		})

	}
}

func (h handler) getUser(ctx *gin.Context) {
	session := sessions.Default(ctx)
	sid := ctx.Param("sid")
	user_id := session.Get(sid)
	var user_details models.User
	h.DB.Where("id = ?", user_id).First(&user_details)

	user_data := make(map[string]interface{})
	user_data["id"] = user_details.ID
	user_data["name"] = user_details.Name
	user_data["email"] = user_details.Email

	ctx.JSON(http.StatusOK, gin.H{
		"Message":      "User details fetched successfully",
		"sid":          sid,
		"user_details": user_data,
	})
}
