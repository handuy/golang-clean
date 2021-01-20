package model

import (
	"errors"
	"log"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id       string
	Email     string
	Password string
}

type NewUser struct {
	Email     string
	Password string
}

func CreateUser(db *gorm.DB, newUser NewUser) (string, error) {
	var user User
	user.Id = uuid.NewV4().String()
	user.Email = newUser.Email
	user.Password = newUser.Password

	err := db.Create(&user).Error
	if err != nil {
		return "", err
	}

	return user.Id, nil
}

func GetUserID(db *gorm.DB, newUser NewUser) (string, error) {
	var user User
	err := db.Where("email = ?", newUser.Email).First(&user).Error
	if err != nil {
		log.Println(err)
		return "", err
	}

	if user.Id == "" {
		return "", errors.New("User chưa tồn tại trong hệ thống")
	}

	comparePass := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newUser.Password))
	if comparePass != nil {
		return "", errors.New("User chưa tồn tại trong hệ thống")
	}

	return user.Id, nil
}

func CheckEmailExist(db *gorm.DB, email string) bool {
	var user User
	err := db.Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return false
	}

	if err != nil {
		return true
	}

	if user.Id == "" {
		return false
	}

	return true
}
