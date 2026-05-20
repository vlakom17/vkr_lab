package handler

import (
	"charts-chart-service/internal/domain/event"
	"charts-chart-service/internal/service"
	"charts-chart-service/internal/transport/http/request"
	"charts-chart-service/internal/transport/http/response"
	"charts-chart-service/internal/utilits"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChartHandler struct {
	service *service.ChartService
}

func NewChartHandler(s *service.ChartService) *ChartHandler {
	return &ChartHandler{service: s}
}

func (h *ChartHandler) CreateChart(c *gin.Context) {
	var req request.CreateChartRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id"})
		return
	}

	chart, err := h.service.CreateChart(
		c.Request.Context(),
		userID,
		req.Title,
		req.Genre,
		req.PositionCount,
		req.Description,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := response.ChartResponse{
		ID:            chart.ID.String(),
		Title:         chart.Title,
		Genre:         chart.Genre,
		PositionCount: chart.PositionCount,
		Description:   chart.Description,
		UserID:        chart.UserID.String(),
	}

	c.JSON(http.StatusCreated, resp)
}
func (h *ChartHandler) PatchChart(c *gin.Context) {
	var req request.PatchChartRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id"})
		return
	}

	idParam := c.Param("id")
	chartID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	chart, err := h.service.PatchChart(
		c.Request.Context(),
		userID,
		chartID,
		req.Title,
		req.Genre,
		req.PositionCount,
		req.Description,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chart)
}

func (h *ChartHandler) GetChartByID(c *gin.Context) {

	idParam := c.Param("id")
	chartID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var userID *uuid.UUID

	if userIDRaw, exists := c.Get("user_id"); exists {
		if id, ok := userIDRaw.(uuid.UUID); ok {
			userID = &id
		}
	}

	chart, err := h.service.GetChartByID(
		c.Request.Context(),
		chartID,
		userID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chart)
}

func (h *ChartHandler) GetChartByIDWithoutView(c *gin.Context) {

	idParam := c.Param("id")
	chartID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	chart, err := h.service.GetChartByIDWithoutView(
		c.Request.Context(),
		chartID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chart)
}

func (h *ChartHandler) GetMyChart(c *gin.Context) {

	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id"})
		return
	}

	chart, err := h.service.GetMyChart(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chart)
}

func (h *ChartHandler) SetReaction(c *gin.Context) {

	idParam := c.Param("id")
	chartID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chart id"})
		return
	}

	var req request.SetReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id"})
		return
	}

	err = h.service.SetReaction(
		c.Request.Context(),
		userID,
		chartID,
		event.ReactionType(req.Type),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ChartHandler) CreateEpisode(c *gin.Context) {

	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id"})
		return
	}

	chartID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chart id"})
		return
	}

	var req request.CreateEpisodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.CreateEpisodeSnapshot(
		c.Request.Context(),
		userID,
		chartID,
		req.Tracks,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ChartHandler) GetMostPopularCharts(c *gin.Context) {

	limitStr := c.Query("limit")
	limit := 25

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	charts, err := h.service.GetMostPopularCharts(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]response.ChartResponse, 0, len(charts))
	for _, ch := range charts {
		resp = append(resp, utilits.MapChartToResponse(ch))
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ChartHandler) GetMyLikedCharts(c *gin.Context) {

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	userID := userIDVal.(uuid.UUID)

	charts, err := h.service.GetUserLikedCharts(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	resp := make([]response.ChartResponse, 0, len(charts))
	for _, ch := range charts {
		resp = append(resp, utilits.MapChartToResponse(ch))
	}

	c.JSON(200, resp)
}

func (h *ChartHandler) GetMyDislikedCharts(c *gin.Context) {

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	userID := userIDVal.(uuid.UUID)

	charts, err := h.service.GetUserDislikedCharts(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	resp := make([]response.ChartResponse, 0, len(charts))
	for _, ch := range charts {
		resp = append(resp, utilits.MapChartToResponse(ch))
	}

	c.JSON(200, resp)
}
