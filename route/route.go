package route

import (
	"github.com/gin-gonic/gin"
	"golang-crud/controller"
)

func SetupRoute() *gin.Engine {
	env, err := controller.NewEnv()
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	groupRoutes := r.Group("/")
	userRoutes(groupRoutes, env)
	noteRoutes(groupRoutes, env)

	return r
}
