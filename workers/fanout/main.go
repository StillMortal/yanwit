package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
	"yanwit/api/repository"
)

type TweetEvent struct {
	TweetID  int64 `json:"tweet_id"`
	AuthorID int64 `json:"author_id"`
}

var redisClient *redis.Client

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
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Проверка Redis
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis: ", err)
	}
	log.Println("Redis connected successfully")

	// Подключение к RabbitMQ
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://yanwit:yanwit_pass@localhost:5672/"
	}
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ: ", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel: ", err)
	}
	defer ch.Close()

	// Объявляем очередь
	_, err = ch.QueueDeclare(
		"tweet_fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to declare queue: ", err)
	}

	// Начинаем потреблять сообщения
	msgs, err := ch.Consume(
		"tweet_fanout",
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatal("Failed to register consumer: ", err)
	}

	log.Println("Fanout worker started. Waiting for messages...")
	ctx = context.Background()

	// Обрабатываем сообщения
	for msg := range msgs {
		var event TweetEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		log.Printf("Processing tweet_id=%d, author_id=%d", event.TweetID, event.AuthorID)

		// Получаем всех подписчиков автора
		followers, err := repository.GetFollowers(event.AuthorID)
		if err != nil {
			log.Printf("Failed to get followers: %v", err)
			continue
		}

		log.Printf("Author has %d followers", len(followers))

		// Добавляем ID твита в ленту каждого подписчика (Redis)
		for _, followerID := range followers {
			timelineKey := "timeline:" + itoa(followerID)
			err := redisClient.LPush(ctx, timelineKey, event.TweetID).Err()
			if err != nil {
				log.Printf("Failed to add tweet to timeline for user %d: %v", followerID, err)
				continue
			}
			// Ограничиваем длину ленты 800 последними твитами
			redisClient.LTrim(ctx, timelineKey, 0, 799)
		}

		// Также добавляем твит в ленту самого автора
		authorTimelineKey := "timeline:" + itoa(event.AuthorID)
		redisClient.LPush(ctx, authorTimelineKey, event.TweetID)
		redisClient.LTrim(ctx, authorTimelineKey, 0, 799)

		log.Printf("Fanout completed for tweet_id=%d", event.TweetID)
	}
}

// itoa преобразует int64 в строку
func itoa(n int64) string {
	return strconv.FormatInt(n, 10)
}