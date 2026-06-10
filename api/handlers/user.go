package handlers

import (
    "net/http"
	"strconv"

    "github.com/gin-gonic/gin"
    "yanwit/api/models"
    "yanwit/api/repository"
)

// SearchUsers ищет пользователей по username
func SearchUsers(c *gin.Context) {
    query := c.Query("q")
    if query == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
        return
    }

    var users []models.User
    sqlQuery := `SELECT id, username, email, avatar_url, bio 
                 FROM users 
                 WHERE username ILIKE $1 
                 ORDER BY username
                 LIMIT 20`

    err := repository.DB.Select(&users, sqlQuery, "%"+query+"%")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
        return
    }

    // Очищаем пароль
    for i := range users {
        users[i].PasswordHash = ""
    }

    c.JSON(http.StatusOK, users)
}

// GetUserStats возвращает статистику пользователя (количество твитов, подписчиков, подписок)
func GetUserStats(c *gin.Context) {
	// Получаем ID из query-параметра ?user_id=123
    userIDStr := c.Query("user_id")
    if userIDStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "user_id query parameter is required"})
        return
    }

    userID, err := strconv.ParseInt(userIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var tweetCount int
    repository.DB.Get(&tweetCount, "SELECT COUNT(*) FROM tweets WHERE author_id = $1", userID)

    var followersCount int
    repository.DB.Get(&followersCount, "SELECT COUNT(*) FROM follows WHERE followee_id = $1", userID)

    var followingCount int
    repository.DB.Get(&followingCount, "SELECT COUNT(*) FROM follows WHERE follower_id = $1", userID)

    c.JSON(http.StatusOK, gin.H{
        "tweet_count":    tweetCount,
        "followers_count": followersCount,
        "following_count": followingCount,
    })
}