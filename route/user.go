package route

import (
	"github.com/gin-gonic/gin"
	"golang-crud/controller"
)

func userRoutes(rg *gin.RouterGroup, env *controller.Env) {
	notes := rg.Group("/users")

	notes.POST("/signup", env.SignUp)
	notes.POST("/login", env.LogIn)
	notes.POST("/logout", env.LogOut)
}