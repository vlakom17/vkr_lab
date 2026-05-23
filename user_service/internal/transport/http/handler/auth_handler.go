package handler

import (
	"net/http"
	"strings"

	"charts-user-service/internal/service"
	"charts-user-service/internal/transport/http/request"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {

	var req request.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	token := c.GetHeader("Authorization")
	if token != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "already authorized",
		})
		return
	}
	req.Email = strings.ToLower(req.Email)
	user, err := h.authService.Register(
		c.Request.Context(),
		req.Name,
		req.Email,
		req.Password,
		req.About,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *AuthHandler) Login(c *gin.Context) {

	token := c.GetHeader("Authorization")
	if token != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "already authorized",
		})
		return
	}

	var req request.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	req.Email = strings.ToLower(req.Email)
	token, err := h.authService.Login(
		c.Request.Context(),
		req.Email,
		req.Password,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	err := h.authService.Logout(
		c.Request.Context(),
		token,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) GetUserIDByToken(c *gin.Context) {

	userID, _ := c.Get("user_id")

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
	})
}
