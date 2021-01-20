package controller

import (
	"golang-crud/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (env *Env) GetAllNote(c *gin.Context) {
	result, err := model.GetAllNotes(env.Db)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (env *Env) GetNoteById(c *gin.Context) {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)

	result, err := model.GetNoteById(env.Db, idInt)
	if err != nil {
		if err.Error() == "Không tìm thấy note" {
			c.JSON(http.StatusNotFound, model.StatusMessage{
				Message: "Không tìm thấy note",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.StatusMessage{
			Message: "Lỗi server",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (env *Env) CreateNote(c *gin.Context) {
	reqToken := c.GetHeader("Authorization")
	userID, errParseToken := GetUserIDFromToken(reqToken, env.TokenSecret)
	if errParseToken != nil {
		c.JSON(http.StatusBadRequest, model.StatusMessage{
			Message: "Invalid token",
		})
		return
	}

	var newNote model.NewNote
	if err := c.ShouldBindJSON(&newNote); err != nil {
		c.JSON(http.StatusBadRequest, model.StatusMessage{
			Message: "Bad request",
		})
		return
	}

	result, err := model.InsertNote(env.Db, newNote, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.StatusMessage{
			Message: "Không thể tạo note mới",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (env *Env) UpdateNote(c *gin.Context) {
	reqToken := c.GetHeader("Authorization")
	userID, errParseToken := GetUserIDFromToken(reqToken, env.TokenSecret)
	if errParseToken != nil {
		c.JSON(http.StatusBadRequest, model.StatusMessage{
			Message: "Invalid token",
		})
		return
	}

	var updateNote model.UpdatedNote
	if err := c.ShouldBindJSON(&updateNote); err != nil {
		c.JSON(http.StatusBadRequest, model.StatusMessage{
			Message: "Bad request",
		})
		return
	}

	updateNote.UpdatedAt = time.Now()
	err := model.UpdateNote(env.Db, updateNote, userID)
	if err != nil {
		if err.Error() == "Bạn không có quyền cập nhật note" {
			c.JSON(http.StatusUnauthorized, model.StatusMessage{
				Message: "Bạn không có quyền cập nhật note",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.StatusMessage{
			Message: "Không thể cập nhật note",
		})
		return
	}

	c.JSON(http.StatusOK, model.StatusMessage{
		Message: "Cập nhật thành công",
	})
}

func (env *Env) DeleteNote(c *gin.Context) {
	reqToken := c.GetHeader("Authorization")
	userID, errParseToken := GetUserIDFromToken(reqToken, env.TokenSecret)
	if errParseToken != nil {
		c.JSON(http.StatusBadRequest, model.StatusMessage{
			Message: "Invalid token",
		})
		return
	}

	var deleteNote model.DeletedNote
	if err := c.ShouldBindJSON(&deleteNote); err != nil {
		c.JSON(http.StatusBadRequest, model.StatusMessage{
			Message: "Bad request",
		})
		return
	}
	if deleteNote.Id == 0 {
		c.JSON(http.StatusBadRequest, model.StatusMessage{
			Message: "Bad request",
		})
		return
	}

	err := model.DeleteNote(env.Db, deleteNote, userID)
	if err != nil {
		if err.Error() == "Bạn không có quyền xóa note" {
			c.JSON(http.StatusUnauthorized, model.StatusMessage{
				Message: "Bạn không có quyền xóa note",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.StatusMessage{
			Message: "Không thể xóa note",
		})
		return
	}

	c.JSON(http.StatusOK, model.StatusMessage{
		Message: "Xóa thành công",
	})
}
