package controller

import (
	"golang-crud/domain"

	"net/http"

	"github.com/gin-gonic/gin"

	"log"
)

type UserHandler struct {
	UserService domain.UserService
	TokenSecret string
}

func NewUserHandler(rg *gin.RouterGroup, userService domain.UserService, tokenSecret string) {
	userHandler := &UserHandler{
		UserService: userService,
		TokenSecret: tokenSecret,
	}
	users := rg.Group("/users")

	users.POST("/signup", userHandler.SignUp)
	users.POST("/login", userHandler.LogIn)
}

func (userHandler *UserHandler) SignUp(c *gin.Context) {
	var newUser domain.NewUser
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, domain.StatusMessage{
			Message: "Bad request",
		})
		return
	}

	checkEmailExist := userHandler.UserService.CheckEmailExist(newUser.Email)
	if checkEmailExist {
		c.JSON(http.StatusBadRequest, domain.StatusMessage{
			Message: "Email đã tồn tại trong hệ thống",
		})
		return
	}

	hashedPassword, errHash := userHandler.UserService.HashAndSalt([]byte(newUser.Password))
	if errHash != nil {
		log.Println(errHash)
		c.JSON(http.StatusInternalServerError, domain.StatusMessage{
			Message: "Không thể đăng kí tài khoản",
		})
		return
	}

	newUser.Password = hashedPassword
	userID, err := userHandler.UserService.Create(newUser)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, domain.StatusMessage{
			Message: "Không thể đăng kí tài khoản",
		})
		return
	}

	tokenString, err := userHandler.UserService.NewToken(userID, userHandler.TokenSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.StatusMessage{
			Message: "Không thể tạo token",
		})
		return
	}

	c.JSON(200, gin.H{"token": tokenString})
}

func (userHandler *UserHandler) LogIn(c *gin.Context) {
	var user domain.NewUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, domain.StatusMessage{
			Message: "Bad request",
		})
		return
	}

	userID, err := userHandler.UserService.GetId(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.StatusMessage{
			Message: "Sai thông tin đăng nhập",
		})
		return
	}

	tokenString, err := userHandler.UserService.NewToken(userID, userHandler.TokenSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.StatusMessage{
			Message: "Không thể tạo token",
		})
		return
	}

	c.JSON(200, gin.H{"token": tokenString})
}
