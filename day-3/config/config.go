package config

import (
	"alterra-agmc-day-3/models"
	"log"
	"os"

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

func getEnvOrDefault(envKey string, defaultVal string) string {
	env := os.Getenv(envKey)
	if env == "" {
		return defaultVal
	}
	return env
}
