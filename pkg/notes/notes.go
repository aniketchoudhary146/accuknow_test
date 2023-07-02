package notes

import (
	"net/http"
	"project/pkg/notes/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Note struct {
	Sid  string `json:"sid" binding:"required"`
	Note string `json:"note" binding:"required"`
}

type GetNotesResponse struct {
	Note string
	Id   uint
}

func (h handler) createNote(ctx *gin.Context) {
	var note Note
	err := ctx.ShouldBindBodyWith(&note, binding.JSON)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Message":      "Something went wrong, please try again",
			"ErrorMessage": err.Error(),
		})
		return
	}

	session := sessions.Default(ctx)
	user_id := session.Get(note.Sid)

	note_data := models.UserNotes{UserId: user_id.(uint), Message: note.Note}
	result := h.DB.Create(&note_data)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Message":      "Something went wrong, please try again",
			"ErrorMessage": result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Message": "Note added successfully",
		"Note":    note.Note,
	})

}

func (h handler) getNotes(ctx *gin.Context) {
	var sid struct {
		Sid string `json:"sid" binding:"required"`
	}
	err := ctx.ShouldBindBodyWith(&sid, binding.JSON)
	if sid.Sid == "" {
		sid.Sid = ctx.Query("sid")
	}

	if err != nil && sid.Sid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Message":      "Invalid request params",
			"ErrorMessage": err.Error(),
		})
		return
	} else if sid.Sid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Message":      "Invalid request params",
			"ErrorMessage": "sid missing",
		})
		return
	}
	session := sessions.Default(ctx)
	user_id := session.Get(sid.Sid)

	var notes []models.UserNotes

	result := h.DB.Where("user_id = ?", user_id).Find(&notes)

	var response []GetNotesResponse

	for _, note := range notes {
		var item GetNotesResponse
		item.Id = note.ID
		item.Note = note.Message

		response = append(response, item)
	}

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Message":      "Something went wrong, please try again",
			"ErrorMessage": result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (h handler) deleteNote(ctx *gin.Context) {
	var deleteNoteRequest struct {
		Sid string `json:"sid" binding:"required"`
		Id  uint   `json:"id" binding:"required"`
	}

	err := ctx.ShouldBindBodyWith(&deleteNoteRequest, binding.JSON)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Message":      "Invalid request params",
			"ErrorMessage": err.Error(),
		})
		return
	}
	var note models.UserNotes
	note.ID = deleteNoteRequest.Id
	result := h.DB.Delete(&note)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Message":      "Something went wrong, please try again",
			"ErrorMessage": result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Message": "Note removed successfully",
	})
}
