// cmd/server/main.go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Svengalion/Pastebin/internal/handlers"
	"github.com/Svengalion/Pastebin/internal/models"
	"github.com/Svengalion/Pastebin/internal/repos"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("dotenv file not found")
	}

	// Переменные окружения
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	serverPort := os.Getenv("SERVER_PORT")
	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set in the environment")
	}

	// Строка подключения
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to db: %s", err)
	}

	// Автоматическая миграция моделей
	if err := db.AutoMigrate(&models.Paste{}, &models.User{}); err != nil {
		log.Fatalf("Migration error: %s", err)
	}

	// Инициализация репозиториев и хендлеров
	pasteRepo := repos.NewPasteRepos(db)
	pasteHandler := handlers.NewPasteHandler(pasteRepo)
	userRepo := repos.NewUserRepos(db)
	userHandler := handlers.NewUserHandler(userRepo, []byte(jwtSecret))

	// Создание роутера Gin
	router := gin.Default()

	// Маршруты для пользователей
	router.POST("/users/register", userHandler.RegUser)
	router.POST("/users/auth", userHandler.LoginUser)

	// Защищённые маршруты
	authorized := router.Group("/")
	//authorized.Use(middleware.AuthMiddleware([]byte(jwtSecret)))
	authorized.POST("/pastes/new_paste", pasteHandler.CreatePaste)
	authorized.GET("/pastes/", pasteHandler.GetAllPastes)
	authorized.GET("/pastes/:hash", pasteHandler.GetPaste)

	// Запуск сервера
	addr := fmt.Sprintf(":%s", serverPort)
	log.Printf("Server is running at http://localhost%s/", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Can't start server: %s", err)
	}
}
