package controller

import (
	"golang-crud/domain"

	"net/http"
	"strconv"

	"github.com/gin-gonic/contrib/jwt"
	"github.com/gin-gonic/gin"
)

type NoteHandler struct {
	NoteService domain.NoteService
	UserService domain.UserService
	TokenSecret string
}

func NewNoteHandler(rg *gin.RouterGroup, noteService domain.NoteService, userService domain.UserService, tokenSecret string) {
	noteHandler := &NoteHandler{
		NoteService: noteService,
		UserService: userService,
		TokenSecret: tokenSecret,
	}
	notes := rg.Group("/notes")

	notes.GET("/", noteHandler.GetAllNote)
	notes.GET("/:id", noteHandler.GetNoteById)
	notes.POST("/new", jwt.Auth(noteHandler.TokenSecret), noteHandler.CreateNote)
	notes.PUT("/update", jwt.Auth(noteHandler.TokenSecret), noteHandler.UpdateNote)
	notes.DELETE("/delete", jwt.Auth(noteHandler.TokenSecret), noteHandler.DeleteNote)
}

func (noteHandler *NoteHandler) GetAllNote(c *gin.Context) {
	result, err := noteHandler.NoteService.GetAll()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (noteHandler *NoteHandler) GetNoteById(c *gin.Context) {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)

	result, err := noteHandler.NoteService.GetById(idInt)
	if err != nil {
		if err.Error() == "Không tìm thấy note" {
			c.JSON(http.StatusNotFound, domain.StatusMessage{
				Message: "Không tìm thấy note",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, domain.StatusMessage{
			Message: "Lỗi server",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (noteHandler *NoteHandler) CreateNote(c *gin.Context) {
	reqToken := c.GetHeader("Authorization")
	userID, errParseToken := noteHandler.UserService.GetUserIDFromToken(reqToken, noteHandler.TokenSecret)
	if errParseToken != nil {
		c.JSON(http.StatusBadRequest, domain.StatusMessage{
			Message: "Invalid token",
		})
		return
	}

	var newNote domain.NewNote
	if err := c.ShouldBindJSON(&newNote); err != nil {
		c.JSON(http.StatusBadRequest, domain.StatusMessage{
			Message: "Bad request",
		})
		return
	}

	result, err := noteHandler.NoteService.Create(newNote, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.StatusMessage{
			Message: "Không thể tạo note mới",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (noteHandler *NoteHandler) UpdateNote(c *gin.Context) {
	reqToken := c.GetHeader("Authorization")
	userID, errParseToken := noteHandler.UserService.GetUserIDFromToken(reqToken, noteHandler.TokenSecret)
	if errParseToken != nil {
		c.JSON(http.StatusBadRequest, domain.StatusMessage{
			Message: "Invalid token",
		})
		return
	}

	var updateNote domain.Note
	if err := c.ShouldBindJSON(&updateNote); err != nil {
		c.JSON(http.StatusBadRequest, domain.StatusMessage{
			Message: "Bad request",
		})
		return
	}

	err := noteHandler.NoteService.Update(updateNote, userID)
	if err != nil {
		if err.Error() == "Bạn không có quyền cập nhật note" {
			c.JSON(http.StatusUnauthorized, domain.StatusMessage{
				Message: "Bạn không có quyền cập nhật note",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, domain.StatusMessage{
			Message: "Không thể cập nhật note",
		})
		return
	}

	c.JSON(http.StatusOK, domain.StatusMessage{
		Message: "Cập nhật thành công",
	})
}

func (noteHandler *NoteHandler) DeleteNote(c *gin.Context) {
	reqToken := c.GetHeader("Authorization")
	userID, errParseToken := noteHandler.UserService.GetUserIDFromToken(reqToken, noteHandler.TokenSecret)
	if errParseToken != nil {
		c.JSON(http.StatusBadRequest, domain.StatusMessage{
			Message: "Invalid token",
		})
		return
	}

	var deleteNote domain.DeletedNote
	if err := c.ShouldBindJSON(&deleteNote); err != nil {
		c.JSON(http.StatusBadRequest, domain.StatusMessage{
			Message: "Bad request",
		})
		return
	}
	if deleteNote.Id == 0 {
		c.JSON(http.StatusBadRequest, domain.StatusMessage{
			Message: "Bad request",
		})
		return
	}

	err := noteHandler.NoteService.Delete(deleteNote, userID)
	if err != nil {
		if err.Error() == "Bạn không có quyền xóa note" {
			c.JSON(http.StatusUnauthorized, domain.StatusMessage{
				Message: "Bạn không có quyền xóa note",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, domain.StatusMessage{
			Message: "Không thể xóa note",
		})
		return
	}

	c.JSON(http.StatusOK, domain.StatusMessage{
		Message: "Xóa thành công",
	})
}
