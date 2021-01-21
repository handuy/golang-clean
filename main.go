package main

import (
	"golang-crud/config"
	userRepository "golang-crud/user/repository/mysql"
	userServ "golang-crud/user/service"
	userController "golang-crud/user/controller"

	noteRepository "golang-crud/note/repository/mysql"
	noteServ "golang-crud/note/service"
	noteController "golang-crud/note/controller"

	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	dbConn, tokenSecret, err := config.SetUpDbAndSecret(".")
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	groupRoutes := r.Group("/")

	userRepo := userRepository.NewUserRepository(dbConn)
	userService := userServ.NewUserService(userRepo)

	noteRepo := noteRepository.NewNoteRepo(dbConn)
	noteService := noteServ.NewNoteService(noteRepo)

	userController.NewUserHandler(groupRoutes, userService, tokenSecret)
	noteController.NewNoteHandler(groupRoutes, noteService, userService, tokenSecret)
	
	r.Run("localhost:8181")
}
