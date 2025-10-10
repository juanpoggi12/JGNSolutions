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

	// Importar tus propios paquetes
	"github.com/juanpoggi12/JGNSolutions/backend/handlers"
	"github.com/juanpoggi12/JGNSolutions/backend/middleware"
	"github.com/juanpoggi12/JGNSolutions/backend/repositories"
	"github.com/juanpoggi12/JGNSolutions/backend/services"
	"github.com/juanpoggi12/JGNSolutions/backend/utils"
)

type AppConfig struct {
	Port     string
	MongoURI string
	DBName   string
}

func loadConfig() AppConfig {
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
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return &MongoClient{
		Client: client,
		DB:     client.Database(dbName),
	}, nil
}

func main() {
	cfg := loadConfig()

	// Router con logger y recovery
	r := gin.Default()

	// Conexión a MongoDB
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

	jwtCfg := utils.LoadConfig()

	// --- ⚙️ INICIALIZACIÓN DE CAPAS ---

	// 1️⃣ Repositorios
	userRepo := repositories.NewUserRepository(mc.DB)
	adminRepo := repositories.NewAdminRepository(mc.DB)
	exerciseRepo := repositories.NewExerciseRepository(mc.DB)
	workoutEntryRepo := repositories.NewWorkoutEntryRepository(mc.DB)
	workoutSessionRepo := repositories.NewWorkoutSessionRepository(mc.DB)
	routineRepo := repositories.NewRoutineRepository(mc.DB)
	routineExerciseRepo := repositories.NewRoutineExerciseRepository(mc.DB)
	logRepo := repositories.NewLogRepository(mc.DB)
	userProfileRepo := repositories.NewUserProfileRepository(mc.DB)

	// 2️⃣ Servicios
	userService := services.NewUserService(userRepo)
	adminService := services.NewAdminService(adminRepo, userRepo)
	exerciseService := services.NewExerciseService(exerciseRepo)
	workoutEntryService := services.NewWorkoutEntryService(workoutEntryRepo, workoutSessionRepo)
	workoutSessionService := services.NewWorkoutSessionService(workoutSessionRepo)
	routineService := services.NewRoutineService(routineRepo, routineExerciseRepo, exerciseRepo)
	logService := services.NewLogService(logRepo)
	userProfileService := services.NewUserProfileService(userProfileRepo)

	// 3️⃣ Handlers
	userHandler := handlers.NewUserHandler(userService)
	adminHandler := handlers.NewAdminHandler(adminService)
	exerciseHandler := handlers.NewExerciseHandler(exerciseService)
	workoutEntryHandler := handlers.NewWorkoutEntryHandler(workoutEntryService)
	workoutSessionHandler := handlers.NewWorkoutSessionHandler(workoutSessionService)
	routineHandler := handlers.NewRoutineHandler(routineService)
	adminLogsHandler := handlers.NewAdminLogsHandler(logService)
	userProfileHandler := handlers.NewUserProfileHandler(userProfileService)

	// --- 🚀 REGISTRO DE RUTAS ---

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok":   true,
			"db":   cfg.DBName,
			"time": time.Now().Format(time.RFC3339),
		})
	})

	api := r.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})

		// 👇 Aquí registrás las rutas de usuario
		apiUser := api.Group("/users")
		apiUser.Use(middleware.AuthMiddleware(jwtCfg)) // 🔒 Todas requieren token
		{
			apiUser.POST("/", userHandler.CreateUser)      // solo admin (ya lo valida el service internamente)
			apiUser.GET("/:id", userHandler.GetUserByID)   // admin o el mismo user
			apiUser.PUT("/:id", userHandler.UpdateUser)    // admin o el mismo user
			apiUser.DELETE("/:id", userHandler.DeleteUser) // solo admin
			apiUser.PUT("/change-password", userHandler.ChangePassword)
		}

		// Rutas Admin (solo para usuarios con rol admin)
		apiAdmin := api.Group("/admin")
		apiAdmin.Use(middleware.AuthMiddleware(jwtCfg), middleware.CheckAdmin()) // 🔒 Token + Rol admin
		{
			apiAdmin.GET("/users/count", adminHandler.CountUsers)
			apiAdmin.GET("/exercises/count", adminHandler.CountExercises)
			apiAdmin.GET("/routines/count", adminHandler.CountRoutines)
			apiAdmin.GET("/workouts/count", adminHandler.CountWorkoutSessions)
			apiAdmin.GET("/users", adminHandler.ListUsers)
			apiAdmin.GET("/exercises/top", adminHandler.TopExercises)
			apiAdmin.GET("/routines/top", adminHandler.TopRoutines)
			apiAdmin.GET("/logs", adminLogsHandler.GetLogs)
			apiAdmin.GET("/user-profiles", adminHandler.ListProfiles)
			apiAdmin.GET("/user-profiles/:id", adminHandler.GetProfileByID)
			apiAdmin.DELETE("/user-profiles/:id", adminHandler.DeleteProfile)
		}

		// --- Rutas de ejercicios ---
		apiExercises := api.Group("/exercises")
		apiExercises.Use(middleware.AuthMiddleware(jwtCfg)) // 🔒 Todas requieren token
		{
			apiExercises.POST("/", middleware.CheckAdmin(), exerciseHandler.CreateExercise)      // solo admin
			apiExercises.PUT("/:id", middleware.CheckAdmin(), exerciseHandler.UpdateExercise)    // solo admin
			apiExercises.DELETE("/:id", middleware.CheckAdmin(), exerciseHandler.DeleteExercise) // solo admin

			apiExercises.GET("/:id", exerciseHandler.GetExerciseByID) // visible para todos los logueados
			apiExercises.GET("/search", exerciseHandler.SearchExercises)
		}

		// --- Rutas de entradas de entrenamiento ---
		apiWorkoutEntries := api.Group("/workout-entries")
		apiWorkoutEntries.Use(middleware.AuthMiddleware(jwtCfg)) // todas requieren token
		{
			apiWorkoutEntries.POST("/", workoutEntryHandler.CreateEntry)        // user o admin (valida el service)
			apiWorkoutEntries.GET("/:id", workoutEntryHandler.GetEntryByID)     // user o admin
			apiWorkoutEntries.PUT("/:id", workoutEntryHandler.UpdateEntry)      // user o admin
			apiWorkoutEntries.DELETE("/:id", workoutEntryHandler.DeleteEntry)   // user o admin
			apiWorkoutEntries.GET("/search", workoutEntryHandler.SearchEntries) // búsqueda
		}

		// --- Rutas de sesiones de entrenamiento ---
		apiWorkoutSessions := api.Group("/workout-sessions")
		apiWorkoutSessions.Use(middleware.AuthMiddleware(jwtCfg)) // todas requieren token
		{
			apiWorkoutSessions.POST("/", workoutSessionHandler.CreateSession)       // user o admin
			apiWorkoutSessions.GET("/:id", workoutSessionHandler.GetSessionByID)    // user o admin
			apiWorkoutSessions.PUT("/:id", workoutSessionHandler.UpdateSession)     // user o admin
			apiWorkoutSessions.DELETE("/:id", workoutSessionHandler.DeleteSession)  // user o admin
			apiWorkoutSessions.GET("/search", workoutSessionHandler.SearchSessions) // búsqueda
		}

		// --- Rutas de rutinas ---
		apiRoutines := api.Group("/routines")
		apiRoutines.Use(middleware.AuthMiddleware(jwtCfg)) // todas requieren JWT válido
		{
			apiRoutines.POST("/", routineHandler.CreateRoutine)       // user o admin
			apiRoutines.GET("/:id", routineHandler.GetRoutineByID)    // user o admin
			apiRoutines.PUT("/:id", routineHandler.UpdateRoutine)     // user o admin
			apiRoutines.DELETE("/:id", routineHandler.DeleteRoutine)  // user o admin
			apiRoutines.GET("/search", routineHandler.SearchRoutines) // user o admin

			// Ejercicios dentro de rutinas
			apiRoutines.POST("/:id/exercises", routineHandler.AddExerciseToRoutine)
			apiRoutines.PUT("/exercises/:id", routineHandler.UpdateRoutineExercise)
			apiRoutines.DELETE("/exercises/:id", routineHandler.DeleteRoutineExercise)
			apiRoutines.GET("/exercises/search", routineHandler.ListRoutineExercises)

			// Duplicar rutina
			apiRoutines.POST("/:id/duplicate", routineHandler.DuplicateRoutine)
		}

		// --- Rutas de perfil ---
		apiProfile := api.Group("/profile")
		apiProfile.Use(middleware.AuthMiddleware(jwtCfg)) // solo usuarios logueados
		{
			apiProfile.GET("/search", userProfileHandler.GetProfile)
			apiProfile.PUT("/:id", userProfileHandler.UpdateProfile)
			apiProfile.DELETE("/:id", userProfileHandler.DeleteProfile)
			apiProfile.POST("/", userProfileHandler.CreateProfile)
		}
	}

	// --- 🔊 Servidor HTTP ---
	log.Printf("Servidor escuchando en :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}
