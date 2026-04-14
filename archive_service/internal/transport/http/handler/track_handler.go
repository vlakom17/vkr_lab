package handler

import (
	"net/http"

	"charts-archive-service/internal/service"
	"charts-archive-service/internal/transport/http/response"

	"github.com/gin-gonic/gin"
)

type TrackHandler struct {
	service *service.TrackService
}

func NewTrackHandler(s *service.TrackService) *TrackHandler {
	return &TrackHandler{service: s}
}

func (h *TrackHandler) SearchTracks(c *gin.Context) {

	query := c.Query("q")

	tracks, err := h.service.SearchTracks(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to search tracks",
		})
		return
	}

	var result []response.TrackResponse

	for _, t := range tracks {
		result = append(result, response.TrackResponse{
			ID:     t.ID.String(),
			Artist: t.Artist,
			Title:  t.Title,
		})
	}

	c.JSON(http.StatusOK, result)
}
