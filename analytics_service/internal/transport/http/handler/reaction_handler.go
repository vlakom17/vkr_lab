package handler

import (
	"charts-analytics-service/internal/service"
	"charts-analytics-service/internal/transport/http/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReactionHandler struct {
	service *service.ReactionService
}

func NewReactionHandler(s *service.ReactionService) *ReactionHandler {
	return &ReactionHandler{service: s}
}

func (h *ReactionHandler) GetReactionStats(c *gin.Context) {
	chartIDParam := c.Param("chart_id")

	chartID, err := uuid.Parse(chartIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chart_id"})
		return
	}

	stats, err := h.service.GetReactionStats(c.Request.Context(), chartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *ReactionHandler) GetMostPopularChartIDs(c *gin.Context) {
	limitStr := c.Query("limit")

	limit := 25
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	ids, err := h.service.GetMostPopularChartIDs(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.ChartIDsResponse{
		ChartIDs: ids,
	})
}

func (h *ReactionHandler) GetUserLikedChartIDs(c *gin.Context) {
	userIDParam := c.Param("user_id")

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	ids, err := h.service.GetUserLikedChartIDs(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.ChartIDsResponse{
		ChartIDs: ids,
	})
}

func (h *ReactionHandler) GetUserDislikedChartIDs(c *gin.Context) {
	userIDParam := c.Param("user_id")

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	ids, err := h.service.GetUserDislikedChartIDs(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.ChartIDsResponse{
		ChartIDs: ids,
	})
}
