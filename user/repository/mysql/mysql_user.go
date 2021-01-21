package mysql

import (
	"golang-crud/domain"

	"log"

	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepo {
	return &userRepo{db: db}
}

func (u *userRepo) CheckEmailExist(email string) bool {
	var user domain.User
	err := u.db.Where("email = ?", email).First(&user).Error
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

func (u *userRepo) Create(newUser domain.User) (string, error) {
	err := u.db.Create(&newUser).Error
	if err != nil {
		return "", err
	}

	return newUser.Id, nil
}

func (u *userRepo) Get(email string) (domain.User, error) {
	var user domain.User
	err := u.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		log.Println(err)
		return user, err
	}

	return user, nil
}
