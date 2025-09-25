package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/tuusuario/nombre-equipo/backend/internal/config"
	"github.com/tuusuario/nombre-equipo/backend/internal/db"
)

func main() {
	_ = godotenv.Load(".env") // en prod no hace falta

	cfg := config.Load()

	// Conexión a Mongo (solo para validar)
	m, err := db.Connect(context.Background(), cfg.MongoURI, cfg.DBName)
	if err != nil {
		log.Fatalf("Error conectando a Mongo: %v", err)
	}
	defer func() { _ = m.Client.Disconnect(context.Background()) }()

	r := gin.Default()

	// static + templates (los usaremos en próximos pasos)
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*.html")

	// ruta mínima de salud
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"ok":       true,
			"database": cfg.DBName,
		})
	})

	log.Printf("Escuchando en :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
