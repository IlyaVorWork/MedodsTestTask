package auth

import (
	"net/http"

	"MedodsTestTask/internal/pkg"
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

	// Получение id пользователя из параметра запроса
	guid := c.Query("guid")
	err := uuid.Validate(guid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// Получение IP клиента
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

	// Получение access token'а из заголовка авторизации
	accessToken := c.Request.Header.Get("Authorization")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "access token was not provided"})
		return
	}

	// Получение данных из тела запроса
	var queryData RefreshDTO
	err := c.ShouldBindJSON(&queryData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// Проверка на совпадения предоставленных access и refresh токенов
	if accessToken == queryData.RefreshToken {
		c.JSON(http.StatusBadRequest, gin.H{"Error": pkg.ProvidedTokensSimilar.Error()})
	}

	// Получение IP клиента
	ipv4 := c.ClientIP()

	pair, err := handler.service.Refresh(accessToken, queryData.RefreshToken, ipv4)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pair)
	return
}
