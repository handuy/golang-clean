package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DB_USER      string `mapstructure:"DB_USER"`
	DB_PASS      string `mapstructure:"DB_PASS"`
	DB_ADDRESS   string `mapstructure:"DB_ADDRESS"`
	DB_SCHEMA    string `mapstructure:"DB_SCHEMA"`
	TOKEN_SECRET string `mapstructure:"TOKEN_SECRET"`
}

func SetUpDbAndSecret(path string) (*gorm.DB, string, error) {
	var config Config

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, "", err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, "", err
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local`,
			config.DB_USER, config.DB_PASS, config.DB_ADDRESS, config.DB_SCHEMA), // auto configure based on currently MySQL version
	}), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: logger.Info,
				Colorful: false,
			},
		),
	})
	if err != nil {
		return nil, "", err
	}

	return db, config.TOKEN_SECRET, nil
}
