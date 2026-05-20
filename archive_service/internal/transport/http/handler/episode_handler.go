package handler

import (
	"charts-archive-service/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EpisodeHandler struct {
	service *service.EpisodeService
}

func NewEpisodeHandler(s *service.EpisodeService) *EpisodeHandler {
	return &EpisodeHandler{service: s}
}

func (h *EpisodeHandler) GetEpisode(c *gin.Context) {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	episode, err := h.service.GetEpisode(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if episode == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "episode not found"})
		return
	}

	c.JSON(http.StatusOK, episode)
}

func (h *EpisodeHandler) GetLatestEpisodesInternal(c *gin.Context) {
	limitStr := c.Query("limit")

	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	episodes, err := h.service.GetLatestEpisodesWithTracks(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, episodes)
}

func (h *EpisodeHandler) GetEpisodesByChart(c *gin.Context) {
	chartIDParam := c.Param("chart_id")

	chartID, err := uuid.Parse(chartIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chart_id"})
		return
	}

	episodes, err := h.service.GetEpisodesByChart(c.Request.Context(), chartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, episodes)
}

func (h *EpisodeHandler) GetNearestLeftEpisode(c *gin.Context) {
	chartIDStr := c.Query("chart_id")
	dateStr := c.Query("date")

	chartID, err := uuid.Parse(chartIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chart_id"})
		return
	}

	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"})
		return
	}

	episode, err := h.service.GetNearestLeftEpisode(c.Request.Context(), chartID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if episode == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "episode not found"})
		return
	}

	c.JSON(http.StatusOK, episode)
}

func (h *EpisodeHandler) GetLatestEpisodesPage(c *gin.Context) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	episodes, err := h.service.GetLatestEpisodesPage(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, episodes)
}
