package config

import (
	"alterra-agmc-day-5-6/internal/models"
	"log"
	"os"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := getEnvOrDefault("DB_DSN", "root:password@tcp(localhost:3306)/development?charset=utf8mb4&parseTime=True&loc=Local")
	var err error
	DB, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Panic(err)
	}
	log.Println("DB Connected")
	initMigration()
}

func initMigration() {
	DB.AutoMigrate(&models.User{})
}

func GetJWTExpirationTime() int64 {
	envVar := getEnvOrDefault("JWT_EXPIRATION_TIME_IN_MILLIS", "3600000")
	expirationTime, err := strconv.Atoi(envVar)
	if err != nil {
		log.Panic("JWT_EXPIRATION_TIME_IN_MILLIS Env should be a number", err)
	}
	return int64(expirationTime)
}

func GetJWTSecretKey() string {
	return getEnvOrDefault("JWT_SECRET_KEY", "%D*G-KaPdSgVkYp3s5v8y/B?E(H+MbQeThWmZq4t7w9z$C&F)J@NcRfUjXn2r5u8")
}

func getEnvOrDefault(envKey string, defaultVal string) string {
	env := os.Getenv(envKey)
	if env == "" {
		return defaultVal
	}
	return env
}
