package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "yanwit/api/repository"
)

// FollowUser подписывает текущего пользователя на другого
func FollowUser(c *gin.Context) {
    userID := c.GetInt64("user_id")
    
    followeeIDStr := c.Param("id")
    followeeID, err := strconv.ParseInt(followeeIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    
    if userID == followeeID {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot follow yourself"})
        return
    }
    
    // Проверяем, существует ли пользователь
    user, err := repository.GetUserByID(followeeID)
    if err != nil || user == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    
    // Проверяем, не подписан ли уже
    isFollowing, _ := repository.IsFollowing(userID, followeeID)
    if isFollowing {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Already following this user"})
        return
    }
    
    // Подписываем
    if err := repository.Follow(userID, followeeID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Successfully followed user"})
}

// UnfollowUser отписывает текущего пользователя от другого
func UnfollowUser(c *gin.Context) {
    userID := c.GetInt64("user_id")
    
    followeeIDStr := c.Param("id")
    followeeID, err := strconv.ParseInt(followeeIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    
    if err := repository.Unfollow(userID, followeeID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow user"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Successfully unfollowed user"})
}