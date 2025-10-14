package utils

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	JWTSecret        string
	AccessTTLMinutes int
	RefreshTTLDays   int
	CookieSecure     bool
}

func LoadConfig() Config {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("Falta JWT_SECRET en el .env")
	}

	accessTTL := mustParseInt(getEnv("ACCESS_TTL_MINUTES", "15"), "ACCESS_TTL_MINUTES")
	refreshTTL := mustParseInt(getEnv("REFRESH_TTL_DAYS", "30"), "REFRESH_TTL_DAYS")
	secureCookie := getEnv("REFRESH_COOKIE_SECURE", "false")
	secure, err := strconv.ParseBool(secureCookie)
	if err != nil {
		log.Printf("REFRESH_COOKIE_SECURE inválido, usando false: %v", err)
		secure = false
	}

	return Config{
		JWTSecret:        secret,
		AccessTTLMinutes: accessTTL,
		RefreshTTLDays:   refreshTTL,
		CookieSecure:     secure,
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func mustParseInt(val, field string) int {
	n, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("%s debe ser un número: %v", field, err)
	}
	return n
}
