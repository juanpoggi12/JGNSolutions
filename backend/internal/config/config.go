package config

import (
	"log"
	"os"
)

type Config struct {
	Port      string
	MongoURI  string
	DBName    string
	JWTSecret string
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Falta variable de entorno requerida: %s", k)
	}
	return v
}

func Load() Config {
	return Config{
		Port:      getEnv("PORT", "8080"),
		MongoURI:  must("MONGO_URI"),
		DBName:    must("DB_NAME"),
		JWTSecret: must("JWT_SECRET"),
	}
}
