package handler

import (
	"net/http"
	"strconv"

	"charts-analytics-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RecommendationHandler struct {
	service *service.RecommendationService
}

func NewRecommendationHandler(s *service.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{
		service: s,
	}
}

func (h *RecommendationHandler) GetRecommendations(c *gin.Context) {

	limitStr := c.Query("limit")

	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	userIDVal, exists := c.Get("user_id")

	var userID uuid.UUID
	if exists {
		if id, ok := userIDVal.(uuid.UUID); ok {
			userID = id
		} else {
			userID = uuid.Nil
		}
	} else {
		userID = uuid.Nil
	}

	recs, err := h.service.GetRecommendations(
		c.Request.Context(),
		userID,
		limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, recs)
}
