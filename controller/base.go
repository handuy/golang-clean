package controller

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	uuid "github.com/satori/go.uuid"
)

type Env struct {
	Db          *gorm.DB
	TokenSecret string
}

type DbConfig struct {
	DB_USER    string `mapstructure:"DB_USER"`
	DB_PASS    string `mapstructure:"DB_PASS"`
	DB_ADDRESS string `mapstructure:"DB_ADDRESS"`
	DB_SCHEMA  string `mapstructure:"DB_SCHEMA"`
}

var newLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	logger.Config{
		LogLevel: logger.Info, // Log level
		Colorful: false,       // Disable color
	},
)

func loadConfig(path string) (config DbConfig, err error) {
	var dbInfo DbConfig

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&dbInfo)
	return dbInfo, nil
}

func NewEnv() (*Env, error) {
	config, err := loadConfig(".")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local`,
			config.DB_USER, config.DB_PASS, config.DB_ADDRESS, config.DB_SCHEMA), // auto configure based on currently MySQL version
	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	tokenSecret := uuid.NewV4().String()

	return &Env{
		Db:          db,
		TokenSecret: tokenSecret,
	}, nil
}

func CreateToken(userID, tokenSecret string) (string, error) {
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

func HashAndSalt(userPass []byte) (string, error) {
	// Generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(userPass), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(hash), nil
}

func GetUserIDFromToken(token string, tokenSecret string) (string, error) {
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
