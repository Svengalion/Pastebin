package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/Svengalion/Pastebin/cmd/server/docs" // Импорт для Swagger
	"github.com/Svengalion/Pastebin/internal/handlers"
	"github.com/Svengalion/Pastebin/internal/models"
	"github.com/Svengalion/Pastebin/internal/repos"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title           Pastebin Clone API
// @version         1.0
// @description     API документация для Pastebin-клона.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("dotenv file not found")
	}

	//переменные окружения
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	serverPort := os.Getenv("SERVER_PORT")

	//строка подключения
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to db: %s", err)
	}

	if err := db.AutoMigrate(&models.Paste{}, &models.User{}); err != nil {
		log.Fatalf("Migration error: %s", err)
	}

	pasteRepo := repos.NewPasteRepos(db)
	pasteHandler := handlers.NewPasteHandler(pasteRepo)
	userRepo := repos.NewUser(db)
	userHandler := handlers.NewUserHandler(userRepo)

	router := gin.Default()

	router.POST("/pastes/new_paste", pasteHandler.CreatePaste)
	router.GET("/pastes/:hash", pasteHandler.GetPaste)
	router.POST("/users/registration", userHandler.RegUser)
	router.GET("/users/auth", userHandler.AuthUser)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := fmt.Sprintf(":%s", serverPort)
	log.Printf("Server is running at http://localhost%s/", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Can't start server: %s", err)
	}
}
