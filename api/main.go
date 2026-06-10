package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"yanwit/api/cache"
	"yanwit/api/handlers"
	"yanwit/api/middleware"
	"yanwit/api/queue"
	"yanwit/api/repository"
)

func main() {
	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Подключение к PostgreSQL
	if err := repository.InitDB(); err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer repository.CloseDB()

	// Подключение к Redis
	if err := cache.InitRedis(); err != nil {
		log.Println("Warning: Redis connection failed: ", err)
	} else {
		defer cache.CloseRedis()
	}

	// Подключение к RabbitMQ
	if err := queue.InitRabbitMQ(); err != nil {
		log.Println("Warning: RabbitMQ connection failed: ", err)
	} else {
		defer queue.CloseRabbitMQ()
	}

	// Настройка Gin
	router := gin.Default()

	router.Use(middleware.CORS())

	// Health check эндпоинт (публичный)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "ok",
			"service":  "yanwit-api",
			"database": "connected",
		})
	})

	// Публичные маршруты (не требуют авторизации)
	auth := router.Group("/api")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// Защищённые маршруты (требуют JWT токен)
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Профиль
		protected.GET("/profile", func(c *gin.Context) {
			userID := c.GetInt64("user_id")
			username := c.GetString("username")
			user, _ := repository.GetUserByID(userID)
			c.JSON(200, gin.H{
				"user_id":  userID,
				"username": username,
				"user":     user,
			})
		})

		// Твиты
		protected.POST("/tweets", handlers.CreateTweet)
		protected.GET("/timeline", handlers.GetTimeline)
		protected.GET("/users/:username/tweets", handlers.GetUserTweets)
		
		// AI маршруты
		protected.POST("/ai/alternatives", handlers.GenerateAlternatives)
		protected.POST("/ai/detect-manipulation", handlers.DetectManipulation)

		// Поиск пользователя
		protected.GET("/users/search", handlers.SearchUsers)

		// Подписки
		protected.POST("/users/:id/follow", handlers.FollowUser)
		protected.DELETE("/users/:id/follow", handlers.UnfollowUser)

		// Получение статистики пользователя
		protected.GET("/users/stats", handlers.GetUserStats)
	}

	// Запуск сервера
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Yanwit API starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}