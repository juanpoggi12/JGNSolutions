package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// --- Config ---

type AppConfig struct {
	Port     string
	MongoURI string
	DBName   string
}

func loadConfig() AppConfig {
	// En desarrollo: carga .env si existe. En producciÃ³n, las vars vienen del entorno.
	_ = godotenv.Load()

	return AppConfig{
		Port:     getEnv("PORT", "8080"),
		MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:   getEnv("DB_NAME", "jgnsolutions"),
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// --- Mongo ---

type MongoClient struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func connectMongo(ctx context.Context, uri, dbName string) (*MongoClient, error) {
	clientOpts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}
	// Verificamos conexiÃ³n con ping
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &MongoClient{
		Client: client,
		DB:     client.Database(dbName),
	}, nil
}

// --- main ---

func main() {
	cfg := loadConfig()

	// Router con logger y recovery
	r := gin.Default()

	// ConexiÃ³n a Mongo con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mc, err := connectMongo(ctx, cfg.MongoURI, cfg.DBName)
	if err != nil {
		log.Fatalf("Error conectando a MongoDB: %v", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := mc.Client.Disconnect(shutdownCtx); err != nil {
			log.Printf("Error cerrando MongoDB: %v", err)
		}
	}()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok":   true,
			"db":   cfg.DBName,
			"time": time.Now().Format(time.RFC3339),
		})
	})

	// Grupo base para tu API real
	api := r.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
		// Ejemplos a futuro:
		// api.POST("/auth/register", registerHandler(mc.DB))
		// api.POST("/auth/login", loginHandler(mc.DB))
		// api.GET("/exercises", authMiddleware(), listExercisesHandler(mc.DB))
	}

	log.Printf("Servidor escuchando en :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}
