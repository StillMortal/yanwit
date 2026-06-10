package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"yanwit/api/ai"
)

type AlternativesRequest struct {
	Text  string `json:"text" binding:"required"`
	Count int    `json:"count" binding:"min=1,max=5"`
	Style string `json:"style" binding:"oneof=funny professional sarcastic encouraging"`
}

type ManipulationRequest struct {
	Text string `json:"text" binding:"required"`
}

func GenerateAlternatives(c *gin.Context) {
	var req AlternativesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if req.Count == 0 {
		req.Count = 3
	}
	if req.Style == "" {
		req.Style = "funny"
	}
	
	result, err := ai.GetAlternatives(req.Text, req.Count, req.Style)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI service unavailable", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, result)
}

func DetectManipulation(c *gin.Context) {
	var req ManipulationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	result, err := ai.DetectManipulation(req.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI service unavailable", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, result)
}