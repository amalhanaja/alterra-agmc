package config

import (
	"log"
	"os"
	"strconv"
)

func GetJWTExpirationTime() int64 {
	envVar := GetEnvOrDefault("JWT_EXPIRATION_TIME_IN_MILLIS", "3600000")
	expirationTime, err := strconv.Atoi(envVar)
	if err != nil {
		log.Panic("JWT_EXPIRATION_TIME_IN_MILLIS Env should be a number", err)
	}
	return int64(expirationTime)
}

func GetJWTSecretKey() string {
	return GetEnvOrDefault("JWT_SECRET_KEY", "%D*G-KaPdSgVkYp3s5v8y/B?E(H+MbQeThWmZq4t7w9z$C&F)J@NcRfUjXn2r5u8")
}

func GetEnvOrDefault(envKey string, defaultVal string) string {
	env := os.Getenv(envKey)
	if env == "" {
		return defaultVal
	}
	return env
}
