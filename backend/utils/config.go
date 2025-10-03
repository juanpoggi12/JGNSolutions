package utils

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	JWTSecret       string
	JWTExpiresInMin int
}

func LoadConfig() Config {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("Falta JWT_SECRET en el .env")
	}

	expStr := os.Getenv("JWT_EXPIRES_MIN")
	if expStr == "" {
		expStr = "60" // por defecto 60 minutos
	}
	exp, err := strconv.Atoi(expStr)
	if err != nil {
		log.Fatal("JWT_EXPIRES_MIN debe ser un número")
	}

	return Config{
		JWTSecret:       secret,
		JWTExpiresInMin: exp,
	}
}
