package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IService interface {
	Login(guid, ipv4 string) (map[string]string, error)
	Refresh(accessToken, refreshToken, ipv4 string) (map[string]string, error)
}

type Handler struct {
	service IService
}

func NewHandler(service IService) *Handler {
	return &Handler{
		service: service,
	}
}

func (handler *Handler) Login(c *gin.Context) {
	guid := c.Query("guid")
	err := uuid.Validate(guid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	ipv4 := c.RemoteIP()

	pair, err := handler.service.Login(guid, ipv4)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pair)
	return
}

func (handler *Handler) Refresh(c *gin.Context) {
	accessToken := c.Request.Header.Get("Authorization")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "access token was not provided"})
		return
	}

	var queryData RefreshDTO
	err := c.ShouldBindJSON(&queryData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	ipv4 := c.ClientIP()

	pair, err := handler.service.Refresh(accessToken, queryData.RefreshToken, ipv4)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pair)
	return
}
