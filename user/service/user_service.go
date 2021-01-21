package service

import (
	"errors"
	"golang-crud/domain"
	"log"
	"strings"
	"time"

	jwt_lib "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo domain.UserRepo
}

func NewUserService(userRepo domain.UserRepo) domain.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (u *userService) CheckEmailExist(email string) bool {
	result := u.userRepo.CheckEmailExist(email)
	return result
}

func (u *userService) HashAndSalt(password []byte) (string, error) {
	// Generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(hash), nil
}

func (u *userService) Create(newUser domain.NewUser) (string, error) {
	var user domain.User
	user.Id = uuid.NewV4().String()
	user.Email = newUser.Email
	user.Password = newUser.Password

	result, err := u.userRepo.Create(user)
	if err != nil {
		return result, err
	}

	return result, err
}

func (u *userService) NewToken(userID string, tokenSecret string) (string, error) {
	// Create the token
	token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
	// Set some claims
	token.Claims = jwt_lib.MapClaims{
		"ID":  userID,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	}
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		log.Println(err)
		return "", err
	}

	return tokenString, nil
}

func (u *userService) GetId(newUser domain.NewUser) (string, error) {
	user, err := u.userRepo.Get(newUser.Email)
	if err != nil {
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

func (u *userService) GetUserIDFromToken(token string, tokenSecret string) (string, error) {
	var userID string
	splitToken := strings.Split(token, "Bearer ")[1]

	claims := jwt_lib.MapClaims{}

	tkn, err := jwt_lib.ParseWithClaims(splitToken, claims, func(token *jwt_lib.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		log.Println(err)
		return userID, err
	}

	if !tkn.Valid {
		return userID, errors.New("Token không hợp lệ")
	}

	for k, v := range claims {
		if k == "ID" {
			userID = v.(string)
		}
	}

	if userID == "" {
		return userID, errors.New("Token không hợp lệ")
	}

	return userID, nil
}
