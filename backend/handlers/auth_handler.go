package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/services"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos inválidos"})
		return
	}

	if err := h.service.Register(c.Request.Context(), req); err != nil {
		status := http.StatusInternalServerError
		switch err {
		case services.ErrEmailExists:
			status = http.StatusConflict
		default:
			status = http.StatusInternalServerError
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos inválidos"})
		return
	}

	token, expiresIn, setCookie, err := h.service.Login(c.Request.Context(), req, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrInvalidCredentials {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	if setCookie != nil {
		setCookie(c.Writer)
	}

	c.JSON(http.StatusOK, dto.LoginResp{AccessToken: token, ExpiresIn: expiresIn})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refresh, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token requerido"})
		return
	}

	token, expiresIn, setCookie, svcErr := h.service.Refresh(c.Request.Context(), refresh, c.Request.UserAgent(), c.ClientIP())
	if svcErr != nil {
		status := http.StatusInternalServerError
		switch svcErr {
		case services.ErrInvalidRefresh, services.ErrExpiredRefresh:
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": svcErr.Error()})
		return
	}

	if setCookie != nil {
		setCookie(c.Writer)
	}

	c.JSON(http.StatusOK, dto.LoginResp{AccessToken: token, ExpiresIn: expiresIn})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refresh, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token requerido"})
		return
	}

	clearCookie, svcErr := h.service.Logout(c.Request.Context(), refresh)
	if svcErr != nil {
		status := http.StatusInternalServerError
		if svcErr == services.ErrInvalidRefresh {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": svcErr.Error()})
		return
	}

	if clearCookie != nil {
		clearCookie(c.Writer)
	}

	c.Status(http.StatusNoContent)
}
