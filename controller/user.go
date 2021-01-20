package controller

import (
	"golang-crud/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (env *Env) SignUp(c *gin.Context) {
	var newUser model.NewUser
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, model.StatusMessage{
			Message: "Bad request",
		})
		return
	}

	checkEmailExist := model.CheckEmailExist(env.Db, newUser.Email)
	if checkEmailExist {
		c.JSON(http.StatusBadRequest, model.StatusMessage{
			Message: "Email đã tồn tại trong hệ thống",
		})
		return
	}

	hashedPassword, errHash := HashAndSalt([]byte(newUser.Password))
	if errHash != nil {
		log.Println(errHash)
		c.JSON(http.StatusInternalServerError, model.StatusMessage{
			Message: "Không thể đăng kí tài khoản",
		})
		return
	}

	newUser.Password = hashedPassword
	userID, err := model.CreateUser(env.Db, newUser)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, model.StatusMessage{
			Message: "Không thể đăng kí tài khoản",
		})
		return
	}

	tokenString, err := CreateToken(userID, env.TokenSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.StatusMessage{
			Message: "Không thể tạo token",
		})
		return
	}

	c.JSON(200, gin.H{"token": tokenString})
}

func (env *Env) LogIn(c *gin.Context) {
	var user model.NewUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, model.StatusMessage{
			Message: "Bad request",
		})
		return
	}

	userID, err := model.GetUserID(env.Db, user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.StatusMessage{
			Message: "Sai thông tin đăng nhập",
		})
		return
	}

	tokenString, err := CreateToken(userID, env.TokenSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.StatusMessage{
			Message: "Không thể tạo token",
		})
		return
	}

	c.JSON(200, gin.H{"token": tokenString})
}

func (env *Env) LogOut(c *gin.Context) {

}
