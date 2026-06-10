package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"yanwit/api/cache"
	"yanwit/api/models"
	"yanwit/api/queue"
	"yanwit/api/repository"
)

// CreateTweet публикует новый твит
func CreateTweet(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req models.CreateTweetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tweet := &models.Tweet{
		AuthorID:    userID,
		Text:        req.Text,
		ReplyToID:   req.ReplyToID,
		RetweetOfID: req.RetweetOfID,
	}

	if err := repository.CreateTweet(tweet); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tweet"})
		return
	}

	// Публикуем событие в RabbitMQ для асинхронной раздачи
	if err := queue.PublishTweetEvent(tweet.ID, userID); err != nil {
		log.Printf("Warning: Failed to publish tweet event: %v", err)
		// Не возвращаем ошибку пользователю, так как твит уже сохранён
	}

	// Получаем информацию об авторе
	author, _ := repository.GetUserByID(userID)
	tweet.Author = author

	c.JSON(http.StatusCreated, tweet)
}

// GetTimeline возвращает ленту твитов пользователя
func GetTimeline(c *gin.Context) {
	userID := c.GetInt64("user_id")
	ctx := c.Request.Context() // используем контекст запроса

	// Параметры пагинации
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 50 {
		limit = 50
	}

	// Пытаемся получить ID твитов из Redis
	timelineKey := "timeline:" + itoa(userID)
	tweetIDs, err := cache.RedisClient.LRange(ctx, timelineKey, int64(offset), int64(offset+limit-1)).Result()

	var tweets []*models.Tweet

	if err == nil && len(tweetIDs) > 0 {
		// Кэш hit: получаем твиты по ID из БД
		var ids []int64
		for _, idStr := range tweetIDs {
			id, parseErr := strconv.ParseInt(idStr, 10, 64)
			if parseErr == nil {
				ids = append(ids, id)
			}
		}
		
		// Получаем твиты по ID
		if len(ids) > 0 {
			tweets, err = repository.GetTweetsByIDs(ids)
			if err != nil {
				log.Printf("Failed to get tweets by IDs: %v", err)
				tweets = []*models.Tweet{}
			}
		}
	}

	// Если Redis пуст или произошла ошибка, загружаем из БД
	if len(tweets) == 0 {
		followees, err := repository.GetFollowees(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get followees"})
			return
		}

		authorIDs := append(followees, userID)
		tweets, err = repository.GetTweetsByAuthorIDs(authorIDs, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get timeline"})
			return
		}
	}

	// Загружаем информацию об авторах
	for _, tweet := range tweets {
		author, _ := repository.GetUserByID(tweet.AuthorID)
		tweet.Author = author
	}

	hasMore := len(tweets) == limit
	c.JSON(http.StatusOK, models.TimelineResponse{
		Tweets:  tweets,
		HasMore: hasMore,
	})
}

// GetUserTweets возвращает твиты конкретного пользователя
func GetUserTweets(c *gin.Context) {
	username := c.Param("username")

	user, err := repository.GetUserByUsername(username)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	tweets, err := repository.GetTweetsByUserID(user.ID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user tweets"})
		return
	}

	// Загружаем информацию об авторе
	for _, tweet := range tweets {
		tweet.Author = user
	}

	c.JSON(http.StatusOK, gin.H{
		"user":   user,
		"tweets": tweets,
	})
}

// itoa преобразует int64 в строку
func itoa(n int64) string {
	return strconv.FormatInt(n, 10)
}