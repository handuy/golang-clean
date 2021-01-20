package route

import (
	"github.com/gin-gonic/gin"
	"golang-crud/controller"
	"github.com/gin-gonic/contrib/jwt"
)

func noteRoutes(rg *gin.RouterGroup, env *controller.Env) {
	notes := rg.Group("/notes")

	notes.GET("/", env.GetAllNote)
	notes.GET("/:id", env.GetNoteById)
	notes.POST("/new", jwt.Auth(env.TokenSecret), env.CreateNote)
	notes.POST("/update", jwt.Auth(env.TokenSecret), env.UpdateNote)
	notes.POST("/delete", jwt.Auth(env.TokenSecret), env.DeleteNote)
}